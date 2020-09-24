package flowscheduler

import "fmt"

const (
	NODE_TASK = 1
	NODE_RESOURCE = 2
	NODE_SINK = 3
	NODE_AGGNODE = 4
	//unscheduled node is the same as the job node
	NODE_UNSCEDULED = 5 
	NODE_OTHER = 6
)
type GraphNode struct{
	nid NodeID
	excess int64
	outgoingArcs map[NodeID]*GraphArc
	incomingArcs map[NodeID]*GraphArc
	comment string
	ntype NodeType
	// if this node is a task node or job node,we should set this 
	jd *JobDescriptor
	// if this node is a task node,it needs to be setted
	td *TaskDescriptor
	// if this node is a res node, it needs to be setted
	rd *ResDescriptor
}

func (n *GraphNode)IsTaskNode()bool{
	return n.ntype == NODE_TASK
}

func (n *GraphNode)IsResourceNode()bool{
	return n.ntype == NODE_RESOURCE
}

func (node *GraphNode)Export()string{
	ans := ""
	if node.IsResourceNode() {
		ans += fmt.Sprintf("c Resource Node,res aviliable is %d\n",node.rd.ResAvailable.Gpu)
	} else if node.IsTaskNode(){
		ans += fmt.Sprintf("c TaskNode, ResRequest is %d,state : %d\n",node.td.ResRequest.Gpu,node.td.State)
	} else {
		ans += "c "+node.comment
	}
	ans += fmt.Sprintf("n %d %d\n",node.nid,node.excess)
	return ans
}