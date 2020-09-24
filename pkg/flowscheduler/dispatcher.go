package flowscheduler

import "fmt"
import "bufio"
import "strings"
import "os/exec"
import "io/ioutil"
import "container/list"
import "github.com/kubernetes-sigs/poseidon/pkg/config"
import "sort"
//import "vlog"
import "time"
import "github.com/golang/glog"


type Dispatcher struct{
	gm *GraphManager
	runOnce bool
} 


func NewDispatcher(gm *GraphManager)*Dispatcher{
    //vlog.Vlog("New Dispatcher")
    return &Dispatcher{
        runOnce:false,
        gm:gm,
    }
}



// 返回solver的输出
func (dp *Dispatcher)RunSolver() (string,map[NodeID]NodeID,map[TaskID]ResID) {
    //vlog.Vlog("dispatcher - Run Solver")
    switch dp.gm.costmodel.(type){
    case *GPUCostModel:
        //go dp.RunPinPackingSlover(ch)
        _,taskMappings := dp.RunGpuSlover()
        return "==+++==",taskMappings,nil
        //return dp.RunGpuSlover()
    case *PengChengCostModel:
        //go dp.RunPinPackingSlover(ch)
        _,taskMappings,taskUnschedule := dp.RunPCSlover()
        
        return "==+++==",taskMappings,taskUnschedule
    default:
        //vlog.Vlog("Unknow cost models")
    }
    dp.gm.sinkNode.excess = 0
    dp.gm.UpdateFlowGraph()
    input := dp.gm.ExportGraph()
    //vlog.Temlog(input)
    cmd := exec.Command("./flowscheduler/slover/cs2.exe")
    defer cmd.Wait()
    stdout, err := cmd.StdoutPipe()
    defer stdout.Close()
	if err != nil{
        return "",nil,nil
	}
    
    stdin, err := cmd.StdinPipe()
    //defer stdin.Close()
	if err != nil{
        return "",nil,nil
    }
    cmd.Start()
    stdin.Write([]byte(input))
    stdin.Close()
    out_bytes, _ := ioutil.ReadAll(stdout)
    
    //out_bytes,_ := cmd.Output()
    output := string(out_bytes)
    //vlog.Temlog(output)
    flow := dp.ParseOutput(output)
    taskMappings,taskUnschedule := dp.GetMappings(flow) 
    dp.runOnce = true
    return output,taskMappings,taskUnschedule
}

func (dp *Dispatcher)GetMappings(flow map[NodeID](map[NodeID]uint64))(map[NodeID]NodeID,map[NodeID]NodeID){
    taskMappings := make(map[NodeID]NodeID)
    leaves := dp.gm.leafNodes
    for puid,_ := range leaves {
        if _,ok := flow[puid];!ok{
            continue
        }
        toVisit := list.New()
        toVisit.PushBack(puid)
        for ;toVisit.Len() >0;{
            first := toVisit.Front()
            dst := first.Value.(NodeID)
            toVisit.Remove(first)
            node := dp.gm.gcm.GetNode(dst)
            if node.IsTaskNode() {
                taskMappings[dst]=puid
            } else {
                for src,_ := range flow[dst]{
                    toVisit.PushBack(src)
                }
            }
        }
    }
    taskKilled := make(map[TaskID]ResID)
    for dst,srcmap := range flow{
        for src,_ := range srcmap{
            dn := dp.gm.GetNode(dst)
            sn := dp.gm.GetNode(src)
            if sn.IsTaskNode() && dn.ntype == NODE_UNSCEDULED{
                taskKilled[sn.nid] = dn.nid
            } 
        }
    } 
    return taskMappings,taskKilled
}


func (dp *Dispatcher)ParseOutput(output string )map[NodeID](map[NodeID]uint64){
    //vlog.Vlog("dsipatcher -parse solver output")
    ioreader := strings.NewReader(output)
    reader := bufio.NewReader(ioreader)
    flow := make(map[NodeID](map[NodeID]uint64))
    var str string
    var src,dst,f uint64
    for ;; {
        data,_,err := reader.ReadLine();
        if err == nil {
            str = string(data)
            if str[:1] == "f"{
                fmt.Sscanf(str,"f %d %d %d",&src,&dst,&f)
                if f > 0{
                    //fmt.Sprintf(str)
                    if _,ok:=flow[NodeID(dst)];!ok{
                        flow[NodeID(dst)] = make(map[NodeID]uint64)
                    }
                    flow[NodeID(dst)][NodeID(src)] = f
                }
            }
        } else {
            break
        }
    }
    return flow
}



type  GpuNums []uint64
func (nums GpuNums)Len()int{return len(nums)}
func (nums GpuNums)Swap(i,j int){nums[i],nums[j] = nums[j],nums[i]}
func (nums GpuNums)Less(i,j int)bool{return nums[i]>nums[j]}

func (dp *Dispatcher)RunGpuSlover()(string,map[NodeID]NodeID){
    //vlog.Vlog("Run GPUCostModel Slover")
    nums := GpuNums{}
    for tnode,_ := range dp.gm.queue{
        td := tnode.td    
        if td.State == TASK_RUNNING {
            continue
        }
        nums = append(nums,td.ResRequest.Gpu)
    }
    sort.Sort(nums)
  
    for _,rnode := range dp.gm.resNodes{
        rnode.rd.ResReserved.Gpu = 0
    }
    var cur uint64 = 1<<63
    taskMappings := make(map[NodeID]NodeID)
    //vlog.Temlog("task to scedule:" + fmt.Sprintln(nums))
    for _,num := range nums{
        if cur == num{
            continue
        }
        //vlog.Temlog(fmt.Sprintln("process: ",num))
        cur = num
        gcs :=dp.gm.costmodel.(*GPUCostModel) 
        gcs.SetGpuNum(cur)
        dp.gm.sinkNode.excess = 0
        
        dp.gm.UpdateFlowGraph()
       
        _,temTaskMapping := dp.GPULimitSlover()
        
        for tid,rid :=range temTaskMapping{
            tnode := dp.gm.GetNode(tid)
            if _,ok := dp.gm.taskRuningArcs[tnode.td.GetTaskID()];ok {
                continue
            }
            taskMappings[tid]=rid
            //rnode := dp.gm.GetNode(rid)
            //rnode.rd.ResAvailable.Gpu -= cur
            //rnode.rd.ResReserved.Gpu += cur
        }
       
    }
    /*for _,rnode := range dp.gm.resNodes{
        rnode.rd.ResAvailable.Gpu += rnode.rd.ResReserved.Gpu 
        rnode.rd.ResReserved.Gpu = 0
    }*/
    //ch<-taskMappings
    return "",taskMappings;
}


func (dp *Dispatcher)GPULimitSlover()(string,map[NodeID]NodeID){
    //taskMappings := make(map[NodeID]NodeID)
    input := dp.gm.ExportGraph()
    //vlog.Temlog(fmt.Sprintln("__\n"+input+"__\n"))
    
    cmd := exec.Command("./flowscheduler/slover/cs2.exe")
    defer cmd.Wait()
    stdout, err := cmd.StdoutPipe()
	if err != nil{
        fmt.Println(err)
        return "",nil
    }
    defer stdout.Close()
    
    stdin, err := cmd.StdinPipe()
	if err != nil {
        fmt.Println(err)
        return "",nil
    }
    startTime := time.Now()
    cmd.Start()
    stdin.Write([]byte(input))
    stdin.Close()
    out_bytes, err := ioutil.ReadAll(stdout)
    endTime := time.Now()
    endTime.Sub(startTime)
    //vlog.Dlog(fmt.Sprintf("%d",dur.Nanoseconds()/1000000))
    if err != nil {
        fmt.Println(err)
        //vlog.Dlog(fmt.Sprintln(err))
        return "",nil
    }
    //out_bytes,_ := cmd.Output()
    output := string(out_bytes)
    temstrs := strings.Split(output,"\n")
    temoutput := ""
    if len(temstrs) <= 16 {
        fmt.Println(output)
        return "",nil
    }
    for _,v := range temstrs[16:]{
        temoutput+= v + "\n"
    }
    //vlog.Temlog(fmt.Sprintln("\npart output:\n",temoutput))
    flow := dp.ParseOutput(output)
    taskMappings,_ := dp.GetMappings(flow) 
    dp.runOnce = true
    return "",taskMappings
}


func (dp *Dispatcher)ExportGraphWithoutRunningTasks()string{
    return dp.gm.ExportGraphWithoutRunningTasks()
}












func (dp *Dispatcher)RunPCSlover()(string,map[NodeID]NodeID,map[TaskID]ResID){
	//get the information about all the task 
    toschedule := Details{}
    haveschedule := Details{}
    haveschedule = append(haveschedule,TaskDetail{0,0,0,-1,nil})
    for _,tnode := range dp.gm.taskNodes{
        td := tnode.td  
        if tnode.td.State == TASK_RUNNING {
            if td.IsGangSchedule == true{
				continue
			}
			arc,ok := dp.gm.taskRuningArcs[tnode.td.GetTaskID()]
            if ok== false{
                fmt.Println(" in sp. RuPCSlover, Error ")
            }
            rnode :=arc.dst
            rd := rnode.rd
            haveschedule = append(haveschedule,TaskDetail{td.Priority,td.ResRequest.Gpu,rd.ResAvailable.Gpu,td.StartTime,td})
        } else {
			// res avaiable 越大，其优先度越低，保证在gpu需求以及优先级相同时，未被调度的任务的优先度要低于已被调度的任务的优先
			toschedule = append(toschedule,TaskDetail{td.Priority,td.ResRequest.Gpu,1024,td.SubmitTime,td})
        }
    }
	glog.Infof("SCHEDULE INFO|DISPATCHER_INFO,toscheduletask %v\nhavesche task info:%v",toschedule,haveschedule)
    sort.Sort(toschedule)
    sort.Sort(haveschedule)
  
	
    for _,rnode := range dp.gm.resNodes{
        rnode.rd.ResReserved.Gpu = 0
    }
    cur := TaskDetail{uint64(1<<31),uint64(1<<31),0,int64(1),nil}
    taskMappings := make(map[NodeID]NodeID)
    taskUnschedule := make(map[TaskID]ResID)
	n := len(toschedule)
//	fmt.Println(toschedule)
    //try to schedule all task which have same priority and gpu need
	for i:=n-1;i>=0;i-=1{
		dl := toschedule[i]
        if dl.p == cur.p && dl.gpu == cur.gpu{
            continue
        }
        cur = dl
	//	fmt.Println("dp.run PCslover,cur is:",cur)
	//	fmt.Println("dp.RunPCslover,res info: ",dp.gm.GetResInfo())	
        pccm := dp.gm.costmodel.(*PengChengCostModel)
        pccm.SetPriority(cur.p)
        pccm.SetGpunum(cur.gpu)
        dp.gm.sinkNode.excess = 0

		cnt := 0
        for _,v := range toschedule {
            if v.p == cur.p && v.gpu == cur.gpu {
                cnt ++
            }
        }
		// try binary search to check can we schedule some task by kill some task
        l,r := 0,len(haveschedule)-1
		for;(r>0&&Compare(cur,haveschedule[r])); {
			r--
		}
        for ;l<r ;{
            m := int((l+r)/2)
            pccm.SetTinfo(haveschedule[m])
            dp.gm.UpdateFlowGraph()
            _,temTaskMapping,_ := dp.PengChengSlover()
            if len(temTaskMapping) == cnt {
                r = m
            } else {
                l = m+1
            }
        }

        pccm.SetTinfo(haveschedule[l])
        dp.gm.UpdateFlowGraph()
        _,temTaskMapping,_ := dp.PengChengSlover()
        temTaskMapping = dp.CheckScheduleResult(temTaskMapping)
        temTaskKilled := pccm.GetTaskKilled(temTaskMapping)

        for tid,rid := range temTaskKilled{
            dp.gm.HandleTaskKilled(tid)
            taskUnschedule[tid]=rid
        }

        for tnid,rnid := range temTaskMapping{
            tnode := dp.gm.GetNode(tnid)
            if _,ok := dp.gm.taskRuningArcs[tnode.td.GetTaskID()];ok {
                continue
            }
            taskMappings[tnid]=rnid
            rnode := dp.gm.GetNode(rnid)
            rnode.rd.ResAvailable.Gpu -= pccm.Gpunum
			if rnode.rd.ResAvailable.Gpu > 8 {
				glog.Infof("SCHEDULE_INFO : ERROR,Rnode.GPU>8,info: rnode.Name :%s,temtaskmappings : %v",rnode.rd.Name,temTaskMapping)
				glog.Infof("SCHEDULE_INFO : rd.Resavai: %d, rd.ResPreet:%d",rnode.rd.ResAvailable.Gpu,rnode.rd.ResPreempt.Gpu)
			}
            rnode.rd.ResReserved.Gpu += pccm.Gpunum
        }
		
    }
    //fmt.Pffffrintln("test")
    for _,rnode := range dp.gm.resNodes{
        rnode.rd.ResAvailable.Gpu += rnode.rd.ResReserved.Gpu 
        rnode.rd.ResReserved.Gpu = 0
    }
    //ch<-taskMappings
    return "",taskMappings,taskUnschedule;
}


func (dp *Dispatcher)CheckScheduleResult(temTaskMapping map[NodeID]NodeID)map[NodeID]NodeID{
    processed := make(map[NodeID]bool)
    ans := make(map[NodeID]NodeID)
	loginfo := make(map[string]NodeID)
    for tnid,rnid := range temTaskMapping{
        if processed[tnid] {
            continue
        }
        processed[tnid] = true
        tnode := dp.gm.GetNode(tnid)
        if tnode == nil  || tnode.td == nil{
            fmt.Println("Error ,in dispatcher ,tnode or td is nil")
        }
        td := tnode.td
		loginfo[td.Name] = rnid
        jd := td.Jd
        tmp := make(map[NodeID]NodeID)
        schedule := true
        for _,taskdes := range jd.Tasks {
            tid := taskdes.GetTaskID()
            tnode := dp.gm.TaskIDToNode(tid)
            if tnode == nil  || tnode.td == nil{
                fmt.Println("Error ,in dispatcher ,td or tnode tnode is nil")
            }
            if rnid,ok := temTaskMapping[tnode.nid];ok{
                tmp[tnode.nid] =temTaskMapping[tnode.nid]
				loginfo[taskdes.Name] = rnid
            }else {
                schedule = false
            }
        }
        if schedule {
            for k,v := range tmp{
                ans[k]=v
            }
        }
    }
	glog.Infof("SCHEDULE INFO:Task Schedule: %v",loginfo)
    return ans
}  

func (dp *Dispatcher)PengChengSlover()(string,map[NodeID]NodeID,map[NodeID]NodeID){
    input := dp.gm.ExportGraph()
    if config.GetFlowDebug(){
		fmt.Println(input)
	}
    cmd := exec.Command("../../pkg/flowscheduler/slover/cs2.exe")
    //cmd := exec.Command("./flowscheduler/slover/cs2.exe")
	defer cmd.Wait()
    stdout, err := cmd.StdoutPipe()
	if err != nil{
        //fmt.Println(err)
        return "",nil,nil
    }
    defer stdout.Close()
    
    stdin, err := cmd.StdinPipe()
	if err != nil {
        fmt.Println(err)
        return "",nil,nil
    }
    startTime := time.Now()
    cmd.Start()
    stdin.Write([]byte(input))
    stdin.Close()
    out_bytes, err := ioutil.ReadAll(stdout)
    endTime := time.Now()
    endTime.Sub(startTime)
    if err != nil {
        fmt.Println(err)
        return "",nil,nil
    }
    output := string(out_bytes)
	if config.GetFlowDebug(){ 
		fmt.Println("in run pc slover",output)
    }
	temstrs := strings.Split(output,"\n")
    temoutput := ""
    if len(temstrs) <= 16 {
        //fmt.Println(output)
        //vlog.Dlog("len temstrs < 16\n")
        return "",nil,nil
    }
    for _,v := range temstrs[16:]{
        temoutput+= v + "\n"
    }
    //vlog.Temlog(fmt.Sprintln("\npart output:\n",temoutput))
    flow := dp.ParseOutput(output)
    taskMappings,taskUnschedule := dp.GetMappings(flow) 
    dp.runOnce = true
    return "",taskMappings,taskUnschedule
}




/*

func (dp *Dispatcher) RunPinPackingSlover(ch chan map[NodeID]NodeID )map[NodeID]NodeID{
    woods := make([]PinPacking,0,0)
    boxs := make([]PinPacking,0,0)
    for _,tnode := range dp.gm.taskToNodeMap{
        td := tnode.taskDes
        if td.taskType == TASK_ROOT{
            continue
        }
        woods = append(woods,PinPacking{td.resourceRequest.gpu,tnode.id})
    }

    for _,rnode := range dp.gm.resToNodeMap {
        rd := rnode.resDes
        boxs = append(boxs,PinPacking{rd.availableRes.gpu,rnode.id})
    }
    taskMappings := PinPackingSlover(woods,boxs)
    ch<-taskMappings
    return taskMappings
}



func ReadFile() {
    // file, err := os.Open("./test.txt")
    file, err := os.OpenFile("./test.txt", os.O_CREATE|os.O_RDONLY, 0666)
    if err != nil {
        return
    }
    defer file.Close()    //关闭文件

    reader := bufio.NewReader(file)    //带缓冲区的读写
    for {
        str, err := reader.ReadString('\n')    // 以\n为分隔符来读取
        if err != nil {
            return
        }
    }
}
*/
