package main

type CommandChain struct {
	Cmd        string
	Next, Prev *CommandChain
}

func NewCommandChain(cmd string) *CommandChain {
	return &CommandChain{Cmd: cmd}
}

func (c *CommandChain) SetNext(next *CommandChain) {
	c.Next = next
}

func (c *CommandChain) SetPrev(prev *CommandChain) {
	c.Prev = prev
}

func (c *CommandChain) Last() *CommandChain {
	cur := c
	for cur.Next != nil {
		cur = cur.Next
	}
	return cur
}

func (c *CommandChain) Add(chain *CommandChain) {
	lastCmdChain := c.Last()
	lastCmdChain.SetNext(chain)
	chain.SetPrev(lastCmdChain)
}
