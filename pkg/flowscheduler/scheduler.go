package flowscheduler

import "fmt"
import "time"
//import "vlog"
import "sort"
import "github.com/golang/glog"
import "sync"

var jobmutex sync.Mutex
//var taskmutex sync.Mutex


type FlowScheduler struct{
	gm      *GraphManager
	dp      *Dispatcher
	cm      CostModel
	// the pus that was removed during the slover run
	//pusRemoved       map[NodeID]bool
	//taskCompleted    map[NodeID]bool
	jobs      		 map[JobID]*JobDescriptor
	jobuids          map[string]JobID
	toScheduleJobs   []*JobDescriptor
	tasks            map[TaskID]*TaskDescriptor
	ress             map[ResID]*ResDescriptor
	//taskBindings    map[TaskID]ResID
	//resBindings     map[ResID](map[TaskID]bool)
	runnableTask     map[JobID](map[TaskID]bool)
	lastUpdateTime   uint64
	//if time has passed updateFrequncy(ms) ,then we update time based cost
	UpdateFrequency  uint64
	sloverRunCnt     uint64
}


func NewFlowScheduler()*FlowScheduler{
	fs := new(FlowScheduler)
	fs.gm = new(GraphManager)
	fs.gm.Init()
	fs.dp = NewDispatcher(fs.gm)
	//fs.cm = NewGpuCostModel(fs.gm)
	//fs.cm = NewFCFSCostModel(fs.gm)
	//fs.pusRemoved = make(map[NodeID]bool)
	//fs.taskCompleted = make(map[NodeID]bool)
	fs.cm = NewPengChengCostModel(fs.gm)
	fs.jobs = make(map[JobID]*JobDescriptor)
	fs.jobuids = make(map[string]JobID)
	fs.tasks = make(map[TaskID]*TaskDescriptor)
	fs.ress =make(map[ResID]*ResDescriptor)
	//fs.taskBindings = make(map[TaskID]ResID)
	//fs.resBindings = make(map[ResID](map[TaskID]bool))
	fs.runnableTask = make(map[JobID](map[TaskID]bool))
	fs.toScheduleJobs = make([]*JobDescriptor,0,0)
	return fs
}

func (fs *FlowScheduler)UsedResID(rid uint64)bool{
    _,ok := fs.ress[rid]
    return ok
}

func (fs *FlowScheduler)UsedJobID(jid uint64)bool{
    _,ok := fs.jobs[jid]
    return ok
}

func (fs *FlowScheduler)ChangeJobDes(jd *JobDescriptor)bool{
	jobmutex.Lock()
	defer jobmutex.Unlock()
	//taskmutex.Lock()
	//defer taskmutex.Unlock()
	glog.Infof("SCHEDULE INFO: JobChange: %v",jd)
	if jid,ok := fs.jobuids[jd.Uuid];ok{
		changejd,ok := fs.jobs[jid]
		if ok == false {
			fmt.Println("in fs.cahngehobdes Error\n")
		}
		changejd.Tasks = append(changejd.Tasks,jd.Tasks...)
		changejd.GangScheduleNum = len(changejd.Tasks)
		for _,td := range changejd.Tasks {
			td.IsGangSchedule = true
			td.Jd = changejd
			if td.ResRequest.Gpu != 1 {
				tmp := float64(td.ResRequest.Gpu) * 1.51
				rp := int64(tmp+0.5)
				td.Priority =uint64( int(rp) + len(changejd.Tasks))
			} else {
				td.Priority = uint64(len(changejd.Tasks)) + 2
			}
			fs.tasks[td.GetTaskID()] = td
			fs.runnableTask[jid][td.GetTaskID()] = true
		}
		fs.toScheduleJobs = append(fs.toScheduleJobs,changejd)
		return true
	}
	return false
}



func (fs *FlowScheduler) HandleTaskCompeletedByID(tid TaskID){
    glog.Infof("SCHEDULE INFO:COMPELETEByID tid is %v",tid)
	//taskmutex.Lock()
	//defer taskmutex.Unlock()
	td,_ := fs.tasks[tid]
	glog.Infof("SCHEDULE_INFO:TASKCOMPELETED td is %s",fmt.Sprintln(td))
	if td == nil {
		return
	}
	glog.Infof("SCHEDULE INFO:TASK COMPLETE %s| %v",td.Name,td.GetTaskID())
	jd := td.Jd
	for _,task := range jd.Tasks{
		task.State = TASK_COMPLETED
	}
    fs.HandleTaskCompeleted(td)
}

func (fs *FlowScheduler) HandleTaskFailedByID(tid TaskID){
	glog.Infof("SCHEDULE INFO:FAILEDByID tid is %v",tid)	
	//taskmutex.Lock()
//	defer taskmutex.Unlock()
	td,_ := fs.tasks[tid]
	if td == nil {
		return
	}
	if td.State == TASK_COMPLETED {
		//taskmutex.Unlock()
		fs.HandleTaskCompeletedByID(tid)
		return
	}
	//defer taskmutex.Unlock()
	if td == nil || td.State != TASK_FAILED{
		return 
	}
	glog.Infof("SCHEDULE INFO:TASK FAILED %s|%d",td.Name,td.GetTaskID())
	fs.HandleTaskFailed(td)
}


func (fs *FlowScheduler)ScheduleAllJobsToPoseidon()(map[TaskID]ResID,map[TaskID]ResID){
	jobmutex.Lock()
	defer jobmutex.Unlock()
    taskToRes, taskKilled := fs.ScheduleJobs(fs.toScheduleJobs)
    fs.toScheduleJobs = make([]*JobDescriptor,0,0)
//	ans := make(map[TaskID]string)
    //taskKilled := make(map[TaskID]string)
	/*for tid,rid := range taskToRes{
		tnode := fs.gm.TaskIDToNode(tid)
		td := tnode.td
		rd := fs.ResIDToResDes(rid)
		for _,trid := range rd.CurrentRunningTasks{
			rtd := fs.gm.TaskIDToNode(trid).td
			if rtd.Priority > td.Priority {
				continue
			}
		}
	}*/
//	for tid,rid := range taskToRes{
		//tnode := fs.gm.TaskIDToNode(tid)
        //if _,ok := fs.gm.taskRuningArcs[tnode.td.GetTaskID()];ok {
        //    continue
        //}
  //      rd := fs.ResIDToResDes(rid)
    //    ans[tid] = rd.Uuid
//	}
	//`fs.ReplaceKilledTask()
	tmp := make(map[string]uint64)
	for _,rd := range fs.ress {
		tmp[rd.Name]=rd.ResAvailable.Gpu
	}
    return taskToRes,taskKilled
}

func (fs *FlowScheduler)GetResInfo()map[string]uint64{
	return fs.gm.GetResInfo()
}

func (fs *FlowScheduler)GetTaskInfo()map[string]string{
	return fs.gm.GetTaskInfo()
}

func (fs *FlowScheduler)ScheduleAllJobs()(map[TaskID]ResID,map[TaskID]ResID){
	taskToRes,taskKilled := fs.ScheduleJobs(fs.toScheduleJobs)
	fs.toScheduleJobs = make([]*JobDescriptor,0,0)
	return taskToRes,taskKilled
}
func (fs *FlowScheduler)TaskIDToTaskDes(tid TaskID)*TaskDescriptor{
	td,ok := fs.tasks[tid]
    if ok == false {
        return nil
    }
    return td
}


func (fs *FlowScheduler)TaskIDToJobDes(tid TaskID)*JobDescriptor{
	td,ok := fs.tasks[tid]
	if ok == false {
		//vlog.Vlog("fs.TaskIDToJobDes, td is nil")
		return nil
	}
	return td.Jd
}

func (fs *FlowScheduler)ResIDToResDes(rid ResID)*ResDescriptor{
	rd,ok := fs.ress[rid]
	if !ok {
		return nil
	}
	return rd
}
func (fs *FlowScheduler)ScheduleJobs(jobs []*JobDescriptor)(map[TaskID]ResID,map[TaskID]ResID){
	fs.gm.AddOrUpdateJobNodes(jobs)
	//fs.gm.UpdateFlowGraph()
	glog.Infof("SCHEDULE INFO:TASK_INFO  %v",fs.GetTaskInfo())
	glog.Infof("SCHEDULE INFO:RES_INFO  %v",fs.GetResInfo())
	return fs.RunScheduleIteration()
}




func (fs *FlowScheduler)HandleTaskCompeleted(td *TaskDescriptor){
	//if task is abort or failed ,is should be already removed
	if td == nil {
		return 
	}
	tid := td.GetTaskID()
	glog.Infof("SCHEDULE INFO: Task Completed,td.name is :",td.Name)
	if _,ok := fs.tasks[tid];!ok{
		glog.Infof("fs.tasks do not have the task:%s",td.Name)
		return 
	}
	if td.Name != fs.tasks[tid].Name {
		glog.Infof(td.Name)	
		return 
	}
	delete(fs.tasks,tid)
	jid := td.Jd.GetJobID()
	delete(fs.runnableTask[jid],tid)
	//if rid,ok := fs.taskBindings[tid];ok{
		//delete(fs.taskBindings,tid)
		//delete(fs.resBindings[rid],tid)
	//}
	td.State = TASK_COMPLETED
	tnode := fs.gm.TaskIDToNode(td.GetTaskID())
	if tnode == nil {
		glog.Infof("SCHEDULE_INFO:tnode(%s) is nil",td.Name);
	}
	fs.gm.TaskCompleted(tnode)
	fs.CheckJobCompleted(td.Jd)
}


func (fs *FlowScheduler)HandleTaskFailed(td *TaskDescriptor){
	if td == nil{
		return 
	}
	glog.Infof("SCHEDULE_INFO : TaskFailed,td.Name is %v",td.Name)
	tid := td.GetTaskID()
	if td.Name != fs.tasks[tid].Name{
		return 
	}
	//delete(fs.tasks,tid)
	jid := td.Jd.GetJobID()
	delete(fs.runnableTask[jid],tid)
	td.State = TASK_FAILED
	tnode := fs.gm.TaskIDToNode(td.GetTaskID())
	fs.gm.TaskFailed(tnode)
}


func (fs *FlowScheduler)HandleTaskPlacement(tid TaskID,rid ResID){
	_,ok := fs.tasks[tid];
	if !ok {
		//vlog.Dlog(fmt.Sprintf("Error-fs-HandleTaskPlacement,the task(taskID is %d)is not exists",tid))
		return
	}
	rd,ok := fs.ress[rid];
	if !ok {
		//vlog.Dlog(fmt.Sprintf("Error-fs-HandleTaskPlacement,the res(resID is %d) is not exists",rid))
		return
	}
	
	rd.AddCurrentRunningTask(tid)
	fs.gm.TaskScheduled(fs.gm.TaskIDToNode(tid),fs.gm.ResIDToNode(rid))
	//fs.taskBindings[tid]=rid
	//if _,ok := fs.resBindings[rid];!ok{
	//	fs.resBindings[rid] = make(map[TaskID]bool)
	//}
	//fs.resBindings[rid][tid]=true
}

func (fs *FlowScheduler)HandleTaskKilled(tid TaskID) {
	
}


func (fs *FlowScheduler)ReplaceKilledTask(){
	for _ ,tnode := range fs.gm.taskNodes {
		td := tnode.td
		if td.State == TASK_FAILED {
			td.State = TASK_UNSCHEDULED
		}
	}
}


func (fs *FlowScheduler)RunScheduleIteration()(map[TaskID]ResID,map[TaskID]ResID){
	curTime := uint64(time.Now().Unix())
	fs.lastUpdateTime = uint64(curTime)
	/*alljobs := make([]*JobDescriptor,0,0)
	for _,jd := range fs.jobs {
		alljobs = append(alljobs,jd)
	}
	

	fs.gm.AddOrUpdateJobNodes(alljobs)
	*/
	// run slover to get the task to res mapping
	output,taskMappings,taskUnschedule := fs.dp.RunSolver()
	fmt.Sprintf(output)
	SortedTasks := fs.Transform(taskMappings)

	taskKilled := taskUnschedule

	/*for _,v := range taskUnschedule {
		_,ok := fs.gm.taskRuningArcs[v]
		if ok {
			taskKilled = append(taskKilled,v)
			fs.HaaaandleTaskKilled(v)
		}
	}*/
	

	taskToRes := make(map[TaskID]ResID)
	for _,sr := range SortedTasks{
		tnode := fs.gm.gcm.GetNode(sr.tnid)
		rnode := fs.gm.gcm.GetNode(sr.rnid)
		if !fs.HaveEnoughResource(tnode.td,rnode.rd){
			continue
		}
		tid := tnode.td.GetTaskID()
		rid := rnode.rd.GetResID()
		fs.HandleTaskPlacement(tid,rid)
		taskToRes[tid] = rid
	}
	taskMappings = make(map[NodeID]NodeID)

//	for tid,_ := range taskKilled {
//		fs.HandleTaskCompeletedByID(tid)
//	}


	return taskToRes,taskKilled 
}


func (fs *FlowScheduler)HaveEnoughResource(td *TaskDescriptor,rd *ResDescriptor)bool{
	pccm:= fs.cm.(*PengChengCostModel)
	pccm.SetPriority(td.Priority)
	pccm.UpdateResNode(rd)
	rr := td.ResRequest
	ar := rd.ResAvailable
	pr := rd.ResPreempt
	return  rr.Gpu<=ar.Gpu + pr.Gpu  && rr.Cpu<=ar.Cpu + pr.Cpu && rr.Memory <= ar.Memory + pr.Memory
}


func (fs *FlowScheduler)CheckJobCompleted(jd *JobDescriptor){
	if tasks,ok := fs.runnableTask[jd.GetJobID()];ok{
		cnt := len(tasks)
		for tid,_ := range tasks {
			td := fs.TaskIDToTaskDes(tid)
			if td == nil || td.State == TASK_COMPLETED{
				cnt --
				tnode := fs.gm.TaskIDToNode(tid)
				fs.gm.TaskCompleted(tnode)
			}
		}
		if cnt > 0{
			return 
		} else {
			jd.Tasks = make([]*TaskDescriptor,0,0)
		}
	}
	fs.HandleJobCompleted(jd)
}

func (fs *FlowScheduler)HandleJobCompleted(jd *JobDescriptor){
	jid := jd.GetJobID()
	fs.gm.JobCompleted(jd)
	l := -1
	for k,v := range fs.toScheduleJobs{
		if v == jd {
			l = k
		}
	}
	if l!=-1 {
		fs.toScheduleJobs = append(fs.toScheduleJobs[:l],fs.toScheduleJobs[l+1:]...)
	}
	delete(fs.jobs,jid)
	delete(fs.jobuids,jd.Uuid)
	delete(fs.runnableTask,jid)
}

func (fs *FlowScheduler)AddJob(jd *JobDescriptor){
	//taskmutex.Lock()
	jobmutex.Lock()
	//defer taskmutex.Unlock()
	defer jobmutex.Unlock()
	jid := jd.GetJobID()
	fs.jobs[jid] = jd
	fs.jobuids[jd.Uuid] = jid
	if _,ok := fs.runnableTask[jid];!ok{
		mp := make(map[TaskID]bool)
		fs.runnableTask[jid] =  mp
	}
	for _,td := range jd.Tasks{
		if td.ResRequest.Gpu != 1 {
            tmp := float64(td.ResRequest.Gpu) * 1.51
            rp := int64(tmp+0.5)
            td.Priority =uint64( int(rp) + len(jd.Tasks))
        } else {
            td.Priority = uint64(len(jd.Tasks)) + 2
        }
		fs.tasks[td.GetTaskID()] = td
		fs.runnableTask[jid][td.GetTaskID()] = true
	}
	fs.toScheduleJobs = append(fs.toScheduleJobs,jd)
}

func (fs *FlowScheduler)AddJobs(jds []*JobDescriptor){
	glog.Infof("SCHEDULE INFO: ADD JOBS :%v",jds)	
	for _,jd := range jds{
		fs.AddJob(jd)
	}
	fs.gm.AddOrUpdateJobNodes(fs.toScheduleJobs)
}
func (fs *FlowScheduler)AddResource(rd *ResDescriptor){
	rid := rd.GetResID()
	fs.ress[rid] = rd
	rds := make([]*ResDescriptor,0,0)
	rds = append(rds,rd)
	fs.gm.AddOrUpdateAllResNodes(rds)
}

func (fs *FlowScheduler)ExportGraph()string{
	return fs.gm.ExportGraph()
}



//---------------------

type SR struct{
	tnid uint64
	rnid uint64
	Priority uint64
	Gpunum uint64
}

type  ScheduleReslut []SR

func (sr ScheduleReslut)Len()int{return len(sr)}
func (sr ScheduleReslut)Swap(i,j int){sr[i],sr[j] = sr[j],sr[i]}
func (sr ScheduleReslut)Less(i,j int)bool{
	if sr[i].Priority != sr[j].Priority {
		return sr[i].Priority > sr[j].Priority
	}
	return sr[i].Gpunum > sr[j].Gpunum
}


func (fs *FlowScheduler)Transform(taskMappings map[TaskID]ResID)[]SR{
	ans := ScheduleReslut{}
	for tnid,rnid := range taskMappings {
		tnode := fs.gm.GetNode(tnid)
		td := tnode.td
		p := td.Priority
		ans = append(ans,SR{tnid,rnid,p,td.ResRequest.Gpu}) 
	}
	sort.Sort(ans)
	return ans
}
