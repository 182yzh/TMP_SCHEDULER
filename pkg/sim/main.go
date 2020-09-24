package main

import (
	"flag"
	"math"
	//"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"fmt"
	"os"
	"os/exec"
	"encoding/json"
	//"sort"
	"github.com/golang/glog"
	//"k8s.io/apimachinery/pkg/util/wait"
	//"sync"
	"time"
)
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

func (jif *JobInfo)GetGpuNeedNum()int64{
	var ans int64 = 0
	//fmt.Println("jif.Attempts,%v",jif.Attempts)
	for _,det := range jif.Attempts[len(jif.Attempts)-1].Detail{
		//glog.Infof("%v",det.Gpus)
		ans+=int64(len(det.Gpus))
	}
	return ans
}

func (jif *JobInfo)GetTaskNum()int64{
	return int64(len(jif.Attempts[len(jif.Attempts)-1].Detail))
}

func (jif *JobInfo)ExcuteTime()time.Duration{
	att2 := jif.Attempts[len(jif.Attempts)-1]
	att1 := att2
	startTime,err := time.ParseInLocation("2006-01-02 15:04:05",att1.Start_time,time.Local)
	if err != nil {
		panic(err)
	}

	endTime,err := time.ParseInLocation("2006-01-02 15:04:05",att2.End_time,time.Local)
	if err != nil{
		panic(err)
	}
	ans := endTime.Sub(startTime)
	return ans
}

type Monitor struct{
	queue []*JobInfo
	ClientSet kubernetes.Interface
	total int64
	inputfile *os.File
	TimeChangeRate int64
}

func (m *Monitor)Init(filepath string){
//	fcfs.queue = make([]*v1.Pod,0,0)
	m.NewK8sClient("/home/work/.kube/config")
	m.total = 24
	m.TimeChangeRate = 1000000000
	ifile,err := os.Open(filepath)
	m.queue = make([]*JobInfo,0,0)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	m.inputfile = ifile
	dec := json.NewDecoder(m.inputfile)
	_,err = dec.Token()
	if err != nil {
		fmt.Println(err)
	}
	times := make([]int64,0,0)
	for ;dec.More();{
		jinfo := new(JobInfo)
		if err = dec.Decode(jinfo); err != nil {
			panic(err)
		}
		if len(jinfo.Attempts) == 0{
			continue
		}
		att2 := jinfo.Attempts[len(jinfo.Attempts)-1]
		att1 := att2
		startTime := att1.Start_time
		endTime  := att2.End_time
		submitTime := jinfo.Submitted_time
		//if jinfo.Status == "Failed"{
		//	continue
		//}
		if startTime == "None" || endTime == "None" || submitTime == "None" {
			continue
		}
		gpunum := jinfo.GetGpuNeedNum()
		if gpunum == 0 {
			continue
		}
		times = append(times,int64(jinfo.ExcuteTime())/m.TimeChangeRate)
		m.queue = append(m.queue,jinfo)
	}
	//sum := 0
	fmt.Println(len(m.queue))
	var mx,sum,mn int64 = 0,0,111111111
	for _,v := range times{
		sum += v
		if mx < v {mx=v}
		if mn > v {mn=v}
	}
	fmt.Println("Info: average Time,max Time,min Time",sum/int64(len(times)),mx,mn)
}


func (m *Monitor)NewK8sClient(kubeconfig string){
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	m.ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(1 * time.Second)
}

func (m *Monitor) BeginMonitor(){
	order := 1
	for {
		time.Sleep(time.Second*1)
		pods, err := m.ClientSet.CoreV1().Pods("default").List(meta_v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		gpu := int64(0)
		for _,pod := range pods.Items {
			if pod.Status.Phase == "Terminating" {
				continue
			}
			for _,container := range pod.Spec.Containers {
				rq := container.Resources.Requests["nvidia.com/gpu"]
				tmp,_ := (&rq).AsInt64()
				gpu+=tmp
			}
		}
		//	fmt.Println(gpu)
		//glog.Infof("gpu is ")
		for ;gpu < m.total + 6; {
			if len(m.queue) < 0 {
				return
			}
			jinfo := m.queue[0]
			task := jinfo.GetTaskNum()
			gpunum := jinfo.GetGpuNeedNum()/task
			name := fmt.Sprintf("app%d",order)
			if task * gpunum > 24 || gpunum > 8 {
				task = 4
				gpunum = 4
			}
			etime := int64(jinfo.ExcuteTime())/m.TimeChangeRate
			etime = int64(math.Sqrt(float64(etime))/2);
			if etime < 10 {
				etime = 10
			} else if etime > 600 {
				etime = 600
			}

			order ++
			glog.Infof("%s %d %d %d\n",name,gpunum,task,etime)
			cmd := exec.Command("../tmp-test/create_file.sh",name,fmt.Sprintf("%d",gpunum),fmt.Sprintf("%d",task),fmt.Sprintf("%d",etime))
			err := cmd.Run()
			if err != nil{
				glog.Fatalf(err.Error())
				panic(err)
			}
			cmd = exec.Command("kubectl","create","-f",fmt.Sprintf("../tmp-test/"+name+".yaml"))
			err = cmd.Run()
			if err != nil{
				glog.Fatalf(err.Error())
			}
			gpu += gpunum*task
		//	fmt.Println(gpu)
			m.queue = m.queue[1:]
		}
	}
}



func main(){
	fmt.Println("this is  a test ")
	flag.Parse()
	m := new(Monitor)
	m.Init("../tmp-test/jobinfo.txt")
	m.BeginMonitor()
}
