package flowscheduler

//import "fmt"


const (
	JOB_CREATED = 1 // is equal to Job_unscheduled
	JOB_RUNNING = 2
	JOB_COMPLETED = 3
	JOB_FAILED = 4
)


type JobDescriptor struct{
	Name string
	Uuid string
	Jid JobID
	Tasks []*TaskDescriptor
	GangScheduleNum int
	TaskUnscheduledNum int
	SubmitTime uint64
	State JobState
}

func (jd *JobDescriptor)GetJobID()JobID{
	return jd.Jid
}
