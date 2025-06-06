package simple_worker_pool

import (
	"sync"

	sn "github.com/Vladimir220/sync_notification"
)

type IWorkerPool interface {
	SetNumOfWorkers(n uint32) error
	GetNumOfWorkers() uint32
	Stop()
	SetCommand(com ICommand)
}

type workerPool struct {
	numOfWorkers uint32
	com          ICommand
	stopSignal   sn.ISyncNotification
	mu           *sync.RWMutex
}

func CreateWorkerPool(com ICommand) IWorkerPool {
	return &workerPool{
		com:        com,
		stopSignal: sn.CreateSyncNotification(),
		mu:         &sync.RWMutex{},
	}
}

func (wp *workerPool) SetNumOfWorkers(n uint32) error {
	delta := int64(n) - int64(wp.numOfWorkers)
	if delta < 0 {
		wp.stopSignal.Signal(uint32(delta * -1))
	}

	wp.mu.Lock()
	defer wp.mu.Unlock()

	if delta > 0 {
		for range delta {
			err := wp.startWorker()
			if err != nil {
				return err
			}
		}
	}

	wp.numOfWorkers = n

	return nil
}

func (wp workerPool) GetNumOfWorkers() uint32 {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.numOfWorkers
}

func (wp *workerPool) Stop() {
	wp.stopSignal.Broadcast()
	wp.mu.Lock()
	wp.numOfWorkers = 0
	wp.mu.Unlock()
}

func (wp *workerPool) SetCommand(com ICommand) {
	wp.mu.Lock()
	wp.com = com
	wp.mu.Unlock()
}

func (wp *workerPool) startWorker() error {
	waiting, err := wp.stopSignal.GetWaiting()
	if err != nil {
		return err
	}

	go func(worker IWorker) {
		worker.Start()
	}(CreateWorker(&wp.com, wp.mu, waiting))

	return nil
}
