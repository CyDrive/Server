package managers

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CyDrive/models"
	"github.com/CyDrive/rpc"
)

var (
	nextNodeId int32 = 0
)

func GenId() int32 {
	return atomic.AddInt32(&nextNodeId, 1)
}

type Node struct {
	Id                int32
	Usage             int64
	Cap               int64
	LastHeartBeatTime time.Time

	NotifyChan chan interface{}
}

func NewNode(cap, usage int64) *Node {
	return &Node{
		Id:                GenId(),
		Usage:             usage,
		Cap:               cap,
		LastHeartBeatTime: time.Now(),
		NotifyChan:        make(chan interface{}, 100),
	}
}

// type NodeElem struct {
// 	*Node
// 	index int
// }
// type NodePriorityQueue []*NodeElem

// func (q NodePriorityQueue) Len() int { return len(q) }

// func (q NodePriorityQueue) Less(i, j int) bool {
// 	return q[i].Cap-q[i].Usage > q[j].Cap-q[j].Usage
// }

// func (q NodePriorityQueue) Swap(i, j int) {
// 	q[i], q[j] = q[j], q[i]
// }

// func (q *NodePriorityQueue) Push(x interface{}) {
// 	node := x.(*Node)
// 	elem := NodeElem {
// 		Node: node,
// 		index: ,
// 	}
// 	*q = append(*q, elem)
// }

// func (q *NodePriorityQueue) Pop() interface{} {
// 	queue := *q
// 	n := len(queue)
// 	elem := queue[n-1]

// 	queue[n-1] = nil
// 	*q = queue[:n-1]

// 	return elem
// }

const (
	HeartBeatTimeout = 500 // in ms
)

type NodeManager struct {
	nodeMap *sync.Map
	fileMap *sync.Map // map: filePath -> []*Node
	nodeNum int32
}

func NewNodeManager() *NodeManager {
	nodeManager := NodeManager{
		nodeMap: &sync.Map{},
		fileMap: &sync.Map{},
		nodeNum: 0,
	}

	return &nodeManager
}

func (nm *NodeManager) AddNode(node *Node) {
	nm.nodeMap.Store(node.Id, node)
	nm.nodeNum++
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

func (nm *NodeManager) CreateSendFileTask(accountId int32, req *rpc.CreateSendFileTaskNotify) error {
	nodeI, ok := nm.fileMap.Load(accountId)
	if !ok {
		return fmt.Errorf("no such file")
	}

	node := nodeI.(*Node)
	node.NotifyChan <- req

	return nil
}

func (nm *NodeManager) CreateRecvFileTask(account *models.Account, req *rpc.CreateRecvFileTaskNotify) error {
	nodeI, ok := nm.fileMap.Load(account.Id)
	var node *Node
	if !ok {
		// Assign a node to serve the account
		node = nm.PickNode()
		if node == nil {
			return fmt.Errorf("No node to serve!")
		}
		node.Usage += req.FileInfo.Size
		nm.fileMap.Store(account.Id, node)
	} else {
		node = nodeI.(*Node)
	}

	node.NotifyChan <- req
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
