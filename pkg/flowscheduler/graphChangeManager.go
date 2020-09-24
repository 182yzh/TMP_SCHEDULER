package flowscheduler

type GraphChangeManager struct {
	fg *Graph
}


func (gcm *GraphChangeManager)Init(){
	gcm.fg = new(Graph)
	gcm.fg.Init()
}

func (gcm *GraphChangeManager)AddNode(excess int64,ntype uint8,comment string)*GraphNode{
	node := gcm.fg.AddNode(excess,ntype,comment)
	return node
}

func (gcm *GraphChangeManager)AddArc(src,dst *GraphNode,
									 lower,upper uint64,
									 cost uint64,
									 comment string)*GraphArc{
	arc := gcm.fg.AddArc(src,dst,lower,upper,cost,comment)
	return arc
}

//==--== maybe this is a problem
func (gcm *GraphChangeManager)RemoveNode(node *GraphNode){
	gcm.fg.RemoveNode(node)
}

func (gcm *GraphChangeManager)ChangeArc(arc *GraphArc,lower,upper,cost uint64,comment string){
	arc.capLower,arc.capUpper = lower,upper
	arc.cost = cost
	arc.comment = "change arc|"+arc.comment
}
func (gcm *GraphChangeManager)RemoveArc(arc *GraphArc){
	arc.capLower,arc.capUpper = 0,0
	arc.comment = "Remove "+arc.comment
	gcm.fg.RemoveArc(arc)
}

func (gcm *GraphChangeManager)GetArc(src,dst *GraphNode)*GraphArc{
	return gcm.fg.GetArc(src,dst)
}

func (gcm *GraphChangeManager)GetNode(nid NodeID)*GraphNode{
	return gcm.fg.GetNode(nid)
}
