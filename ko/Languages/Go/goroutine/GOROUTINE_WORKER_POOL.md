# Goroutine Worker Pool

## Goroutine Worker Pool이란?

- Goroutine Work Pool은 go의 Concurrency(동시성) 작업을 효율적이고 통제 가능하게 처리하기 위한 디자인 패턴

### Worker Pool이란

- Worker Pool은 정해진 수의 worker를 사용해 큐에 있는 다수의 task(작업)을 실행함으로써 Concurrency를 달성하기 위한 디자인 패턴

- task 실행의 concurrency를 달성하기 위한 디자인 패턴이라는 점에서 Thread Pool과 유사하지만, OS 스레드를 관리하는 Thread Pool과 달리 Worker Pool은 Goroutine을 관리

## 왜 필요한가?

### 1. Resource Management

- Goroutine은 Initial size가 2KB 밖에 안 될만큼 가볍고, Concurrency를 구현하기에 간편함
- 하지만 우리 서버의 물리적 리소스는 무한하지 않기 때문에 무분별한 Goroutine 사용은 성능 저하의 원인이 될 수 있음

### 2. Performance Optimization

- 모든 task마다 Goroutine을 생성하면, 스케줄링에 오버헤드가 발생하여 CPU 효율이 떨어짐
- worker 수를 고정함으로써 스케줄링 오버헤드를 줄여 성능을 최적화할 수 있음

### 3. Achieving stable throughput

- HTTP 요청이나 메시지 큐와 같이 입력 속도가 시스템 처리 속도보다 빠를 수 있는 상황에서 시스템이 요청을 무시하거나 느리게 응답하는 대신, 대기열을 만들어 안정적인 속도로 작업을 처리하도록 보장 가능

## Goroutine Worker Pool은 필수적인가?
위 설명을 보면, Goroutine Worker Pool은 go 프로젝트의 필수 디자인 패턴으로 보임
하지만 다음과 같이 프로젝트마다 Goroutine Worker Pool이 필요한 경우와 아닌 경우가 있음

### Goroutine Worker Pool이 필요한 경우

- 작업 부하가 무한정(unbounded)이거나 대량인 경우
    - Worker Pool이 Goroutine의 무제한 증가를 방지하여, 메모리 소진, Garbage Collection 부하, 예측 불가능한 성능 저하를 막음
- 무한한 동시성이 자원 포화(Resource Saturation) 위험을 초래할 경우
    - 동시 워커 수를 제한하면 부하가 걸리는 상황에서 리소스를 독점하는 것을 방지함
- 안정성을 위해 예측 가능한 병렬 처리가 필요할 경우
    - Concurrency(동시성)을 제한함으로써 트래픽이 몰릴 때에도 일관성 있는 시스템 동작을 유지함
- task가 비교적 균일하고 큐(queue) 처리에 적합할 경우
    - task cost가 일관적이면, 고정된 Worker Pool 크기는 최소한의 오버헤드로 효율적인 스케줄링을 제공함
    - 즉, 조정 없이도 좋은 throughput(처리량)을 제공할 수 있음

### Goroutine Worker Pool이 적합하지 않은 경우

- 각 task가 최소한의 지연 시간(Latency)으로 즉시 처리되어야 할 경우
    - Worker Pool의 queuing은 당연히 지연을 발생시킴
    - 지연 시간에 민감한 task의 경우, Goroutine을 직접 실행하는 것이 나을 수 있음
- Task의 부하가 낮을 경우
    - 워크로드가 가벼운 task일 경우, Worker Pool을 관리하는 오버헤드가 더 클 수 있음
- Task의 부하가 예측 가능할 경우
    - task의 워크로드가 예측 가능하고 제한되어 있는 경우, Goroutine을 직접 실행하는 것이 코드를 더 단순하게 만듦


## Basic Implementation

Goroutine Worker Pool을 구현하기 위해서는 크게 3가지 컴포넌트가 필요함
1. **Job Queue(Channel)**: 처리해야 할 task를 버퍼링하는 queue(channel)
2. **Workers**: 고정된 수로 생성되어 영구적으로 대기하는 goroutine. 이 worker들은 Job Queue에서 task를 읽고 처리한 후 다음 task를 기다림
3. **Result Channel**: worker들의 작업 결과를 메인 프로그램으로 보내기 위한 channel(선택사항)

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	maxWorkers  int
	queuedTask  chan func()
	taskCounter *sync.WaitGroup
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		maxWorkers:  maxWorkers,
		queuedTask:  make(chan func(), maxWorkers*2),
		taskCounter: &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) ExecuteWorker(workerId int) {
	for task := range wp.queuedTask {
		fmt.Printf("[Worker %d] Task Start\n", workerId)
		task()
		fmt.Printf("[Worker %d] Task Finished\n", workerId)
		wp.taskCounter.Done()
	}
}

func (wp *WorkerPool) Run() {
	for i := 1; i <= wp.maxWorkers; i++ {
		go wp.ExecuteWorker(i)
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.taskCounter.Add(1)
	wp.queuedTask <- task
}

func (wp *WorkerPool) Wait() {
	wp.taskCounter.Wait()
}

func main() {
	const totalTasks = 10
	const maxConcurrency = 3

	startTime := time.Now()

	pool := NewWorkerPool(maxConcurrency)
	pool.Run()

	for i := 1; i <= totalTasks; i++ {
		taskId := i
		taskFunc := func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("[Task %d] Processed after 1 second\n", taskId)
		}
		pool.Submit(taskFunc)
		fmt.Printf(" -> Submit Task %d\n", taskId)
	}

	pool.Wait()

	elapsed := time.Since(startTime)
	fmt.Println("---")
	fmt.Printf("Finish. Total Tasks: %d, Max Workers: %d\n", totalTasks, maxConcurrency)
	fmt.Print("Elapsed: %.2f\n", elapsed.Seconds())
}
```

- `WorkPool`
    - Goroutine Worker Pool을 WorkerPool이란 구조체로 추상화 함
    - `maxWorkers int`
        - worker pool에서 가용할 수 있는 최대 worker(goroutine) 수를 정의
    - `queuedTask chan func()`
        - 외부에서 task(작업)를 받기 위한 channel
        - worker가 이 channel에서 task를 꺼내 goroutine으로 실행함
        - 실제 Job Queue의 역할을 함
- `ExecuteWorker(workerId int)`
    - 실제 task를 실행하기 위한 함수
    - `queuedTask` 채널로 들어오는 task를 처리함
- `Run()`
    - `go wp.ExecuteWorker(i)`로 worker 수 만큼 worker(goroutine)을 생성함
- `Submit(task func())`
    - 외부에서 task를 `queuedTask`로 넣어주기 위한 함수


## References

- https://goperf.dev/01-common-patterns/worker-pool/
- https://syafdia.medium.com/go-concurrency-pattern-worker-pool-a437117025b1
- https://erfansahaf.medium.com/managing-goroutines-with-gouroutine-pooling-in-go-9b3596e23225
- https://dev.to/zeedu_dev/worker-pool-design-pattern-explanation-3kil