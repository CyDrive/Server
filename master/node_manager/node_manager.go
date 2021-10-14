package node_manager

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var (
	nextNodeId int32 = 0
)

func GenId() int32 {
	return atomic.AddInt32(&nextNodeId, 1)
}

type Node struct {
	Id                int32
	Cap               int64
	Usage             int64
	LastHeartBeatTime time.Time

	Conn *websocket.Conn
}

func NewNode(cap, usage int64) *Node {
	return &Node{
		Id:                GenId(),
		Cap:               cap,
		Usage:             usage,
		LastHeartBeatTime: time.Now(),
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
	nodeNum int32
}

var nodeManager *NodeManager = NewNodeManager()

func NewNodeManager() *NodeManager {
	nodeManager := NodeManager{
		nodeMap: &sync.Map{},
		nodeNum: 0,
	}

	return &nodeManager
}

func GetNodeManager() *NodeManager {
	return nodeManager
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
