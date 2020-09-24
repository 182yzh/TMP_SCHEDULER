package flowscheduler

import (
	"fmt"
	"time"
//	"github.com/kubernetes-sigs/poseidon/pkg/config"
)


type JobInfo struct{
	p uint64
	g []uint64
	num int
}

type JIS []JobInfo 

var curid uint64
func GenerateRes(resavi map[uint64]uint64)[]*ResDescriptor{
	ans := make([]*ResDescriptor,0,0)
	var i uint64 = 0
	for k,v := range resavi{
		for i =0 ;i<v;i++{
			ans = append(ans,NewRes(k))
		}
	}
	return ans
}

func GenerateJob(resreq JobInfo)[]*JobDescriptor{
	ans := make([]*JobDescriptor,0,0)
	//var i uint64 
	for k:=0;k< resreq.num; k++ {
		ans = append(ans,NewJob(resreq.g,resreq.p))
	}
	return ans
}

func NewJob(reqs []uint64,p uint64)*JobDescriptor{
	jd := new(JobDescriptor)
	
	jd.Jid = curid 
	curid ++ 
	jd.Tasks = make([]*TaskDescriptor,0,0)
	jd.Name = fmt.Sprintf("job-%d",jd.Jid)
	jd.State = JOB_CREATED
	jd.SubmitTime = uint64(time.Now().UnixNano())

	for _,req := range reqs{
		td := new(TaskDescriptor)
		td.Tid = curid
		curid ++ 
		jd.Tasks = append(jd.Tasks,td)
		td.Name = fmt.Sprintf("job-%d task-%d",jd.Jid,td.Tid)
		td.State = TASK_UNSCHEDULED
		td.Jd = jd
		td.SubmitTime =  time.Now().UnixNano()
		td.ResRequest.Gpu = req
		//fmt.Println(fmt.Sprintf("task:%d,gpu need:%d",td.Tid,req))
		td.Priority = p
		if len(reqs) > 1 {
			td.IsGangSchedule = true
		}
	}
	
	return jd
}

func NewRes(avi uint64)*ResDescriptor{
	rd := new(ResDescriptor)
	rd.ResAvailable.Gpu = avi
	rd.Rid = curid
	rd.Name=fmt.Sprintf("Node%d",curid)
	curid ++ 
	rd.ResTotal.Gpu = avi
	rd.Rtype = RES_PU
//	fmt.Println(fmt.Sprintf("res:%d, gpu avi:%d",rd.Rid,avi))
	return rd
}




func TestGraphMain(){
	//return
	//fmt.Println(config.GetEnablePriority(),config.GetEnableTaskResNeed())


	curid = 1
	
	resavi := make(map[uint64]uint64)
	resavi[8]=3
	

	resreq := make([]JIS,0,0)
	
	tem := JIS{}
	tem = append(tem,JobInfo{1,[]uint64{1},5})
	tem = append(tem,JobInfo{2,[]uint64{8},1})
	resreq = append(resreq,tem)

	tem = JIS{}
	tem = append(tem,JobInfo{6,[]uint64{2,2,2,2},1})
	tem = append(tem,JobInfo{4,[]uint64{4,4},1})
	tem = append(tem,JobInfo{1,[]uint64{1},4})
	resreq = append(resreq,tem)
	

//	tem = JIS{}
	//tem = append(tem,JobInfo{4,[]uint64{4,4},1})
//	resreq = append(resreq,tem)

	
	
	ress := GenerateRes(resavi)
	//jobs := GenerateJob(resreq)
///	for _,v := range ress{
//		fmt.Println(v)
//	}

	fs := NewFlowScheduler()
	for _,v := range ress{
		fs.AddResource(v)
	}
	
	
	for _,jis := range resreq{
		for _,ji := range jis{
			jobs := GenerateJob(ji)

			for _,jd := range jobs{
				fmt.Sprintln(jd)
			}
			fs.AddJobs(jobs)
		}
			//fs.gm.UpdateFlowGraph()
			//fmt.Println(fs.ExportGraph())
			shc,tkk := fs.ScheduleAllJobs()
			fmt.Println("task killed :",tkk)
			fmt.Println("task scheduled",shc)
	}
	//fs.AddJobs(jobs)
	
	//fmt.Println("num: ",len(shc))
}
