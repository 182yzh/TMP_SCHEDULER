/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8sclient

import (
	"github.com/kubernetes-sigs/poseidon/pkg/firmament"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/golang/glog"
	config2 "github.com/kubernetes-sigs/poseidon/pkg/config"
	"k8s.io/apimachinery/pkg/util/wait"
	"sync"
	"time"
)

var ClientSet kubernetes.Interface

// BindPodToNode call Kubernetes API to place a pod on a node.
func BindPodToNode() {
	for {
		bindInfo := <-BindChannel
		err := ClientSet.CoreV1().Pods(bindInfo.Namespace).Bind(&v1.Binding{
			TypeMeta: meta_v1.TypeMeta{},
			ObjectMeta: meta_v1.ObjectMeta{
				Name: bindInfo.Name,
			},
			Target: v1.ObjectReference{
				Namespace: bindInfo.Namespace,
				Name:      bindInfo.Nodename,
			}})
		if err != nil {
			glog.Errorf("Could not bind pod:%s to nodeName:%s, error: %v", bindInfo.Name, bindInfo.Nodename, err)
		} else {
			glog.Infof("Bind pod to node : %s to %s",bindInfo.Name, bindInfo.Nodename)
		}


		//pod,err := ClientSet.CoreV1().Pods(bindInfo.Namespace).Get(bindInfo.Name, meta_v1.GetOptions{})
	    //if err != nil {
    	//    glog.Fatalf("Could not get pod:%s in namespace:%s, error: %v", bindInfo.Name, bindInfo.Namespace, err)
    	//}
		//glog.Infof("SCHEDULE INFO: PODINFO After bind,%v",pod)

	}
}

// DeletePod calls Kubernetes API to delete a Pod by its namespace and name.
func DeletePod(podName string, namespace string) {
	
/*
	pod,err := ClientSer.CoreV1().Pods(namespace).Get(podName,&meta_v1.GetOptions{})
	if err != nil {
        glog.Fatalf("Could not get pod:%s in namespace:%s, error: %v", podName, namespace, err)
    }
*/
	//jobname := pod.Labels[]
	policy := meta_v1.DeletePropagationBackground
	err := ClientSet.CoreV1().Pods(namespace).Delete(podName, &meta_v1.DeleteOptions{PropagationPolicy: &policy})
	if err != nil {
		glog.Fatalf("Could not delete pod:%s in namespace:%s, error: %v", podName, namespace, err)
	}
}


func WaitForResReady(namespace,podName,nodeName string){
	pod,err := ClientSet.CoreV1().Pods(namespace).Get(podName, meta_v1.GetOptions{})
	if err != nil {
		glog.Fatalf("Could not get pod:%s in namespace:%s, error: %v", podName, namespace, err)
	}
	req := int64(0)
	for _,container := range pod.Spec.Containers {
		rq := container.Resources.Requests["nvidia.com/gpu"]
		tmp,_ := (&rq).AsInt64()
		req+=tmp
	}
//	glog.Infof("SCHEDULE INFO: %s gpu request is :%d",podName,req);
//	glog.Infof("SCHEDULE INFO: PODINFO :%v",pod)
//	node,_ := ClientSet.CoreV1().Nodes().Get(nodeName, meta_v1.GetOptions{})
//	glog.Infof("SCHEDULE INFO: NODEINFO: %v",node)

	all := int64(0)
	avi := int64(0)
	node,err := ClientSet.CoreV1().Nodes().Get(nodeName, meta_v1.GetOptions{})
	if err != nil {
	    glog.Fatalf("Could not get nodes:%s in namespace:%s, error: %v", nodeName, namespace, err)
    }
    rq := node.Status.Allocatable["nvidia.com/gpu"]
    all,_ = (&rq).AsInt64()
    glog.Infof("SCHEDULE INFO: %s have res allocated gpu :%d",nodeName,all)
	for;avi < req;{
		avi = all
		pods,err := ClientSet.CoreV1().Pods(namespace).List(meta_v1.ListOptions{})
		if err != nil {
			glog.Infof("err ,client can not get pods")
		}
		for _,pod := range pods.Items{
			bindName := pod.Spec.NodeName
			if bindName != nodeName {
				continue
			}
			used := int64(0)
			for _,container := range pod.Spec.Containers {
				rq := container.Resources.Requests["nvidia.com/gpu"]
				tmp,_ := (&rq).AsInt64()
				used+=tmp
			}
			avi -= used
		}
	//	glog.Infof("SCHEDULE INFO: %s have res available gpu :%d",nodeName,all)
	}



	return
}



// GetClientConfig returns a kubeconfig object which to be passed to a Kubernetes client on initialization.
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

// New initializes a firmament and Kubernetes client and starts watching Pod and Node.
func New(schedulerName string, kubeConfig string, kubeVersionMajor, kubeVersionMinor int, firmamentAddress string,fc firmament.FirmamentSchedulerClient){

	config, err := GetClientConfig(kubeConfig)
	if err != nil {
		glog.Fatalf("Failed to load client config: %v", err)
	}
	//

	config.QPS = config2.GetQPS()
	config.Burst = config2.GetBurst()

	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create connection: %v", err)
	}
	
	glog.Info("k8s newclient called")
	stopCh := make(chan struct{})
	go NewPodWatcher(kubeVersionMajor, kubeVersionMinor, schedulerName, ClientSet, fc).Run(stopCh, 10)
	go NewNodeWatcher(ClientSet, fc).Run(stopCh, 10)
	go NewK8sPodWatcher(kubeVersionMajor, kubeVersionMinor, schedulerName, ClientSet, fc).controller.Run(stopCh)

	// We block here.
	<-stopCh
}

func init() {

	glog.Info("k8sclient init called")
	BindChannel = make(chan BindInfo, 1000)
	PodToK8sPodLock = new(sync.Mutex)
	ProcessedPodEventsLock = new(sync.Mutex)
	PodToK8sPod = make(map[PodIdentifier]*v1.Pod)
	ProcessedPodEvents = make(map[PodIdentifier]*v1.Pod)
}

// Run starts a pod watcher.
func BindPodWorkers(stopCh <-chan struct{}, nWorkers int) {

	for i := 0; i < nWorkers; i++ {
		go wait.Until(BindPodToNode, time.Second, stopCh)
	}

	<-stopCh
	glog.Info("Stopping RunBindPods")
}
