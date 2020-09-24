package flowscheduler

//import "vlog"
import "time"
import "sort"
import "fmt"
import "math/rand"

const (
	PRIORITYCOST = uint64(1<<12)
	GPUNEEDCOST  = uint64(1<<6)
	PREEMPTCOST  = uint64(1<<18)
	NOTAVICOST   = uint64(1<<24)
	NOTSCHECOST  = uint64(1<<24)
)

type PengChengCostModel struct{
	gm *GraphManager
	
	Gpunum uint64
	Priority uint64
	Tinfo TaskDetail
	// cost will add a random value n, where 0<= n <= max(minRandLim, task Node Number)
	minRandLim  uint64
	// no more than maxschedulenum tasks can be scheduled to a resource node in every time schedule,it do not work in current settings
	MaxScheduleNum uint64
	seed int64
}

func NewPengChengCostModel(gm_ *GraphManager)*PengChengCostModel{
	//vlog.Vlog("New Gpu Cost model")
	return &PengChengCostModel{
		gm:gm_,
		Priority:1,
		Gpunum:1,
		minRandLim: 30,
		MaxScheduleNum : 3,
		seed : time.Now().UnixNano(),
		Tinfo : TaskDetail{0,0,0,-1,nil},
	}
}

func (pccm *PengChengCostModel)SetGpunum(gpunum uint64){
	pccm.Gpunum = gpunum
}

func (pccm *PengChengCostModel)SetMaxScheduleNum(msn uint64){
	pccm.MaxScheduleNum = msn
}

func (pccm *PengChengCostModel)SetminRandLim(mrl uint64){
	pccm.minRandLim = mrl
}

func (pccm *PengChengCostModel)SetPriority(p uint64){
	pccm.Priority = p
	fmt.Sprintf("")
}



func (pccm *PengChengCostModel)SetTinfo(detail TaskDetail){
	pccm.Tinfo = detail
}

func (pccm *PengChengCostModel)TaskToResCost(td *TaskDescriptor,rd *ResDescriptor)uint64{
	// resource need cost
	rnct :=  uint64(0)
	tg,rg1,rg2 := td.ResRequest.Gpu,rd.ResAvailable.Gpu,rd.ResPreempt.Gpu
	if tg > rg1+rg2 {
		rnct = NOTAVICOST
	} else {
		p := uint64(1)
		q,r := uint64(0),uint64(0)
		if tg > rg1 {
			p = 4
			q,r = uint64((rg1+rg2)/tg),uint64((rg1+rg2)%tg)
		} else {
			q,r = uint64(rg1/tg),uint64(rg1%tg)
		}
		if 1 + rd.ResAvailable.Gpu == 0 {
			rnct =  p*(q*2 +r)*GPUNEEDCOST
		} else {
			rnct = p*(q*2 +r)*GPUNEEDCOST/(1+rd.ResAvailable.Gpu)
		}
	}
	// priority cost
	/*pct := PRIORITYCOST
	fmt.Println(pccm.Priority,td.Priority)
	pct = pct*uint64(pccm.Priority - td.Priority)
	*/
	// random cost
	rand.Seed(pccm.seed)
	pccm.seed = rand.Int63()
	mrl := uint64(len(pccm.gm.taskNodes))
	if mrl < pccm.minRandLim {
		mrl = pccm.minRandLim
	}
	rct := uint64(rand.Intn(int(mrl)))

	if rnct + rct > 1e8 {
		fmt.Println("++++++++++++++++",rnct,rct)
	}
	
	return rnct + rct
}

func (pccm *PengChengCostModel)LeafRescourceToSink(rd *ResDescriptor)*ArcDescriptor{
	lim := uint64((rd.ResAvailable.Gpu+rd.ResPreempt.Gpu)/pccm.Gpunum)
	//num :=uint64(len(rd.CurrentRunningTasks))
	cu := lim
	return &ArcDescriptor{
		capUpper:cu,
		capLower:0,
		cost : 0,//+++
	}
}

func (pccm *PengChengCostModel)TaskContinuation(td *TaskDescriptor) *ArcDescriptor{
	ct := uint64(0)
	cl,cu := uint64(0),uint64(0)
	if td.State == TASK_RUNNING {
		if td.IsGangSchedule {
			cl,cu = 1,1
		} else if pccm.Priority > td.Priority{
			cl,cu = 0,1
		} else {
			cl,cu = 1,1
		}
	} else{
		fmt.Println("Error! Task is not running!")
		cl,cu = 0,1
	}
	ct = uint64(64-td.Priority)
	
	return &ArcDescriptor{
		capLower:cl,
		capUpper:cu,
		cost:ct,	
	}
}

// this arc will be updated during updating job
func (pccm *PengChengCostModel)UnschedAggToSink(jd *JobDescriptor) *ArcDescriptor {
	return &ArcDescriptor{
		capLower:0,
		capUpper:uint64(len(jd.Tasks)),
		cost:0,
	}
}

func (pccm *PengChengCostModel)TaskNodeToResource(td *TaskDescriptor,rd *ResDescriptor)*ArcDescriptor{
	cl,cu := uint64(0),uint64(1)
	if (td.State == TASK_RUNNING) && (td.IsGangSchedule  || td.Priority >= pccm.Priority){
		cl = 1
	}
	//fmt.Println("-------pccm tasknode to resource")
	return &ArcDescriptor{
		cost: pccm.TaskToResCost(td,rd),
		capLower :cl,
		capUpper: cu,
	}
}

func (pccm *PengChengCostModel)TaskPreferdResource(td *TaskDescriptor) []*ResDescriptor{
	if td.Priority != pccm.Priority  || td.ResRequest.Gpu != pccm.Gpunum{
		return make([]*ResDescriptor,0,0)
	}
	ans := make([]*ResDescriptor,0,0)
	for _,rnode := range pccm.gm.resNodes{
		rd := rnode.rd
		if rd.ResAvailable.Gpu + rd.ResPreempt.Gpu >= td.ResRequest.Gpu {
			ans = append(ans,rnode.rd)
		}
	}
	return ans
}


func (pccm *PengChengCostModel)TaskToUnscheduledAgg(td *TaskDescriptor)*ArcDescriptor{
	timecost := time.Now().UnixNano() - td.SubmitTime
	// time past cost 
	tct := timecost/1000000000
	
	// preempt cost
	pct := uint64(0)
	if td.State == TASK_RUNNING && td.IsGangSchedule {
		pct =  NOTSCHECOST
	} else if td.State == TASK_RUNNING && !td.IsGangSchedule{
		pct = (PREEMPTCOST*(pccm.Priority - td.Priority))
	} else if td.State != TASK_RUNNING && td.IsGangSchedule{
		pct = td.Priority*PRIORITYCOST
	} else if td.State != TASK_RUNNING && !td.IsGangSchedule{
		if td.Priority > 20{
			pct = uint64(1<<20)
		} else {
			pct = uint64(1<<td.Priority)
		}
	}
	pct += PRIORITYCOST
	rct := uint64(0)//td.ResRequest.Gpu*PREEMPTCOST

	return &ArcDescriptor{
		capLower:0,
		capUpper:1,
		cost : uint64(tct)+pct+rct,
	}
}

func (pccm *PengChengCostModel)UpdateTaskNode(td *TaskDescriptor){
	if td.State != TASK_RUNNING &&pccm.Priority == td.Priority {
		tid := td.GetTaskID()
		tnode := pccm.gm.TaskIDToNode(tid)
		for _,arc := range tnode.outgoingArcs{
			if arc.dst.IsResourceNode(){
				pccm.gm.gcm.ChangeArc(arc,0,0,arc.cost,"update tasknode in pccm")
			}
		}
	}

	node := pccm.gm.TaskIDToNode(td.GetTaskID())
	if (td.State == TASK_UNSCHEDULED && td.Priority == pccm.Priority && td.ResRequest.Gpu == pccm.Gpunum){
		node.excess = 1
		pccm.gm.sinkNode.excess --
	} else if td.State == TASK_RUNNING && td.Priority < pccm.Priority {
		node.excess = 0
	} else {
		node .excess = 0
	}
	return 
}


func (pccm *PengChengCostModel)UpdateResNode(rd *ResDescriptor){
	res := uint64(0)
	curd := pccm.Tinfo
	for _,v := range rd.CurrentRunningTasks{
		tnode := pccm.gm.TaskIDToNode(v)
		td := tnode.td
		d := TaskDetail{td.Priority, td.ResRequest.Gpu,rd.ResAvailable.Gpu,td.StartTime,td}
		if Compare(d,curd) {
			res += td.ResRequest.Gpu
		}
	}
	rd.ResPreempt.Gpu = res
	return 
}

func (pccm *PengChengCostModel)GetTaskKilled(temTaskMapping map[NodeID]NodeID)map[TaskID]ResID{
	res2Task := make(map[NodeID][]NodeID)
	taskUnshedule := make(map[TaskID]ResID)
	for k,v := range temTaskMapping{
		if _,ok := res2Task[v];ok{
			res2Task[v] = append(res2Task[v],k)
		} else {
			tasks := make([]NodeID,0,0)
			res2Task[v] = append(tasks,k)
		}
	}

	//fmt.Println(res2Task)
	
	for rnodeid,tasks := range res2Task {
		rnode := pccm.gm.gcm.GetNode(rnodeid)
		tds := Details{}
		for _,tnid := range tasks{
			tnode := pccm.gm.gcm.GetNode(tnid)
			//fmt.Println(tnode)
			td := tnode.td
			tds = append(tds,TaskDetail{td.Priority,td.ResRequest.Gpu,0,td.SubmitTime,td})
		}
		for _,tid := range rnode.rd.CurrentRunningTasks{
			tnode := pccm.gm.TaskIDToNode(tid)
			td := tnode.td
			temtdl :=  TaskDetail{td.Priority,td.ResRequest.Gpu,rnode.rd.ResAvailable.Gpu,td.StartTime,td}
			if Compare(temtdl,pccm.Tinfo){
				tds = append(tds,temtdl)
			}
		}
		sort.Sort(tds)
		total := rnode.rd.ResAvailable.Gpu + rnode.rd.ResPreempt.Gpu
		for i:=len(tds)-1;i>=0;i-- {
			tdl := tds[i]
			if total >= tdl.gpu{
				total -= tdl.gpu
			} else{
				taskUnshedule[tdl.td.GetTaskID()]= rnode.rd.GetResID()
			}
		}
	}
	return taskUnshedule
}
	





































/*

func (pccm *PengChengCostModel)UpdateTaskNode(td *TaskDescriptor){
	if td.State != TASK_RUNNING &&pccm.Priority == td.Priority {
		tid := td.GetTaskID()
		tnode := pccm.gm.TaskIDToNode(tid)
		for _,arc := range tnode.outgoingArcs{
			if arc.dst.IsResourceNode(){
				pccm.gm.gcm.ChangeArc(arc,0,0,arc.cost,"update tasknode in pccm")
			}
		}
	}

	node := pccm.gm.TaskIDToNode(td.GetTaskID())
	if (td.State != TASK_RUNNING && td.Priority == pccm.Priority && td.ResRequest.Gpu == pccm.Gpunum){
		node.excess = 1
		pccm.gm.sinkNode.excess --
	} else if td.State == TASK_RUNNING && td.Priority < pccm.Priority {
		num := int(td.ResRequest.Gpu)
		if num != 1 {
			for i:=0;i<num ;i++{
				tn := pccm.gm.gcm.AddNode(1,NODE_TASK,"temNode for task"+td.Name)
				tn.td = new(TaskDescriptor)
				tn.td.Tid = td.Tid
				tn.td.ResRequest = ResVector{0,1,0}
				pccm.temNode = append(pccm.temNode,tn)
				ad1 := pccm.TaskToUnscheduledAgg(td)
				ad2 := pccm.TaskContinuation(td)
				usn := pccm.gm.JobIDToUnscheNode(td.GetJobID())
				pccm.gm.gcm.AddArc(tn,usn,ad1.capLower,ad1.capUpper,ad1.cost,"tem node to unschedulenode arc")
				arc,ok := pccm.gm.taskRuningArcs[td.GetTaskID()]
				if !ok {
					fmt.Println("error task is not running,in pccm.updatetasknode.")
					return 
				}
				rn := arc.dst
				pccm.gm.gcm.AddArc(tn,rn,ad2.capLower,ad2.capUpper,ad2.cost,"tem node "+td.Name+" to res node arc")
			}
			node.excess = 0
			pccm.gm.sinkNode.excess -= int64(td.ResRequest.Gpu)
		} else {
			node.excess = 1
			pccm.gm.sinkNode.excess --
		}
	} else {
		node .excess = 0
	}
	return 
}




*/
