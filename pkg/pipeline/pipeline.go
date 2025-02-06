package pipeline

import "time"

type Pipeline interface {
	Start() error
	Close()

	Add(node ...Node)
}

type pipeline struct {
	nodes []*nodeCtx

	inputChannel    chan Msg
	nodeTtInterval  time.Duration
	enableTtChecker bool
}

func (p *pipeline) Add(nodes ...Node) {
	for _, node := range nodes {
		p.addNode(node)
	}
}

func (p *pipeline) addNode(node Node) {
	nodeCtx := NewNodeCtx(node)

	if p.enableTtChecker {
		//
		// nodeCtx.Checker = NewChecker(name, manager)
	}

	if len(p.nodes) != 0 {
		p.nodes[len(p.nodes)-1].Next = nodeCtx
	} else {
		p.inputChannel = nodeCtx.InputChannel
	}

	p.nodes = append(p.nodes, nodeCtx)
}

func (p *pipeline) Start() error {
	return nil
}

func (p *pipeline) Close() {
	for _, node := range p.nodes {
		if node.Checker != nil {
			node.Checker.Close()
		}
	}
}

func (p *pipeline) process() {
	if len(p.nodes) == 0 {
		return
	}

	curNode := p.nodes[0]
	for curNode != nil {
		if len(curNode.InputChannel) == 0 {
			return
		}

		input := <-curNode.InputChannel
		output := curNode.node.Operate(input)
		if curNode.Checker != nil {
			curNode.Checker.Check()
		}

		if curNode.Next != nil && output != nil {
			curNode.Next.InputChannel <- output
		}

		curNode = curNode.Next
	}
}
