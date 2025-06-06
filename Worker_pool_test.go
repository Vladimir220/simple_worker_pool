package simple_worker_pool

import (
	"sync"
	"testing"
	"time"
)

func checkOutput(t *testing.T, output *chan string, isOutputExpected bool) {
	select {
	case <-time.After(time.Second):
		if isOutputExpected {
			t.Fatal("a value was expected in the channel, but it was not received")
		}
	case str := <-*output:
		if isOutputExpected {
			t.Log("output:", str)
		} else {
			t.Fatal("an unexpected message was received in the channel", str)
		}
	}
}

func getCommandFunc(param string, input, output *chan string, mu *sync.Mutex) func() {
	return func() {
		mu.Lock()
		str, ok := <-*input
		mu.Unlock()
		if ok {
			*output <- param + " " + str
			time.Sleep(3 * time.Second)
		}
	}
}

func TestWorkerPool(t *testing.T) {
	input := make(chan string)
	output := make(chan string)
	mu := &sync.Mutex{}

	command := CreateCommand(getCommandFunc("hello", &input, &output, mu))

	wp := CreateWorkerPool(command)

	// проверяем, что workers вообще работают
	t.Log("checking the work of workers ---------------------------")
	wp.SetNumOfWorkers(2)
	time.Sleep(time.Second)

	input <- "world"
	input <- "world"

	for range 2 {
		checkOutput(t, &output, true)
	}

	checkOutput(t, &output, false)

	// проверяем работу при уменьшении количества workers
	t.Log("сhecking the reduction in the number of workers ---------------------------")
	close(input)
	wp.SetNumOfWorkers(0)
	time.Sleep(time.Second * 3)
	mu.Lock()
	input = make(chan string)
	mu.Unlock()

	go func() { input <- "world" }()

	checkOutput(t, &output, false)

	if wp.GetNumOfWorkers() != 0 {
		t.Fatal("the number of workers is greater than expected:", wp.GetNumOfWorkers())
	}
	<-input

	// проверяем работу при увеличении количества workers
	t.Log("сhecking the increase in the number of workers ---------------------------")
	close(input)
	wp.SetNumOfWorkers(2)
	time.Sleep(time.Second * 3)
	mu.Lock()
	input = make(chan string)
	mu.Unlock()

	input <- "world"
	input <- "world"
	for range 2 {
		checkOutput(t, &output, true)
	}

	checkOutput(t, &output, false)

	if wp.GetNumOfWorkers() != 2 {
		t.Fatal("the number of workers is not equal to 2:", wp.GetNumOfWorkers())
	}

	// проверяем работу wg.Stop()
	t.Log("wg.Stop() check ---------------------------")
	close(input)
	wp.Stop()
	time.Sleep(time.Second * 3)
	mu.Lock()
	input = make(chan string)
	mu.Unlock()

	go func() { input <- "world" }()

	checkOutput(t, &output, false)

	if wp.GetNumOfWorkers() != 0 {
		t.Fatal("the number of workers is greater than expected:", wp.GetNumOfWorkers())
	}
	<-input

	// проверяем работу при смене команды пула
	t.Log("command change check ---------------------------")
	wp.SetNumOfWorkers(1)
	time.Sleep(time.Second)
	close(input)
	command = CreateCommand(getCommandFunc("goodbye", &input, &output, mu))
	wp.SetCommand(command)
	time.Sleep(time.Second * 3)
	mu.Lock()
	input = make(chan string)
	mu.Unlock()

	input <- "world"

	select {
	case <-time.After(time.Second):
		t.Fatal("a value was expected in the channel, but it was not received")
	case str := <-output:
		if str == "goodbye world" {
			t.Log("output:", str)
		} else {
			t.Fatal("unexpected data:", str)
		}
	}

	wp.Stop()
	close(input)
	time.Sleep(time.Second * 3)
}
