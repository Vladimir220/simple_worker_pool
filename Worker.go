package simple_worker_pool

import (
	"sync"

	sn "github.com/Vladimir220/sync_notification"
)

type IWorker interface {
	Start()
}

func CreateWorker(comPtr *ICommand, mu *sync.RWMutex, waiting sn.IWaiting) IWorker {
	return worker{
		waiting: waiting,
		comPtr:  comPtr,
		mu:      mu,
	}
}

type worker struct {
	waiting sn.IWaiting
	comPtr  *ICommand
	mu      *sync.RWMutex
}

func (w worker) Start() {
	for {
		select {
		case <-w.waiting.GetSignalChan():
			w.waiting.Done()
			return
		default:
		}
		w.mu.RLock()

		(*w.comPtr).Execute()

		w.mu.RUnlock()
	}
}
