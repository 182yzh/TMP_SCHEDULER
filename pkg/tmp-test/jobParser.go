package main

import (
	"fmt"
	"encoding/json"
	"os"
	"time"
	"sort"
	"strconv"
)

type JobParser struct {
	Jobs []*JobInfo
	inputfile *os.File
	outputfile *os.File
	//num int
}

func NewJobParser(inputfilename,outputfilename string)*JobParser{
	jp := new(JobParser)
	jp.Jobs = make([]*JobInfo,0,0)
	inputfile,err := os.OpenFile(inputfilename,os.O_RDWR|os.O_CREATE,0)
	//inputfile := os.Open(inputfilename)
	if err != nil {
		panic(err.Error())
	}
	outputfile,err := os.OpenFile(outputfilename,os.O_RDWR|os.O_CREATE,0)
	if err != nil {
		panic(err.Error())
	}
	jp.outputfile = outputfile
	//jp.num = num
	jp.inputfile = inputfile
	return jp
}

type DetailInfo struct {
    Ip string `json:"ip"`
    Gpus []string  `json:"gpus"`
}

type AttemptInfo struct{
    Start_time string `json:"start_time"`
    End_time string `json:"end_time"`
    Detail []DetailInfo `json:"detail"`
}

type JobInfo struct{
    Status string `json:"status"`
    Vc string `json:"vc"`
    Jobid string `json:"jobid"`
    Attempts []AttemptInfo   `json:"attempts"`
    Submitted_time string `json:"submitted_time"`
	User string `json:"user"`
	Test string `json:"server"`
}
type SortJobs []*JobInfo

func(s SortJobs)Len()int{return len(s)}
func(s SortJobs)Swap(i,j int){s[i],s[j] = s[j],s[i]}
func(s SortJobs)Less(i,j int)bool{
	return cmp(s[i].Submitted_time,s[j].Submitted_time) < 0
}


func checkDetail(att AttemptInfo,tem map[string]bool)bool{
	for _,v := range att.Detail{
		ip := v.Ip
		if ip == ""{
			return true
		}
		tem[ip] = true
		num,_ := strconv.Atoi(ip[1:])
		if num >= 424{
			fmt.Println(ip)
			return false
		}
	}
	return true
}

func (jif *JobInfo)ExcuteTime()time.Duration{
	
	att2 := jif.Attempts[len(jif.Attempts)-1]
	att1 := att2
	startTime,err := time.ParseInLocation("2006-01-02 15:04:05",att1.Start_time,time.Local)
	if err != nil {
		fmt.Println(err)
	}

	endTime,err := time.ParseInLocation("2006-01-02 15:04:05",att2.End_time,time.Local)
	if err != nil{
		fmt.Println(err)
	}
	ans := endTime.Sub(startTime)
	
	return ans
}

func (jif *JobInfo)GetGpuNeedNum()uint64{
	var ans uint64 = 0
	for _,det := range jif.Attempts[len(jif.Attempts)-1].Detail{
		ans+=uint64(len(det.Gpus))
	}
	return ans
}

func cmp(isTime,jsTime string)int64{
	itime,err := time.ParseInLocation("2006-01-02 15:04:05",isTime,time.Local)
	if err != nil {
		panic(err.Error())
	}
	jtime,err := time.ParseInLocation("2006-01-02 15:04:05",jsTime,time.Local)
	if err != nil{
		panic(err.Error())
	}

	ans := itime.Sub(jtime)
	//vlog.Dlog(fmt.Sprintln(isTime,jsTime,ans,int64(ans.Seconds())))
	return int64(ans.Seconds()) 
}



func (jp *JobParser)NextJobInfo(){
	dec := json.NewDecoder(jp.inputfile)
	t,err := dec.Token()
//	cnt := 0
	if err != nil {
		fmt.Println(err)
	}
	fmt.Sprintln(t)
	for ;dec.More();{
//		cnt++
		jinfo := new(JobInfo)
		if err = dec.Decode(jinfo); err != nil {
			fmt.Println(err)
			return
		}
		if len(jinfo.Attempts) == 0 {
			continue
		}
		att2 := jinfo.Attempts[len(jinfo.Attempts)-1]
		startTime := att2.Start_time
		endTime  := att2.End_time
		submitTime := jinfo.Submitted_time
		if startTime == "None" || endTime == "None" || submitTime == "None" {
			continue
		}
		if cmp(submitTime,"2017-10-10 00:00:00") < 0 || cmp(endTime ,"2017-12-01 00:00:00")>0{
			continue
		}
		gpunum := jinfo.GetGpuNeedNum()
		if gpunum == 0 { 
			continue
		}
		jp.Jobs = append(jp.Jobs,jinfo)
	}
	t,err = dec.Token()
	if err != nil {
		fmt.Println(err)
	}
	sjobs := SortJobs{}
	sjobs = append(sjobs,jp.Jobs...)
	fmt.Println(len(sjobs))
	sort.Sort(sjobs)
	jobs := make([]*JobInfo,0,0)
	for i,job:=range sjobs{
		if i%40 == 0 {
			jobs = append(jobs,job)
		}
	}
//	fmt.Println(len(jobs))
	_,err = jp.outputfile.Write([]byte("["));
   // fmt.Println(err)
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
	n := len(jobs)
	for i:=0;i<n;i++{
		jsons, errs := json.Marshal(jobs[i])
	    if errs != nil {
	        panic(err.Error())
	    }
		_,err = jp.outputfile.Write(jsons);
		// fmt.Println(err)
		if err != nil {
			fmt.Println(err)
			panic(err.Error())
		}
		s := ",\n"
		if i == n-1 {
			s="]\n"
		}
        _,err = jp.outputfile.Write([]byte(s))
        if err != nil {
            fmt.Println(err)
            panic(err.Error())
        }
	}

/*	jsons, errs := json.Marshal(jobs)
	if errs != nil {
		panic(err.Error())
	}
	_,err = jp.outputfile.Write(jsons)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
*/
	return
}

func (jp *JobParser)test(){
	dec := json.NewDecoder(jp.inputfile)
    t,err := dec.Token()
    cnt := 0
    if err != nil {
        fmt.Println(err)
    }
    fmt.Sprintln(t)
	jinfo := &JobInfo{}
    for ;dec.More();{
		cnt++
		//jinfo := &JobInfo{}
		dec.Decode(jinfo)
		//if cnt >1000000{
		//	break
		//}
	}
	fmt.Println(jinfo)
	fmt.Println(cnt)

}

func main() {
	jp := NewJobParser("cluster_job_log","jobinfo.txt")
	fmt.Println("rest")
	jp.NextJobInfo()
//	jp.test()
	return
}

