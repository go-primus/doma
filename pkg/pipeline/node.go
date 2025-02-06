package pipeline

type Node interface {
	Name() string
	Operate(in Msg) Msg
}

type Checker interface {
	Check()
	Close()
}

type nodeCtx struct {
	node         Node
	InputChannel chan Msg

	Next    *nodeCtx
	Checker Checker
}

func NewNodeCtx(node Node) *nodeCtx {
	return &nodeCtx{
		node:         node,
		InputChannel: make(chan Msg),
	}
}

type BaseNode struct {
	name           string
	maxQueueLength int32
}

func (node *BaseNode) Name() string {
	return node.name
}

func (node *BaseNode) MaxQueueLength() int32 {
	return node.maxQueueLength
}

func NewBaseNode(name string) *BaseNode {
	return &BaseNode{
		name: name,
	}
}
