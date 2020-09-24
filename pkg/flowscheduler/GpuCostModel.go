package flowscheduler

//import "vlog"
import "time"
import "fmt"
//import "math/rand"

type GPUCostModel struct{
	gm *GraphManager
	//gpuNum is the number of gpu need in cur round schedule. 
	gpuNum uint64
}


func NewGpuCostModel(gm_ *GraphManager)*GPUCostModel{
	//vlog.Vlog("New Gpu Cost model")
	return &GPUCostModel{
		gm:gm_,
		gpuNum :1,
	}
}

func (gcs *GPUCostModel)SetGpuNum(num uint64){
	gcs.gpuNum = num
	fmt.Sprintf("")
}

func (gcs *GPUCostModel)TaskToResCost(td *TaskDescriptor,rd *ResDescriptor)uint64{
	if gcs.gpuNum == 0 {
		return 0
	}
	r := rd.ResAvailable.Gpu % gcs.gpuNum
	q := rd.ResAvailable.Gpu / gcs.gpuNum
	return 5*r+20*q+r*q
}


func (gcs *GPUCostModel)LeafRescourceToSink(rd *ResDescriptor)*ArcDescriptor{
	return &ArcDescriptor{
		capUpper:rd.ResAvailable.Gpu/gcs.gpuNum,
		capLower:0,
		cost : 0,//+++
	}
}


func (gcs *GPUCostModel)TaskContinuation(td *TaskDescriptor) *ArcDescriptor{
	return &ArcDescriptor{
		capLower:0,
		capUpper:1,
		cost:(1<<28),	
	}
}

// this arc will be updated during updating job
func (gcs *GPUCostModel)UnschedAggToSink(jd *JobDescriptor) *ArcDescriptor {
	return &ArcDescriptor{
		capLower:0,
		capUpper:0,
		cost:50,
	}
}

//arc from task to res.
func (gcs *GPUCostModel)TaskNodeToResource(td *TaskDescriptor,rd *ResDescriptor)*ArcDescriptor{
	//ct := uint64(rand.Intn(100))
	return &ArcDescriptor{
		cost: gcs.TaskToResCost(td,rd),
		capLower :0,
		capUpper: 1,
	}
}


func (gcs *GPUCostModel)TaskPreferdResource(td *TaskDescriptor) []*ResDescriptor{
	ans := make([]*ResDescriptor,0,0)
	for _,rnode := range gcs.gm.resNodes{
		rd := rnode.rd
		if rd.ResAvailable.Gpu >= gcs.gpuNum{
			ans = append(ans,rnode.rd)	
		}
	}
	return ans
}


// Feb 18, Changed
func GetTaskNumberWithGPULimit(jd *JobDescriptor,num uint64)uint64{
	var ans uint64
	for _,td := range jd.Tasks{
		if td.ResRequest.Gpu == num  && td.State != TASK_RUNNING{
			ans ++
		}
	}
	return ans 
}


func (gcs *GPUCostModel)TaskToUnscheduledAgg(td *TaskDescriptor)*ArcDescriptor{
	timecost := time.Now().UnixNano() - td.SubmitTime
	return &ArcDescriptor{
		capLower:0,
		capUpper:1,
		cost:uint64(timecost/1000),
	}
}

func (gcs *GPUCostModel)UpdateTaskNode(td *TaskDescriptor){
	node := gcs.gm.TaskIDToNode(td.GetTaskID())
	if td.ResRequest.Gpu == gcs.gpuNum && td.State != TASK_RUNNING{
		node.excess = 1
		gcs.gm.sinkNode.excess --
	} else {
		node.excess = 0
	}
}

func (gcs *GPUCostModel)UpdateResNode(rd *ResDescriptor){
	return

}
/*
// 非root task：添加root task to task的边
func (gcs *GPUCostModel)56(td *TaskDescriptor){
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