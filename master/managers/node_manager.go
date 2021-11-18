package managers

import (
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CyDrive/config"
	"github.com/CyDrive/network"
	"github.com/CyDrive/types"
	"github.com/CyDrive/utils"
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

	NotifyChan chan interface{}
}

func NewNode(cap, usage int64, addr string) *Node {
	return &Node{
		Id:                GenId(),
		Addr:              addr,
		Usage:             usage,
		Cap:               cap,
		LastHeartBeatTime: time.Now(),
		NotifyChan:        make(chan interface{}, 100),
	}
}

const (
	HeartBeatTimeout = 500 // in ms
)

type NodeManager struct {
	nodeMap *sync.Map // map: nodeId -> *Node
	fileMap *sync.Map // map: filePath -> []*Node

	fileTransferor *network.FileTransferor
	nodeNum        int32
}

func NewNodeManager(fileTransferor *network.FileTransferor) *NodeManager {
	nodeManager := NodeManager{
		nodeMap:        &sync.Map{},
		fileMap:        &sync.Map{},
		fileTransferor: fileTransferor,
		nodeNum:        0,
	}

	return &nodeManager
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

func (nm *NodeManager) GetNodesByFilePath(filePath string) []*Node {
	nodesI, ok := nm.fileMap.Load(filePath)
	if !ok {
		return nil
	}

	nodes := nodesI.([]*Node)

	return nodes
}

func (nm *NodeManager) DropNode(node *Node) {
	panic("unimplemented")
}

func (nm *NodeManager) NodeHealthMaintenance() {
	panic("unimplemented")
	// for {
	// 	nodes := nm.nodes
	// 	removedNodes := make([]*Node, 0, 1)

	// 	for _, node := range nodes {
	// 		// HeartBeat timeout, drop this node
	// 		if time.Now().Sub(node.LastHeartBeatTime).Milliseconds() >= HeartBeatTimeout {
	// 			removedNodes = append(removedNodes, node)
	// 		}
	// 	}

	// }
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
	node.NotifyChan <- utils.PackCreateTransferFileTaskNotification(taskId, config.IpAddr, filePath)

	return nil
}

func (nm *NodeManager) GetNotifyChan(nodeId int32) (<-chan interface{}, bool) {
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
	if num >= int(nm.nodeNum) {
		num = int(nm.nodeNum)
	}

	allNodes := make([]*Node, 0, nm.nodeNum)
	nm.nodeMap.Range(func(key, value interface{}) bool {
		allNodes = append(allNodes, value.(*Node))
		return true
	})
	sort.Slice(allNodes, func(i, j int) bool {
		return allNodes[i].Cap-allNodes[i].Usage > allNodes[j].Cap-allNodes[j].Usage
	})

	nodes := make([]*Node, 0, num)

	nodes = append(nodes, allNodes[:num]...)

	return nodes
}

func (nm *NodeManager) PickNode() *Node {
	nodes := nm.PickNodes(1)
	if len(nodes) == 0 {
		return nil
	}

	return nodes[0]
}
