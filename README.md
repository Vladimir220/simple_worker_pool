# Simple worker pool #
## Содержание ##
- [Установка](#установка)
- [Интерфейсы](#интерфейсы)
- [Конструкторы](#конструкторы)
- [Используемая библиотека собственного производства](#используемая-библиотека-собственного-производства)
- [Предупреждение](#предупреждение)
---

## Установка ##
```
go get github.com/Vladimir220/simple_worker_pool
```

## Интерфейсы ##
```
type IWorkerPool interface {
	SetNumOfWorkers(n uint32) error
	GetNumOfWorkers() uint32
	Stop()
	SetCommand(com ICommand)
}
```
```
type ICommand interface {
	Execute()
}
```
```
type IWorker interface {
	Start()
}
```

## Конструкторы ##
```
func CreateWorkerPool(com ICommand) IWorkerPool
```
```
func CreateCommand(f func()) ICommand
```
```
func CreateWorker(comPtr *ICommand, mu *sync.RWMutex, waiting sn.IWaiting) IWorker
```

Для большей гибкости можете реализовать свою структуру для интерфейса ```ICommand``` и передать её в ```CreateWorkerPool()```.

## Используемая библиотека собственного производства ##
Для управления Workers был написан инструмент синхнонизирующих уведомлений [sync_notification](https://github.com/Vladimir220/sync_notification/tree/v0.1.3). 

Этот инструмент похож на ```sync.Cond```, но у него нет работы с мьютексом и нет метода ```Wait()```. Вместо метода ```Wait()``` у нас канал ожидания, получаемый из метода интерфейса ```IWaiting```. Это удобно для прослушивания в операторе ```select```.

## Предупреждение ##
Следите, чтобы функции, установленные в Command, не зависали перед выполнением WorkerPool-операций (например, на чтении или записи в канал), иначе WorkerPool будет бесконечно ожидать завершения установленной логики. 

**Тукущая версия WorkerPool реалезует только кооперативную многозадачность.**