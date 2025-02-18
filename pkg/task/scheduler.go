package task


type Scheduler interface {
	Start() 
	Stop() 

	AddExecutor(nodeId int64) 
	RemoveExecutor(nodeId int64) 

	Add(task Task) error 
	Dispatch(node int64) 
	RemoveByNode(node int64) 
}