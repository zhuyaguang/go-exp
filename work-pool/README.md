### 1.Demo on creating worker pool in GoLang

> 通过 chan queue 作为任务队列，work 数组作为并发数量，通过StopChan 停止任务
>
> [实现代码](https://github.com/zhuyaguang/go-exp/tree/main/work-pool/demo1)

* [一个 demo 学会 workerPool](https://mp.weixin.qq.com/s/YCl7r7l3Ty3wbnImVWRLxg)

### 2.The Case For A Go Worker Pool

![A visualization of a worker pool: few workers working many work items.](https://brandur.org/assets/images/go-worker-pool/worker-pool.svg)





#### 1. Base case

~~~ go
package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 9; a++ {
		<-results
	}
}

~~~



#### 2.增加错误处理

> 该版本有bug，如果设置 errors := make(chan error, 1)  会产生死锁。

~~~go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup, jobs <-chan int, results chan<- int, errors chan<- error) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)

		if j%2 == 0 {
			results <- j * 2
		} else {
			errors <- fmt.Errorf("error on job %v", j)
		}
		wg.Done()
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	errors := make(chan error, 100)

	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		go worker(w, &wg, jobs, results, errors)
	}

	for j := 1; j <= 9; j++ {
		jobs <- j
		wg.Add(1)
	}
	close(jobs)

	wg.Wait()

	select {
	case err := <-errors:
		fmt.Println("finished with error:", err.Error())
	default:
	}
}
~~~



#### 3.更健壮版本(demo2/case1.go)

原文地址 ：[Go Work Pool 几个典型例子](https://brandur.org/go-worker-pool)

> 同例子1类似，一个task数组存放任务并同步到 task channel中，一个pool确定并发的数量。



### 3.workPool

[Go语言的并发与WorkerPool - 第一部分](https://mp.weixin.qq.com/s/5pQS82nE9ivF6NjaXFsolQ)

> 串行处理耗时太大--->无限并发处理内存消耗大-->work pool 平衡

[Go语言的并发与WorkerPool - 第二部分](https://mp.weixin.qq.com/s?__biz=MzI2MDA1MTcxMg==&mid=2648468414&idx=1&sn=8efed31baa411f2e63e4fe043f207c41&chksm=f2474dd1c530c4c71f94dda44bb97201164df4a9730a5045534cdf354b54b096321e1f1b91a7&cur_album_id=1506050738668486658&scene=189#rd)

> 第一部分的升级版，增加了错误处理、退出信号 ，同例子1、例子2类似，比较特别的是增加了长连接任务。

源码GitHub：https://github.com/Joker666/goworkerpool



**1.以上只针对 生产者-消费者 同类型任务模型，如果有不同类型的task，该怎么办？**

> 1.可以通过统一grpc方法，由一个双流方法发送接收一个统一的message结构体，屏蔽了不同方法的差异。（这样不同的算法就可以放在一个task queue里面，不用放在不同的queue channel里面）
>
> 再封装一个分发的接口，负责发送到不同的后端算法，避免A任务发送到B后端。

~~~go
type Task struct {
	Err  error
	Data interface{}
	f    func(interface{}) error
}
~~~

通过在Task结构体里面，定义一个成员函数，可以让不同的task做不同的事情。

**2.如果worker处理任务方式不一样，怎么办？**

> task定义一个 f 参数，参数传入一个func ，对应不同的worker。

**3.如果生产者速度，大于消费者速度，造成了，产品堆积，超过了channel的缓冲区，怎么办？**

1.加入grpc的主动健康检测机制，服务不可用会停止生产  **生产者限流**

2.task buffer满了，会阻塞传入。 

3.采用redis消息队列处理，redis会把未处理的任务持久化。然后一个个消费。如果redis消息堆积，可以定义消息长度，来丢弃消息：中间件丢弃旧消息，只保留固定长度的新消息

**4.如果生产者(worker)数量是动态的，怎么办？**

> 统一message结构体后，任务就没有类型的区分，worker也没有类型之分，worker的数量就是并发量，可以自定义设置。



### 4.任务队列

1.taskq （解决高并发场景）

>  Golang asynchronous task/job queue with Redis, SQS, IronMQ, and in-memory backends
>
>  用 redis、SQS、IronMQ、缓存作为后端的 Golang 异步 task/job 队列
>
> https://github.com/vmihailenco/taskq

2.work pool 和 task queue 结合

https://gist.github.com/harlow/dbcd639cf8d396a2ab73



3.参考链接

[Golang 任务队列策略 -- 读《JOB QUEUES IN GO》](https://blog.csdn.net/zhizhengguan/article/details/107358568)

[Golang队列中间件开发总结](https://blog.csdn.net/qq_30145355/article/details/82322238?utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EsearchFromBaidu%7Edefault-4.pc_relevant_baidujshouduan&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EsearchFromBaidu%7Edefault-4.pc_relevant_baidujshouduan)

[golang中job队列实现方法](https://blog.csdn.net/wdy_yx/article/details/78964267)

https://twinnation.org/articles/39/go-concurrency-goroutines-worker-pools-and-throttling-made-simple



### 企业级 goroutine 池

> ants是一个高性能的 goroutine 池，实现了对大规模 goroutine 的调度管理、goroutine 复用，允许使用者在开发并发程序的时候限制 goroutine 数量，复用资源，达到更高效执行任务的效果。

[ants github 地址](https://github.com/panjf2000/ants/blob/master/README_ZH.md)

[Goroutine 并发调度模型深度解析之手撸一个高性能 Goroutine 池](https://www.infoq.cn/article/XF6v3Vapqsqt17FuTVst)