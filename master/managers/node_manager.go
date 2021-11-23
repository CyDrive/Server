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
	HeartBeatTimeout = 500  // in ms
	OfflineTimeout   = 3000 // in second
)

type NodeManager struct {
	env     envs.Env
	nodeMap *sync.Map // map: nodeId -> *Node
	fileMap *sync.Map // map: filePath -> []*Node

	fileTransferor *network.FileTransferor
	nodeNum        int32
	runningNodeNum int32
	replicationNum int32
}

func NewNodeManager(fileTransferor *network.FileTransferor, replicationNum int32) *NodeManager {
	nodeManager := NodeManager{
		nodeMap:        &sync.Map{},
		fileMap:        &sync.Map{},
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
	var nodes []*Node
	if !ok {
		nodes = []*Node{}
		nm.fileMap.Store(filePath, nodes)
	} else {
		nodes := make([]*Node, 0, 1)
		for _, node := range nodesI.([]*Node) {
			if node.State != consts.NodeState_Dropping {
				nodes = append(nodes, node)
			}
		}
		nm.fileMap.Store(filePath, nodes)
	}

	return nodes
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

			if len(nodes) < int(nm.replicationNum) {
				assignNodes := nm.PickNodesExcept(int(nm.replicationNum)-len(nodes), nodes)
				if len(assignNodes) == 0 {
					log.Warnf("no enough node for file: %s", filePath)
				} else {
					for _, assignNode := range assignNodes {
						go func() {
							nm.replica(filePath, nodes[0], assignNode)
							nm.AssignFile(filePath, assignNode)
						}()
					}
				}
			}

			nm.fileMap.Store(key, nodes)

			return true
		})

		// Replica to make the number of replications satisfied
		// to keep durability

	}
}

func (nm *NodeManager) replica(filePath string, src *Node, dest *Node) error {
	fileInfo, err := nm.env.Stat(filePath)
	if err != nil {
		return err
	}

	file := envs.NewPipeFile(fileInfo)

	srcTask := nm.fileTransferor.CreateTask(fileInfo, file, consts.DataTaskType_Upload, 0)
	nm.fileTransferor.CreateTask(fileInfo, file, consts.DataTaskType_Download, 0)

	cond := sync.NewCond(&sync.Mutex{})
	srcTask.OnEnd = func() {
		file.Close()
		cond.Signal()
	}

	cond.Wait()

	return nil
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
