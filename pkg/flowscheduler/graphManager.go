package flowscheduler

import (
	"fmt"
	//"vlog"
	"github.com/golang/glog"
)

// flowgraphmanager only handle the follow:
// map the task to node
// map the res to node
// sinknode,leafnodes,leaf resource id
// running arc from running task to res
// arc from unaggnode to sinknode
// other arc that need to be add will be handle by costmodel

/*
gm中可能需要面对的情况：
添加，删除，更新任务节点
：添加，通过给定的jobdes中的taskdes 添加
：删除，指定任务，删除相应的节点，释放占用的资源（如果在running）
添加，删除，更新资源节点
：由flow schdduler 给定删除任务。（删除指令谁通知flowscheduler？
：添加以及更新，通过resdes进行更新，如果模拟调度则在调度过程中进行更改（只需要更新resdes即可
处理 任务结束/job结束
：任务结束：释放相应资源（通过修该resdes）
处理 task调度
：占用相应资源（实际上调度过程已经更改？task调度后占用资源怎么体现？），
   slover求解过程中的资源是直接占用还是保留
question
调度过程的调度结果出现问题：如果一部分资源不能满足（cpu，memory出现问题，如何处理？）

*/
type GraphManager struct {
	sinkNode       *GraphNode
	taskNodes      map[TaskID]*GraphNode
	resNodes       map[ResID]*GraphNode
	jobUnscheds    map[JobID]*GraphNode
	taskRuningArcs map[TaskID]*GraphArc
	leafNodes      map[NodeID]bool
	leafResID      map[ResID]bool
	//runningTasks   map[TaskID]*GraphNode
	queue          map[*GraphNode]bool
	gcm            *GraphChangeManager
	costmodel      CostModel
}

func (gm *GraphManager) TaskIDToNode(tid TaskID) *GraphNode {
	return gm.taskNodes[tid]
}

func (gm *GraphManager) ResIDToNode(rid ResID) *GraphNode {
	return gm.resNodes[rid]
}

func (gm *GraphManager) JobIDToUnscheNode(jid JobID) *GraphNode {
	return gm.jobUnscheds[jid]
}

func (gm *GraphManager) Init() {
	gm.gcm = new(GraphChangeManager)
	gm.gcm.Init()
	gm.sinkNode = gm.gcm.AddNode(0, NODE_SINK, fmt.Sprintf("Sink Node\n"))
	gm.taskNodes = make(map[TaskID]*GraphNode)
	gm.resNodes = make(map[ResID]*GraphNode)
	gm.jobUnscheds = make(map[JobID]*GraphNode)
	gm.taskRuningArcs = make(map[TaskID]*GraphArc)
	gm.leafNodes = make(map[NodeID]bool)
	gm.leafResID = make(map[ResID]bool)
	gm.costmodel = NewPengChengCostModel(gm)
	gm.queue = make(map[*GraphNode]bool)
	//gm.runningTasks = make(map[TaskID]*GraphNode)
}

func (gm *GraphManager) SetCostModel(cs CostModel) {
	gm.costmodel = cs
}

func (gm *GraphManager) GetResInfo()map[string]uint64{
	ans := make(map[string]uint64)
	for _,node := range gm.resNodes{
		rd := node.rd
		ans[rd.Name] = rd.ResAvailable.Gpu + rd.ResPreempt.Gpu
	}
	return ans
}

func (gm *GraphManager)GetTaskInfo()map[string]string{
	ans := make(map[string]string)
	for tid,tnode := range gm.taskNodes{
		td := tnode.td
		arc,ok := gm.taskRuningArcs[tid]
		info := td.Name
		if ok {
			info+="|R"
		}else {
			info+="|P"
		}
		info+=fmt.Sprintf("|%d|%d",td.Priority,td.ResRequest.Gpu)
		if ok {
			rnode := arc.dst
			ans[info]=rnode.rd.Name + fmt.Sprintf(",TaskState:%v    ",td.State)
		} else {
			ans[info]="Unschedule" + fmt.Sprintf(",TaskState:%v    ",td.State)
		}
	}
	return ans
}

func (gm *GraphManager) AddTaskNode(td *TaskDescriptor) *GraphNode {
	if td == nil {
		//vlog.Dlog("Error: gm.addTaskNode, td is nil")
	}
	if node, ok := gm.taskNodes[td.GetTaskID()]; ok {
		return node
	}
	//fmt.Println("in gm.add TaskNode,gm.tasknodes,add taskNOde:",gm.taskNodes,td.Name)
	node := gm.gcm.AddNode(0, NODE_TASK, "node for "+td.Name+"\n")
	//vlog.Vlog(fmt.Sprintf("gm-add task Node %d", node.nid))
	node.td = td
	node.jd = td.Jd
	gm.taskNodes[td.GetTaskID()] = node
	gm.queue[node] = true
	return node
}

func (gm *GraphManager) RemoveTaskNode(td *TaskDescriptor) NodeID {
	if td == nil {
		return 0
	}
	tid := td.GetTaskID()
	node := gm.TaskIDToNode(tid)

	delete(gm.taskNodes, tid)
	gm.gcm.RemoveNode(node)
	return node.nid
}

func (gm *GraphManager) AddUnscehduleAggNode(jd *JobDescriptor) *GraphNode {
	if jd == nil {
		fmt.Println("140 gm.AddUnscheduleAggNode ,Error Jd is nil")
//vlog.Dlog("Error: gm.AddUnscehduleAggNode, jd is nil")
	}
	jid := jd.GetJobID()
	if unaggnode, ok := gm.jobUnscheds[jid]; ok {
		return unaggnode
	}
	comment := "this is an unscheduled node for job\n"
	node := gm.gcm.AddNode(0, NODE_UNSCEDULED, comment)
	node.jd = jd
	gm.jobUnscheds[jid] = node
	return node
}


func (gm *GraphManager) RemoveUnscehduleAggNode(jd *JobDescriptor) *GraphNode {
	if jd == nil {
		//vlog.Dlog("Error: gm.RemoveUnscehduleAggNode, jd is nil")
	}
	jid := jd.GetJobID()
	node, ok := gm.jobUnscheds[jid]
	if !ok {
		return nil
	}
	gm.gcm.RemoveNode(node)
	//vlog.Vlog(fmt.Sprintf("gcm remove unscheduled agg node %d for job %d", node.nid, jid))
	delete(gm.jobUnscheds, jid)
	return node
}

func (gm *GraphManager) AddRescourceNode(rd *ResDescriptor) *GraphNode {
	if rd == nil {
		//vlog.Dlog("Error ,gm.AddrescourceNode, td is nil")
	}
	rid := rd.GetResID()
	if rnode, ok := gm.resNodes[rid]; ok {
		return rnode
	}
	rnode := gm.gcm.AddNode(0, NODE_RESOURCE, fmt.Sprintf("Add res node for res:%d\n", rid))
	//vlog.Vlog(fmt.Sprintf("gcm AddRescourceNode, add node %d for res %d", rnode.nid, rid))
	rnode.rd = rd
	gm.resNodes[rid] = rnode
	if rd.Rtype == RES_PU {
		gm.leafNodes[rnode.nid] = true
		gm.leafResID[rid] = true
	}
	return rnode
}

func (gm *GraphManager) RemoveRescourceNode(rd *ResDescriptor) {
	if rd == nil {
		//vlog.Dlog("Error ,gm.RemoveRescourceNode, td is nil")
	}
	rid := rd.GetResID()
	rnode, ok := gm.resNodes[rid]
	if !ok {
		//vlog.Dlog("Error,gm.RemoveResourceNode rnode is not exit")
		return
	}
	delete(gm.resNodes, rid)
	//vlog.Vlog(fmt.Sprintf("gcm RemoveRescourceNode, remove res node %d,resid- %d", rnode.nid, rid))
	if rd.Rtype == RES_PU {
		delete(gm.leafNodes, rnode.nid)
		delete(gm.leafResID, rid)
	}
	gm.gcm.RemoveNode(rnode)
}

func (gm *GraphManager) UpdateUnscheduledNode(node *GraphNode, cap_delta int64) {
	if node == nil|| node.ntype != NODE_UNSCEDULED{
		//vlog.Dlog("Error: gm.updateunschedulednode,node is nill or error type\n")
	}
	//vlog.Vlog(fmt.Sprintf("gm - UpdateUnscheduledAggNode, add node(%d) to sinknode by %d", node.nid, cap_delta))
	arc := gm.gcm.GetArc(node, gm.sinkNode)
	if arc == nil {
		ad := gm.costmodel.UnschedAggToSink(node.jd)
		capupper := int64(ad.capUpper) + cap_delta
		arc = gm.gcm.AddArc(node, gm.sinkNode, ad.capLower, uint64(capupper), ad.cost, "unscheduled to sinknode\n")
	} else {
		capupper := int64(arc.capUpper) + int64(cap_delta)
		gm.gcm.ChangeArc(arc, arc.capLower, uint64(capupper), arc.cost, "change arc\n")
	}
}

func (gm *GraphManager) UpdateTaskToResourceArcs(tnode *GraphNode) {
	if tnode == nil {
		fmt.Println("Error,gm.UpdateTaskToResourceArcs, tnode is nil")
	}
	//fmt.Println("gm.UpdateTaskToResourceArcs, tnode 0",tnode.td.Name)
	preRes := gm.costmodel.TaskPreferdResource(tnode.td)
	//fmt.Println(preRes)
	for _, rd := range preRes {
		rnode := gm.ResIDToNode(rd.GetResID())
		ad := gm.costmodel.TaskNodeToResource(tnode.td, rd)
		if arc := gm.gcm.GetArc(tnode, rnode); arc != nil {
			gm.gcm.ChangeArc(arc, ad.capLower, ad.capUpper, arc.cost, "update arc for task prefer arc\n")
		} else {
			gm.gcm.AddArc(tnode, rnode, ad.capLower, ad.capUpper, ad.cost, "add arc to this resource\n")
		}
	}
}

func (gm *GraphManager) GetNode(nid NodeID) *GraphNode {
	return gm.gcm.GetNode(nid)
}

func (gm *GraphManager) GetArc(src, dst *GraphNode) *GraphArc {
	return gm.gcm.GetArc(src, dst)
}


func (gm *GraphManager) AddOrUpdateAllResNodes(rds []*ResDescriptor)map[*GraphNode]bool{
	resNodes := make(map[*GraphNode]bool)
	for _, rd := range rds {
		rnode, ok := gm.resNodes[rd.GetResID()]
		if !ok {
			rnode = gm.AddRescourceNode(rd)
		}
		rnode.rd = rd
		resNodes[rnode] = true
		//gm.UpdateResourceNode(rnode)
	}
	return resNodes
}

func (gm *GraphManager) UpdateResourceNode(rnode *GraphNode) {
	if rnode == nil {
		//vlog.Dlog("Error gm.UpdateResourceNode rnode is nil")
	}
	rd := rnode.rd
	//vlog.Vlog(fmt.Sprintf("gm-UpdateResourceNode %d", rd.GetResID()))
	gm.costmodel.UpdateResNode(rnode.rd)
	if rd.Rtype == RES_PU {
		ad := gm.costmodel.LeafRescourceToSink(rd)
		arc := gm.gcm.GetArc(rnode, gm.sinkNode)
		if arc == nil {
			gm.gcm.AddArc(rnode, gm.sinkNode, ad.capLower, ad.capUpper, ad.cost, "arc from res_pu to sink\n")
		} else {
			gm.gcm.ChangeArc(arc, ad.capLower, ad.capUpper, ad.cost, "arc from res_pu to sink\n")
		}
	}
}

func (gm *GraphManager) AddOrUpdateJobNodes(jobs []*JobDescriptor)map[*GraphNode]bool{
	//fmt.Println("in gm.addorupdatejobnodes",jobs)
	newNode := make(map[*GraphNode]bool)
	for _, jd := range jobs {
		jid := jd.GetJobID()
		var unschedulednode *GraphNode
		var ok bool
		if unschedulednode, ok = gm.jobUnscheds[jid]; !ok {
			unschedulednode = gm.AddUnscehduleAggNode(jd)
		}
		gm.UpdateUnscheduledNode(unschedulednode, 0)
		gm.AddAllTask(jd,newNode)
	}
	return newNode
}

func (gm *GraphManager) AddAllTask(jd *JobDescriptor, newNode map[*GraphNode]bool) {
	//fmt.Println(jd.Tasks)
	for _, td := range jd.Tasks {
		node := gm.AddTaskNode(td)
		if td.State == TASK_RUNNING && node.excess == 0{
			continue
		}
		newNode[node] = true
	}
}


func (gm *GraphManager) UpdateFlowGraph() {
	gm.sinkNode.excess = 0 
	newTaskNode := make(map[*GraphNode]bool)
	for _,v := range gm.taskNodes{
		newTaskNode[v]=true
	}
	rds := make([]*ResDescriptor,0,0)
	for _,rnode := range gm.resNodes{
		rds = append(rds,rnode.rd)
	}
	newResNode := gm.AddOrUpdateAllResNodes(rds)
	for rnode,_ := range newResNode{
		if rnode.IsTaskNode() {
			gm.UpdateTaskNode(rnode)
		} else if rnode.IsResourceNode() {
			gm.UpdateResourceNode(rnode)
		}
	}

	for tnode, _ := range newTaskNode {
		if tnode.IsTaskNode() {
			gm.UpdateTaskNode(tnode)
		} else if tnode.IsResourceNode() {
			gm.UpdateResourceNode(tnode)
		}
		//delete(newNode, tnode)
	}
}

func (gm *GraphManager) UpdateTaskNode(node *GraphNode) {
	if _, ok := gm.taskRuningArcs[node.td.GetTaskID()]; ok {
		gm.UpdateRunningTask(node)
		gm.costmodel.UpdateTaskNode(node.td)
	} else {
		// UpdateTaskNode mast early than UpdateTaskToResourceArcs, 
		// becouse UpdateTaskNode will Reset the Arcs
		gm.costmodel.UpdateTaskNode(node.td)
		gm.UpdateTaskToResourceArcs(node)
		//gm.costmodel.UpdateTaskNode(node.td)
		ad := gm.costmodel.TaskToUnscheduledAgg(node.td)
		unschednode := gm.JobIDToUnscheNode(node.jd.GetJobID())
		//fmt.Println(unschednode)
		//fmt.Println("in fm.update task node, unsch:",unschednode)
		if node == nil || unschednode == nil {
			fmt.Println(node)
			fmt.Println(unschednode)
		}
		arc := gm.GetArc(node,unschednode)
		if arc != nil { 
			gm.gcm.ChangeArc(arc,ad.capLower,ad.capUpper,ad.cost,"update|"+arc.comment)
		} else {
			gm.gcm.AddArc(node,unschednode,ad.capLower,ad.capUpper,ad.cost,"add arc to unscheduled node\n")
		}
		gm.UpdateUnscheduledNode(unschednode,1)
	}
}

func (gm *GraphManager) UpdateRunningTask(node *GraphNode) {
	if node == nil {
		//vlog.Dlog("Error gm.updateRunningTask ,node is nil")
		return
	}
	arc, ok := gm.taskRuningArcs[node.td.GetTaskID()]
	if ok {
		ad := gm.costmodel.TaskContinuation(node.td)
		gm.gcm.ChangeArc(arc, ad.capLower, ad.capUpper, ad.cost, "update running arc\n")
	} else {
		//vlog.Dlog("Error : gm.UpdateRunningTask. this task is not running")
	}
}

//删除节点
func (gm *GraphManager) TaskCompleted(tnode *GraphNode) {
	if tnode == nil {
		return
	}
	glog.Infof("SCHEDULE INFO:Task Completed,tnode is %v,tnode.name is ",tnode.td.Name)
	td := tnode.td
	if td == nil {
		fmt.Println("in gm.taskcompleted, error td is nil")
		return 
	}
	runningArc, ok := gm.taskRuningArcs[td.GetTaskID()]
	glog.Infof("td.name : %s, td.State :%d",td.Name,td.State)
	if ok == false && td.State == TASK_COMPLETED {
		loc := 0
	    jd := td.Jd
		loc = 0
	    for k,v := range jd.Tasks{
	        if v == td{
	            loc = k
	            break
			}
		}
		jd.Tasks = append(jd.Tasks[:loc],jd.Tasks[loc+1:]...)
	    gm.RemoveTaskHelper(tnode)
	}
	if ok == false {
		//vlog.Vlog("Error gm.taskcompleted ,this task is not running")
		return
	}
	rnode := runningArc.dst
	rd := rnode.rd
	glog.Infof("SCHEDULE_INFO : ResNode Name is:%s",rd.Name)
	rd.ResAvailable.Gpu += td.ResRequest.Gpu
	rd.ResAvailable.Cpu += td.ResRequest.Cpu
	rd.ResAvailable.Memory += td.ResRequest.Memory
	loc := 0
	tid := td.GetTaskID()
    for k,v := range rd.CurrentRunningTasks{
        if v ==tid {
            loc = k
            break
        }
    }
    rd.CurrentRunningTasks = append(rd.CurrentRunningTasks[:loc],rd.CurrentRunningTasks[loc+1:]...)	
	delete(gm.taskRuningArcs, td.GetTaskID())
	td.State = TASK_COMPLETED
	jd := td.Jd
	loc = 0
	for k,v := range jd.Tasks{
		if v == td{
			loc = k
			break
		}
	}
	jd.Tasks = append(jd.Tasks[:loc],jd.Tasks[loc+1:]...)
	gm.RemoveTaskHelper(tnode)
	glog.Infof("SCHEDULE_INFO : After TaskComPleted, tasks is :%v",gm.GetTaskInfo())
}

func (gm *GraphManager) TaskFailed(tnode *GraphNode){
	if tnode == nil {
		return
	}
	glog.Infof("SCHEDULE INFO:Task Failed,tnode is %v,tnode.name is ",tnode.td.Name)
	td := tnode.td
	if td == nil {
		return
	}
	runningArc, ok := gm.taskRuningArcs[td.GetTaskID()]
    if ok == false {
        return
    }
    rnode := runningArc.dst
    rd := rnode.rd
    rd.ResAvailable.Gpu += td.ResRequest.Gpu
    rd.ResAvailable.Cpu += td.ResRequest.Cpu
    rd.ResAvailable.Memory += td.ResRequest.Memory
    loc := 0
    tid := td.GetTaskID()
    for k,v := range rd.CurrentRunningTasks{
        if v ==tid {
            loc = k
            break
        }
    }
    rd.CurrentRunningTasks = append(rd.CurrentRunningTasks[:loc],rd.CurrentRunningTasks[loc+1:]...) 
    delete(gm.taskRuningArcs, td.GetTaskID())
    jd := td.Jd
    loc = 0
    for k,v := range jd.Tasks{
        if v == td{
            loc = k
            break
        }
    }
    jd.Tasks = append(jd.Tasks[:loc],jd.Tasks[loc+1:]...)
	td.State = TASK_FAILED
    gm.RemoveTaskHelper(tnode)
	glog.Infof("SCHEDULE_INFO : After TaskComPleted, tasks is :%v",gm.GetTaskInfo())
}

func (gm *GraphManager) TaskScheduled(tnode, rnode *GraphNode) {
	//vlog.Vlog(fmt.Sprintf("schedule ,%d to %d", tnode.nid, rnode.nid))
	if tnode.td.State == TASK_RUNNING {
		gm.UpdateArcsForScheduledTask(tnode, rnode)
		return
	}
	tnode.td.State = TASK_RUNNING
	gm.UpdateArcsForScheduledTask(tnode, rnode)
	rd := rnode.rd
	td := tnode.td
	delete(gm.queue,tnode)
	//tnode.excess = 0
	rd.ResAvailable.Gpu -= td.ResRequest.Gpu
	rd.ResAvailable.Cpu -= td.ResRequest.Cpu
	rd.ResAvailable.Memory -= td.ResRequest.Memory
	if rnode.rd.ResAvailable.Gpu > 8 {
		glog.Infof("SCHEDULE_INFO : ERROR,Rnode.GPU>8,info: rnode.Name %s",rnode.rd.Name)
		glog.Infof("SCHEDULE_INFO : rd.Resavai: %d, rd.ResPreet:%d",rnode.rd.ResAvailable.Gpu,rnode.rd.ResPreempt.Gpu)
    }

}


func (gm *GraphManager)HandleTaskKilled(tid TaskID){
	gm.TaskFailed(gm.TaskIDToNode(tid))
/*	tnode := gm.TaskIDToNode(tid)
	td := tnode.td
	td.State = TASK_FAILED
	arc := gm.taskRuningArcs[tid]
	fmt.Println(arc)
	//fmt.Println(arc.dst)
	rd := arc.dst.rd
	rd.ResAvailable.Gpu += td.ResRequest.Gpu
	loc := 0
	for k,v := range rd.CurrentRunningTasks{
		if v ==tid {
			loc = k
			break
		}
	}
	rd.CurrentRunningTasks = append(rd.CurrentRunningTasks[:loc],rd.CurrentRunningTasks[loc+1:]...)
	
	delete(gm.taskRuningArcs,tid)
*/
}

func (gm *GraphManager) UpdateArcsForScheduledTask(tnode, rnode *GraphNode) {
	if tnode == nil || rnode == nil {
		//vlog.Dlog("Error: gm.updateArcsforScheduledTask tnode or rnode is nil")
		return
	}
	tid := tnode.td.GetTaskID()
	//fmt.Println("in fm.UpdateArcsFor...,taskrunningarcs:",gm.taskRuningArcs)
	if arc, ok := gm.taskRuningArcs[tid]; ok {
		arc.atype = ARC_RUNNING
		arcDes := gm.costmodel.TaskContinuation(tnode.td)
		gm.gcm.ChangeArc(arc, arcDes.capLower, arcDes.capUpper, arcDes.cost, "update this arc\n")
	} else {
		gm.PinTaskToResNode(tnode, rnode)
	}
}

func (gm *GraphManager) PinTaskToResNode(tnode, rnode *GraphNode) {
	//fmt.Println("in gm.PinTaskToresnode,tid,rid = ",tnode.td.Name,rnode.rd.Uuid)
	if tnode == nil || rnode == nil {
		//vlog.Dlog("error gm.pintasktoresnode,rnode or tnode is mil")
		return
	}
	//vlog.Vlog(fmt.Sprintf("gm-PinTaskToResNode, task node /%d to res node/%d", tnode.nid, rnode.nid))
	for dstID, arc := range tnode.outgoingArcs {
		if dstID == rnode.nid {
			arc.atype = ARC_RUNNING
			arcDes := gm.costmodel.TaskContinuation(tnode.td)
			gm.gcm.ChangeArc(arc, arcDes.capLower, arcDes.capUpper, arcDes.cost, "turn this arv to running arc!\n")
			gm.taskRuningArcs[tnode.td.GetTaskID()] = arc
		} else if arc.dst.ntype == NODE_RESOURCE{
			gm.gcm.RemoveArc(arc)
		}
	}
	tnode.td.State = TASK_RUNNING
}

func (gm *GraphManager) RemoveTaskHelper(tnode *GraphNode) {
	if tnode == nil {
		return
	}
	glog.Infof("SCHEDULE_INFO : remove task node , tname is :%v",tnode.td.Name)
	delete(gm.taskRuningArcs, tnode.td.GetTaskID())
	gm.RemoveTaskNode(tnode.td)
}

func (gm *GraphManager) JobCompleted(jd *JobDescriptor) {
	//remove unscheduled node |remove agg node
	if jd == nil {
		//vlog.Dlog("Error gm.Jobcompleted, jd is nil")
		return
	}
	jid := jd.GetJobID()
	unschedulednode, ok := gm.jobUnscheds[jid]
	if ok == false {
		//vlog.Dlog("Error gm.Jobcompleted, unscheduled is nil")
		return 
	}
	gm.gcm.RemoveNode(unschedulednode)
	delete(gm.jobUnscheds, jid)
}


func (gm *GraphManager) ExportGraph() string {
	return gm.gcm.fg.ExportGraph()
}

func (gm *GraphManager)GetUnExportNodes()map[NodeID]bool{
	ans := make(map[NodeID]bool)
	for _,runningArc := range gm.taskRuningArcs {
		tnode := runningArc.src
		jd := tnode.jd
		ans[tnode.nid] = true
		jobUnsched := gm.JobIDToUnscheNode(jd.GetJobID())
		ans[jobUnsched.nid] = true
	} 
	return ans
}

/// this func didn't handle : a job has more than two tasks
func (gm *GraphManager)ExportGraphWithoutRunningTasks()string{
	unExportNodes := gm.GetUnExportNodes()
	pre := "c This is a max-flow min-cost problem\n"
	var nodeNum,arcNum int64
	ans := ""
	ans += "c Nodes\n"
	g := gm.gcm.fg
	for _,node := range g.nodes{
		if _,ok := unExportNodes[node.nid];ok{
			continue
		}
		ans += node.Export()
		nodeNum++
	}
	ans += "c arcs\n"
	for arc,_ := range g.arcSet{
		srcid,dstid := arc.srcID,arc.dstID
		if _,ok := unExportNodes[srcid];ok{
			continue
		}
		if _,ok := unExportNodes[dstid];ok{
			continue
		}
		arcNum++
		ans += arc.Export()
	}
	pre += fmt.Sprintf("p min %d %d\n",nodeNum,arcNum)
	//vlog.Dlog(fmt.Sprintf("un export node nums : %d",len(unExportNodes)))
	//vlog.Dlog(fmt.Sprintf("node num : %d arc num:%d",nodeNum,arcNum))
	return pre+ans
}
//Taskfailed
//Jobfailed


/*
func (gm *FlowGraphManager)NodeBindingToSchedulingDeltas(
								tnid NodeID,rnid NodeID,
								taskBindings map[TaskID]ResID,
								)*SchedulingDelta{
	tnode := gm.GetNodeByID(tnid)
	rnode := gm.GetNodeByID(rnid)
	if !tnode.IsTaskNode() || rnode.nodeType == NODE_PU {
		return nil
	}
	td := tnode.taskDes
	resID,ok := taskBindings[td.GetTaskID()]
	var sd *SchedulingDelta
	if ok {
		if resID != rnode.GetResID(){
			sd = NewSchedulingDelta(td.GetTaskID(),rnode.GetResID())
			sd.deltaType = DELTA_MIGRATE
			return sd
		} else {
			rnode.resDes.AddCurrentRunningTask(td.GetTaskID())
		}
	} else {
		sd = NewSchedulingDelta(td.GetTaskID(),resID)
		sd.deltaType = DELTA_PLACE
		return sd
	}
	return nil
}



//failed/migarated 默认已经在运行中
func (gm *FlowGraphManager)TaskRemoved(tid TaskID) {
	LOG(fmt.Sprintf("gm-TaskRemoved, task id/%d",tid))
	if _,isRunning := gm.taskToRuningArc[tid];isRunning {
		tnode := gm.taskToNodeMap[tid]
		jid := tnode.taskDes.GetJobID()
		jnode := gm.jobToUnschedMap[jid]
		gm.UpdateUnscheduledAggNode(jnode,-1)
	}
	gm.RemoveTaskHelper(tid)
}

*/
//*-------------
