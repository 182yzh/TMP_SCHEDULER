package flowscheduler

//import "vlog"
import "fmt"

type Graph struct{
	nodes map[NodeID]*GraphNode
	arcSet  map[*GraphArc]bool
	curID uint64
	unusedID map[NodeID]bool
}

func (g *Graph)Init(){
	g.nodes = make(map[NodeID]*GraphNode)
	g.arcSet = make(map[*GraphArc]bool)
	g.unusedID = make(map[NodeID]bool)
	g.curID = 0
}

func (g *Graph)GetNextID()NodeID{
	var ans NodeID
	if len(g.unusedID) == 0{
		g.curID++
		ans =  g.curID
	} 
	for id,_ := range g.unusedID{
		ans = id
		delete(g.unusedID,id)
		break
	}
	return ans
}

func (g *Graph)GetNodeByID(nid NodeID)*GraphNode{
	node,_ := g.nodes[nid]
	return node
}

func (g *Graph)AddNode(excess int64,ntype NodeType,comment string)*GraphNode{
	node := new(GraphNode)
	node.nid = g.GetNextID()
	g.nodes[node.nid] = node
	node.excess = excess
	node.ntype = ntype
	node.comment = comment
	node.incomingArcs = make(map[NodeID]*GraphArc)
	node.outgoingArcs = make(map[NodeID]*GraphArc)
	return node
}

func (g *Graph)RemoveNode(node *GraphNode){
	if node == nil {
		//vlog.Dlog("Error g.RemoveNode, node is nil")
	}
	for _,arc := range node.outgoingArcs{
		g.RemoveArc(arc)
	}
	for _,arc := range node.incomingArcs{
		g.RemoveArc(arc)
	}
	delete(g.nodes,node.nid)
	g.unusedID[node.nid] = true
}

func (g *Graph)AddArc(src_,dst_ *GraphNode,
					  lower,upper uint64,
					  cost_ uint64,
					  comment_ string)*GraphArc{
	arc := &GraphArc{
		srcID:src_.nid,
		dstID:dst_.nid,
		capLower:lower,
		capUpper:upper,
		cost:cost_,
		comment:comment_,
		src:src_,
		dst:dst_,
		atype:ARC_OTHER,
	}
	src_.outgoingArcs[dst_.nid] = arc
	dst_.incomingArcs[src_.nid] = arc
	g.arcSet[arc] = true
	return arc
}

func (g *Graph)RemoveArc(arc *GraphArc){
	if arc == nil {
		//vlog.Dlog("Error g.RemoveArc, arc is nil")
	}
	src,dst := arc.src,arc.dst 
	delete(src.outgoingArcs,dst.nid)
	delete(dst.incomingArcs,src.nid)
	delete(g.arcSet,arc)
}

func (g *Graph)ChangeArc(arc *GraphArc,
						 lower,upper uint64,
						 cost uint64){
	if arc == nil {
		//vlog.Dlog("Error g.RemoveArc, arc is nil")
	}
	arc.capLower = lower
	arc.capUpper = upper
	arc.cost = cost
}

func (g *Graph)GetArc(src,dst *GraphNode)*GraphArc{
	if src == nil || dst == nil {
		//vlog.Dlog("Error g.GetArc, src or dst is nil")
		fmt.Println("src or dst is nil\n")
	}
	arc,_ := src.outgoingArcs[dst.nid]
	return arc
}

func (g *Graph)GetNode(nid NodeID)*GraphNode{
	node,_ := g.nodes[nid]
	return node
} 

func (g *Graph)ExportGraph()string{
	ans := "c This is a max-flow min-cost problem\n"
	//ans += fmt.Sprintf("p min %d %d\n",g.curID,len(g.arcSet))
	n :=0
	nodes  := "c Nodes\n"
	for _,node := range g.nodes{
		//if node.excess == 0{
		//	continue
		//}
		nodes += node.Export()
	}
	arcs := "c arcs\n"
	for arc,_ := range g.arcSet{
		if arc.capUpper == 0 {
			continue
		}
		node := arc.src
		if node.ntype ==  NODE_TASK && node.excess == 0{
			continue
		}
		n++
		arcs += arc.Export()
	}
	ans += fmt.Sprintf("p min %d %d\n",g.curID,n)
	return ans + nodes + arcs
}
/*
func (g *Graph)ExportGraph()string{
	ans := "c This is a max-flow min-cost problem\n"
	ans += fmt.Sprintf("p min %d %d\n",g.curID,len(g.arcSet))
	ans += "c Nodes\n"
	for _,node := range g.nodes{
		ans += node.Export()
	}
	ans += "c arcs\n"
	for arc,_ := range g.arcSet{
		ans += arc.Export()
	}
	return ans 
}
*/
