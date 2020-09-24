package flowscheduler

const (
	RES_PU = 1
	RES_CORE = 2
	RES_OTHER = 3
)

type ResDescriptor struct{
	Rid ResID
	Uuid string
	Name string
	CurrentRunningTasks []TaskID
	Rtype ResType
	//labels []Label
	ResAvailable ResVector
	ResReserved  ResVector
	ResPreempt   ResVector 
	ResTotal     ResVector
}

func (rd *ResDescriptor)GetResID()ResID{
	return rd.Rid
}

func (rd *ResDescriptor)AddCurrentRunningTask(tid TaskID){
	rd.CurrentRunningTasks = append(rd.CurrentRunningTasks,tid)
}
