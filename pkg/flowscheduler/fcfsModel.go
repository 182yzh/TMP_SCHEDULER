
package flowscheduler

//import "vlog"
import "time"
import "fmt"
//import "math/rand"

type FCFSModel struct{
	gm *GraphManager
	//gpuNum is the number of gpu need in cur round schedule. 
	limit uint64
}

func (fm *FCFSModel)SetLimit(t uint64){
	fm.limit = t
}

func NewFCFSCostModel(gm_ *GraphManager)*FCFSModel{
	//vlog.Vlog("New fcfs Cost model")
	fmt.Sprintf("this is a new ---\n")
	return &FCFSModel{
		gm:gm_,
		limit:4,
	}
}


func (fm *FCFSModel)TaskToResCost(td *TaskDescriptor,rd *ResDescriptor)uint64{
	return td.ResRequest.Gpu*10
}

func (fm *FCFSModel)LeafRescourceToSink(rd *ResDescriptor)*ArcDescriptor{
	mn := fm.limit
	if rd.ResAvailable.Gpu < fm.limit{
		mn = rd.ResAvailable.Gpu
	}
	return &ArcDescriptor{
		capUpper:mn,
		capLower:0,
		cost : 0,//+++
	}
}


func (fm *FCFSModel)TaskContinuation(td *TaskDescriptor) *ArcDescriptor{
	return &ArcDescriptor{
		capLower:0,
		capUpper:1,
		cost:(1<<28),	
	}
}

// this arc will be updated during updating job
func (fm *FCFSModel)UnschedAggToSink(jd *JobDescriptor) *ArcDescriptor {
	return &ArcDescriptor{
		capLower:0,
		capUpper:0,
		cost:50,
	}
}

//arc from task to res.
func (fm *FCFSModel)TaskNodeToResource(td *TaskDescriptor,rd *ResDescriptor)*ArcDescriptor{
	//ct := uint64(rand.Intn(100))
	return &ArcDescriptor{
		cost: fm.TaskToResCost(td,rd),
		capLower :0,
		capUpper: 1,
	}
}


func (fm *FCFSModel)TaskPreferdResource(td *TaskDescriptor) []*ResDescriptor{
	ans := make([]*ResDescriptor,0,0)
	for _,rnode := range fm.gm.resNodes{
		rd := rnode.rd
		if rd.ResAvailable.Gpu >= td.ResRequest.Gpu{
			ans = append(ans,rnode.rd)	
		}
	}
	return ans
}




func (fm *FCFSModel)TaskToUnscheduledAgg(td *TaskDescriptor)*ArcDescriptor{
	timecost := time.Now().UnixNano() - td.SubmitTime
	return &ArcDescriptor{
		capLower:0,
		capUpper:1,
		cost:uint64(timecost/1000),
	}
}

func (fm *FCFSModel)UpdateTaskNode(td *TaskDescriptor){
	node := fm.gm.TaskIDToNode(td.GetTaskID())
	if td.State != TASK_RUNNING{
		node.excess = 1
		fm.gm.sinkNode.excess --
	} else {
		node.excess = 0
	}
}

func (fm *FCFSModel)UpdateResNode(rd *ResDescriptor){
	return 
}
/*
// 非root task：添加root task to task的边
func (fm *FCFSModel)56(td *TaskDescriptor){
	if td.taskType == TASK_ROOT{
		gcs.TaskToUnscheduledAgg(td)
		return 
	}
	LOG("Find or not++++++")
	//add or update the arc from root task(job agg node) to task
	gm := gcs.gm
	tnode := gm.TaskIDToNode(td.GetTaskID())
	jd := td.jobDes
	rtnode := gm.TaskIDToNode(jd.rootTask.GetTaskID())
	arc := gcs.gm.gcm.GetArc(rtnode,tnode)
	var cap uint64 = 1
	if td.resourceRequest.gpu != gcs.gpuNum{
		cap = 0
	}
	var cost uint64 = 0 //+++ cost should not be zero;
	if arc == nil {
		gcs.gm.gcm.AddArc(rtnode,tnode,0,cap,cost,"add arc from rtnode to tnode\n")
	} else {
		gcs.gm.gcm.ChangeArc(arc,0,cap,cost,"change arc from rtnode to tnode\n")
	}
}
*/
