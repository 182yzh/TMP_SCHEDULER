package main


import (
	"flag"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"fmt"
	"sort"
	"github.com/golang/glog"
	//"k8s.io/apimachinery/pkg/util/wait"
	//"sync"
	"time"
)

type NodeInfo struct {
	gpu int64
	name string
}

type  NodesInfo []NodeInfo

func (ni NodesInfo)Len()int{return len(ni)}
func (ni NodesInfo)Swap(i,j int){ni[i],ni[j]=ni[j],ni[i]}
func (ni NodesInfo)Less(i,j int)bool{return ni[i].gpu < ni[j].gpu}


type FcfsScheduler struct{
	queue []*v1.Pod
	ClientSet kubernetes.Interface
}

func (fcfs *FcfsScheduler)Init(){
	fcfs.queue = make([]*v1.Pod,0,0)
	fcfs.NewK8sClient("/home/work/.kube/config")
}

func (fcfs *FcfsScheduler)NewK8sClient(kubeconfig string){
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	fcfs.ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	
	pods, err := fcfs.ClientSet.CoreV1().Pods("default").List(meta_v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	//fmt.Println(pods.Items)
	time.Sleep(1 * time.Second)
//	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

func (fcfs *FcfsScheduler) BeginSchedule(){
	for {
		time.Sleep(time.Second*1)
		pods, err := fcfs.ClientSet.CoreV1().Pods("default").List(meta_v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		cnt := 0
		pendingPods := make([]v1.Pod,0,0)
		for _,pod := range pods.Items {
			if pod.Status.Phase == "Pending" {
				cnt ++
				pendingPods = append(pendingPods,pod)
			}
		}
		
	//	fmt.Println(pendingPods)
		if cnt == 0 {
			continue
		}
		
		gpu:= int64(0)
		pod := pendingPods[0]
		fmt.Println(pod.Name)
		for _,container := range pod.Spec.Containers {
			rq := container.Resources.Requests["nvidia.com/gpu"]
			tmp,_ := (&rq).AsInt64()
			gpu+=tmp
		}
		nodes := fcfs.ResState()
		ninfo := NodesInfo{}
		for k,v := range nodes{
			ninfo = append(ninfo,NodeInfo{v,k})
		}
		sort.Sort(ninfo)
		fmt.Println(nodes,gpu,cnt)
		//假设所有任务都是同一个namespace
		result := make(map[string]string)
		for _,nif := range ninfo{
			k := nif.name
			v := nif.gpu
			for;cnt>0 && v>=gpu;{
				result[pod.Name]=k
				cnt--
				v--
			}
		}
		if cnt != 0{
			continue
		}
		for podname,nodename := range result{
			err := fcfs.ClientSet.CoreV1().Pods("default").Bind(&v1.Binding{
			TypeMeta: meta_v1.TypeMeta{},
			ObjectMeta: meta_v1.ObjectMeta{
				Name: podname,
			},
			Target: v1.ObjectReference{
				Namespace: "default",
				Name:      nodename,
			}})
			if err != nil {
				glog.Errorf("Could not bind pod:%s to nodeName:%s, error: %v", podname, nodename, err)
			} else {
				glog.Infof("Bind pod to node : %s to %s",podname, nodename)
			}
		}
		time.Sleep(time.Second*1)
	}
}



func (fcfs *FcfsScheduler)ResState()map[string]int64{
	nodes, err := fcfs.ClientSet.CoreV1().Nodes().List(meta_v1.ListOptions{})
    if err != nil {
        panic(err)
    }
	ans := make(map[string]int64)
	for _,node := range nodes.Items {
		rto := node.Status.Allocatable["nvidia.com/gpu"]
		ans[node.Name],_ = (&rto).AsInt64()
	}

	pods,err := fcfs.ClientSet.CoreV1().Pods("default").List(meta_v1.ListOptions{})
	if err != nil {
			panic(err)
	}
	for _,pod := range pods.Items{
		if pod.Spec.NodeName == ""{
			continue
		}
		req := int64(0)
	    for _,container := range pod.Spec.Containers {
		    rq := container.Resources.Requests["nvidia.com/gpu"]
			tmp,_:= (&rq).AsInt64()
			req+=tmp
		}
		//fmt.Println(pod.Name,pod.Spec.NodeName,req)
		ans[pod.Spec.NodeName] -= req
	}
	return ans
}

func main(){
	fmt.Println("this i s  a test ")
	flag.Parse()
	fcfs := new(FcfsScheduler)
	fcfs.Init()
	fcfs.BeginSchedule()
}
