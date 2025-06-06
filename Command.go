package simple_worker_pool

type ICommand interface {
	Execute()
}

func CreateCommand(f func()) ICommand {
	return &command{f: f}
}

type command struct {
	f func()
}

func (c command) Execute() {
	c.f()
}
