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




### 4.企业级 goroutine 池

> ants是一个高性能的 goroutine 池，实现了对大规模 goroutine 的调度管理、goroutine 复用，允许使用者在开发并发程序的时候限制 goroutine 数量，复用资源，达到更高效执行任务的效果。

[ants github 地址](https://github.com/panjf2000/ants/blob/master/README_ZH.md)

[Goroutine 并发调度模型深度解析之手撸一个高性能 Goroutine 池](https://www.infoq.cn/article/XF6v3Vapqsqt17FuTVst)



**1.以上只针对 生产者-消费者 同类型任务模型，如果有不同类型的worker-task，该怎么办？**

**2.如果生产者速度，大于消费者速度，造成了，产品堆积，超过了channel的缓冲区，怎么办？**

**3.如果生产者(worker)数量是动态的，怎么办？**



### 5.任务队列

1.taskq

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



 