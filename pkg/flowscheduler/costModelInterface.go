package flowscheduler

type CostModel interface {
	UnschedAggToSink(jd *JobDescriptor) *ArcDescriptor
	TaskContinuation(td *TaskDescriptor) *ArcDescriptor
	LeafRescourceToSink(rd *ResDescriptor) *ArcDescriptor
	TaskToUnscheduledAgg(td *TaskDescriptor) *ArcDescriptor
	TaskPreferdResource(td *TaskDescriptor)  []*ResDescriptor
	TaskNodeToResource(td *TaskDescriptor,rd *ResDescriptor)*ArcDescriptor
	UpdateTaskNode(td *TaskDescriptor)
	UpdateResNode(rd *ResDescriptor)
} 




