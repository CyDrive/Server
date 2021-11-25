package managers

import (
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/envs"
	"github.com/CyDrive/network"
	"github.com/CyDrive/rpc"
	"github.com/CyDrive/types"
	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
)

var (
	nextNodeId int32 = 0
)

func GenId() int32 {
	return atomic.AddInt32(&nextNodeId, 1)
}

type Node struct {
	Id                int32
	Addr              string
	Usage             int64
	Cap               int64
	LastHeartBeatTime time.Time
	State             consts.NodeState

	NotifyChan chan *rpc.Notify
}

func NewNode(cap, usage int64, addr string) *Node {
	return &Node{
		Id:                GenId(),
		Addr:              addr,
		Usage:             usage,
		Cap:               cap,
		LastHeartBeatTime: time.Now(),
		State:             consts.NodeState_Starting,
		NotifyChan:        make(chan *rpc.Notify, 100),
	}
}

const (
	HeartBeatTimeout = 1000 // in ms
	OfflineTimeout   = 300  // in second
)

type NodeManager struct {
	env            envs.Env
	nodeMap        *sync.Map // map: nodeId -> *Node
	fileMap        *sync.Map // map: filePath -> []*Node
	replicatingMap *sync.Map // map: filePath -> []*Node

	fileTransferor *network.FileTransferor
	nodeNum        int32
	runningNodeNum int32
	replicationNum int32
}

func NewNodeManager(fileTransferor *network.FileTransferor, replicationNum int32) *NodeManager {
	nodeManager := NodeManager{
		nodeMap:        &sync.Map{},
		fileMap:        &sync.Map{},
		replicatingMap: &sync.Map{},

		fileTransferor: fileTransferor,
		nodeNum:        0,
		runningNodeNum: 0,
		replicationNum: replicationNum,
	}

	go nodeManager.healthMaintenance()

	return &nodeManager
}

func (nm *NodeManager) SetEnv(env envs.Env) {
	nm.env = env
}

func (nm *NodeManager) changeNodeState(node *Node, state consts.NodeState) {
	if node.State == state {
		return
	}
	if node.State == consts.NodeState_Running {
		atomic.AddInt32(&nm.runningNodeNum, -1)
	}
	if state == consts.NodeState_Running {
		atomic.AddInt32(&nm.runningNodeNum, 1)
	}
	node.State = state
}

func (nm *NodeManager) AddNode(node *Node) {
	nm.nodeMap.Store(node.Id, node)
	atomic.AddInt32(&nm.nodeNum, 1)
}

func (nm *NodeManager) GetNode(id int32) *Node {
	value, ok := nm.nodeMap.Load(id)
	if !ok {
		return nil
	}

	return value.(*Node)
}

func (nm *NodeManager) ChangeNodeState(id int32, state consts.NodeState) {
	node := nm.GetNode(id)
	if node != nil {
		nm.changeNodeState(node, state)
	}
}

func (nm *NodeManager) GetNodesByFilePath(filePath string) []*Node {
	nodesI := nm.getNodesByFilePath(filePath)

	nodes := make([]*Node, 0, 1)
	for _, node := range nodesI {
		if node.State == consts.NodeState_Running {
			nodes = append(nodes, node)
		}
	}

	if len(nodes) == 0 { // this is a new file, assign a node to serve it
		node := nm.PickNode()
		nodes = nm.AssignFile(filePath, node)
	}

	return nodes
}

func (nm *NodeManager) getNodesByFilePath(filePath string) []*Node {
	nodesI, ok := nm.fileMap.Load(filePath)
	ret := make([]*Node, 0, 2)
	if !ok {
		ret = []*Node{}
	} else {
		nodes := nodesI.([]*Node)
		for _, node := range nodes {
			if node.State != consts.NodeState_Dropping {
				ret = append(ret, node)
			}
		}
	}

	nm.fileMap.Store(filePath, ret)
	return ret
}

func (nm *NodeManager) AssignFile(filePath string, node *Node) []*Node {
	nodes := nm.getNodesByFilePath(filePath)
	nodes = append(nodes, node)
	nm.fileMap.Store(filePath, nodes)

	return nodes
}

func (nm *NodeManager) dropNode(id int32) {
	node := nm.GetNode(id)
	if node.State == consts.NodeState_Dropping {
		nm.nodeMap.Delete(id)
		atomic.AddInt32(&nm.nodeNum, -1)
	} else {
		panic("forget to set the node state to Dropping")
	}
}

func (nm *NodeManager) healthMaintenance() {
	for {
		// Remove offline nodes reaching the OfflineTimeout
		removedNodes := make([]int32, 0, 1)
		nm.nodeMap.Range(func(key, value interface{}) bool {
			id := key.(int32)
			node := value.(*Node)
			if time.Since(node.LastHeartBeatTime).Milliseconds() >= HeartBeatTimeout {
				nm.changeNodeState(node, consts.NodeState_Offline)
			}
			if time.Since(node.LastHeartBeatTime).Seconds() >= OfflineTimeout {
				nm.changeNodeState(node, consts.NodeState_Dropping)
				removedNodes = append(removedNodes, id)
			}
			return true
		})
		for _, id := range removedNodes {
			nm.dropNode(id)
		}

		// Remove dropped node from fileMap
		nm.fileMap.Range(func(key, value interface{}) bool {
			filePath := key.(string)
			old := value.([]*Node)

			nodes := make([]*Node, 0, 1)
			for _, node := range old {
				if node.State != consts.NodeState_Dropping {
					nodes = append(nodes, node)
				}
			}

			// Need to replica
			if len(nodes) < int(nm.replicationNum) {
				assignNodes := filterNodesByState(
					nm.PickNodesExcept(int(nm.replicationNum)-len(nodes), nodes),
					consts.NodeState_Running)

				assignNodes = nm.filterNodesByReplicatingMap(filePath, assignNodes)

				srcNodes := filterNodesByState(nodes, consts.NodeState_Running)

				if len(assignNodes) == 0 {
					log.Warnf("no enough node for file: %s, nodes_num=%d replication_num=%d", filePath, len(nodes), nm.replicationNum)
				} else if len(srcNodes) == 0 {
					log.Infof("All nodes are still starting, replica later")
				} else {
					for _, assignNode := range assignNodes {
						nm.addReplicas(filePath, assignNode)
						go func(filePath string, src, dest *Node) {
							nm.replica(filePath, src, dest)
							nm.AssignFile(filePath, dest)
						}(filePath, srcNodes[0], assignNode)
					}
				}
			}

			nm.fileMap.Store(key, nodes)

			return true
		})

		time.Sleep(HeartBeatTimeout * time.Millisecond)
	}
}

func (nm *NodeManager) replica(filePath string, src *Node, dest *Node) error {
	log.Infof("replica file=%s from src=%+v to dest=%+v", filePath, src, dest)

	fileInfo, err := nm.env.Stat(filePath)
	if err != nil {
		return err
	}

	file := envs.NewPipeFile(fileInfo)

	srcTask := nm.fileTransferor.CreateTask(fileInfo, file, consts.DataTaskType_Upload, 0)
	destTask := nm.fileTransferor.CreateTask(fileInfo, file, consts.DataTaskType_Download, 0)

	src.NotifyChan <- utils.PackTransferFileNotification(srcTask.Id, config.IpAddr, filePath, consts.DataTaskType_Upload)
	dest.NotifyChan <- utils.PackTransferFileNotification(destTask.Id, config.IpAddr, filePath, consts.DataTaskType_Download)

	wg := sync.WaitGroup{}
	wg.Add(1)
	srcTask.OnEnd = func() {
		defer wg.Done()
		file.Close()
	}

	wg.Wait()

	return nil
}

func (nm *NodeManager) addReplicas(filePath string, node *Node) {
	nodes := make([]*Node, 0, 2)
	nodesI, ok := nm.replicatingMap.Load(filePath)
	if ok {
		nodes = nodesI.([]*Node)
	}
	nodes = append(nodes, node)

	nm.replicatingMap.Store(filePath, nodes)
}

func (nm *NodeManager) PrepareReadFile(taskId types.TaskId, filePath string) error {
	nodesI, ok := nm.fileMap.Load(filePath)
	if !ok {
		return os.ErrNotExist
	}

	// todo: load-balance
	// now we just pick the first node
	nodes := nodesI.([]*Node)
	node := nodes[0]

	notify := utils.PackTransferFileNotification(taskId, config.IpAddr, filePath, consts.DataTaskType_Upload)
	log.Infof("send notification to channel, notify=%+v", notify)
	node.NotifyChan <- notify

	return nil
}

func (nm *NodeManager) PrepareWriteFile(taskId types.TaskId, filePath string) error {
	nodes := nm.GetNodesByFilePath(filePath)

	// todo: load-balance
	// now we just pick the first node
	for _, node := range nodes {
		node.NotifyChan <- utils.PackTransferFileNotification(taskId, config.IpAddr, filePath, consts.DataTaskType_Download)
	}

	return nil
}

func (nm *NodeManager) NotifyDeleteFile(filePath string) {
	nm.fileMap.Delete(filePath)

	nodes := nm.GetNodesByFilePath(filePath)
	for _, node := range nodes {
		node.NotifyChan <- utils.PackDeleteFileNotification(filePath)
	}
}

func (nm *NodeManager) GetNotifyChan(nodeId int32) (<-chan *rpc.Notify, bool) {
	nodeI, ok := nm.nodeMap.Load(nodeId)
	if !ok {
		return nil, false
	}

	node := nodeI.(*Node)
	return node.NotifyChan, true
}

// return: the first num nodes with highest (cap - usage)
// if num is greater than the number of nodes, return all nodes
// note: the returned slice is not sorted!
func (nm *NodeManager) PickNodes(num int) []*Node {
	if num >= int(nm.runningNodeNum) {
		num = int(nm.runningNodeNum)
	}

	allNodes := make([]*Node, 0, nm.nodeNum)
	nm.nodeMap.Range(func(key, value interface{}) bool {
		allNodes = append(allNodes, value.(*Node))
		return true
	})
	sort.Slice(allNodes, func(i, j int) bool {
		if allNodes[i].State == allNodes[j].State {
			return allNodes[i].Cap-allNodes[i].Usage > allNodes[j].Cap-allNodes[j].Usage
		}
		return allNodes[i].State == consts.NodeState_Running
	})

	nodes := make([]*Node, 0, num)

	nodes = append(nodes, allNodes[:num]...)

	return nodes
}

func (nm *NodeManager) PickNodesExcept(num int, exceptNodes []*Node) []*Node {
	nodes := nm.PickNodes(num)
	ret := make([]*Node, 0, len(nodes))
	for _, node := range nodes {
		shouldAdd := true
		for _, except := range exceptNodes {
			if node.Id == except.Id {
				shouldAdd = false
				break
			}
		}

		if shouldAdd {
			ret = append(ret, node)
		}
	}

	return ret
}

func (nm *NodeManager) PickNode() *Node {
	nodes := nm.PickNodes(1)
	if len(nodes) == 0 {
		return nil
	}

	return nodes[0]
}

func (nm *NodeManager) filterNodesByReplicatingMap(filePath string, nodes []*Node) []*Node {
	ret := make([]*Node, 0, len(nodes))
	replicatingNodesI, ok := nm.replicatingMap.Load(filePath)
	replicatingNodes := []*Node{}
	if ok {
		replicatingNodes = replicatingNodesI.([]*Node)
	}

	for _, node := range nodes {
		if !isInSlice(node, replicatingNodes) {
			ret = append(ret, node)
		}
	}

	return ret
}

func isInSlice(checkNode *Node, nodes []*Node) bool {
	for _, node := range nodes {
		if checkNode == node {
			return true
		}
	}
	return false
}

func filterNodesByState(nodes []*Node, state consts.NodeState) []*Node {
	ret := make([]*Node, 0, len(nodes))
	for _, node := range nodes {
		if node.State == state {
			ret = append(ret, node)
		}
	}

	return ret
}
