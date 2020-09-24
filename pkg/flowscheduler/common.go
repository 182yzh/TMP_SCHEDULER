package flowscheduler

import "github.com/kubernetes-sigs/poseidon/pkg/config"

type NodeID = uint64
type NodeType = uint8
type TaskID = uint64
type TaskType = uint8
type JobID = uint64
type ResID = uint64
type JobState = uint8
type ResType = uint8
type TaskState = uint8
type ArcType = uint8

type ResVector struct {
	Cpu uint64
	Memory uint64
	Gpu uint64
}

type Label struct{
	key string
	value string
}

type ArcDescriptor struct {
	capLower uint64
	capUpper uint64
	cost     uint64
}

func Test()string {
	return "test \n"
}


type TaskDetail struct{
    p uint64
	gpu uint64
	resavi uint64
	Time int64
	td *TaskDescriptor
}

type  Details []TaskDetail
func(tdl Details)Len()int{return len(tdl)}
func (tdl Details)Swap(i,j int){tdl[i],tdl[j] = tdl[j],tdl[i]}
func (tdl Details)Less(i,j int)bool{
    return Compare(tdl[i],tdl[j])
}

func Compare(i,j TaskDetail)bool{
    if config.GetEnablePriority() &&  i.p != j.p {
        return i.p < j.p
    }
    if config.GetEnableTaskResNeed()  &&  i.gpu != j.gpu{
        return i.gpu < j.gpu
	}
	if config.GetEnableResAvai() && i.resavi != j.resavi {
		return i.resavi > j.resavi
	}
	if config.GetEnableRunTime()  && i.Time != j.Time {
		return i.Time >= j.Time
	}
	return true
}

