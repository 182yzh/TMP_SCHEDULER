package flowscheduler

import "fmt"

const (
	ARC_RUNNING = 1
	ARC_OTHER = 2
)
type GraphArc struct{
	srcID NodeID
	dstID NodeID
	src   *GraphNode
	dst   *GraphNode
	cost   uint64
	capLower uint64
	capUpper uint64
	comment string
	atype ArcType
}

func (arc *GraphArc)Export()string{
	var ans string
	//ans += "c "+arc.comment
	ans += fmt.Sprintf("a %d %d %d %d %d\n",arc.srcID,arc.dstID,arc.capLower,arc.capUpper,arc.cost)
	return ans 
}