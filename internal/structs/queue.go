package structs

type QueueStatusRequest struct {
	BusyWorkers    int
	FailureTasks   int
	SubmittedTasks int
	SuccessTasks   int
}
