## 请对你创建的 goroutine 负责

### 不要创建一个你不知道何时退出的 goroutine

请阅读下面这段代码，看看有什么问题？

> 为什么先从下面这段代码出发，是因为在之前的经验里面我们写了大量类似的代码，之前没有意识到这个问题，并且还因为这种代码出现过短暂的事故

```go
// Week03/blog/01/01.go
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func setup() {
	// 这里面有一些初始化的操作
}

func main() {
	setup()

	// 主服务
	server()

	// for debug
	pprof()

	select {}
}

func server() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		// 主服务
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Panicf("http server err: %+v", err)
			return
		}
	}()
}

func pprof() {
	// 辅助服务，监听了其他端口，这里是 pprof 服务，用于 debug
	go http.ListenAndServe(":8081", nil)
}

```

灵魂拷问来了，请问：

- 如果 `server` 是在其他包里面，如果没有特殊说明，你知道这是一个异步调用么？
- `main` 函数当中最后在哪里空转干什么？会不会存在浪费？
- 如果线上出现事故，debug 服务已经退出，你想要 debug 这时你是否很茫然？
- 如果某一天服务突然重启，你却找不到事故日志，你是否能想起这个 `8081` 端口的服务？

#### 请将选择权留给对方，不要帮别人做选择

请把是否并发的选择权交给你的调用者，而不是自己就直接悄悄的用上了 goroutine
下面这次改动将两个函数是否并发操作的选择权留给了 main 函数

```go
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func setup() {
	// 这里面有一些初始化的操作
}

func main() {
	setup()

	// for debug
	go pprof()

	// 主服务
	go server()

	select {}
}

func server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// 主服务
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panicf("http server err: %+v", err)
		return
	}
}

func pprof() {
	// 辅助服务，监听了其他端口，这里是 pprof 服务，用于 debug
	http.ListenAndServe(":8081", nil)
}
```

#### 请不要作为一个旁观者

一般情况下，不要让主进程成为一个旁观者，明明可以干活，但是最后使用了一个 `select` 在那儿空跑
感谢上一步将是否异步的选择权交给了我( `main` )，在旁边看着也怪尴尬的

```go
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func setup() {
	// 这里面有一些初始化的操作
}

func main() {
	setup()

	// for debug
	go pprof()

	// 主服务
	server()
}

func server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// 主服务
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panicf("http server err: %+v", err)
		return
	}
}

func pprof() {
	// 辅助服务，监听了其他端口，这里是 pprof 服务，用于 debug
	http.ListenAndServe(":8081", nil)
}
```

#### 不要创建一个你永远不知道什么时候会退出的 goroutine

我们再做一些改造，使用 `channel` 来控制，解释都写在代码注释里面了

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func setup() {
	// 这里面有一些初始化的操作
}

func main() {
	setup()

	// 用于监听服务退出
	done := make(chan error, 2)
	// 用于控制服务退出，传入同一个 stop，做到只要有一个服务退出了那么另外一个服务也会随之退出
	stop := make(chan struct{}, 0)
	// for debug
	go func() {
		done <- pprof(stop)
	}()

	// 主服务
	go func() {
		done <- app(stop)
	}()

	// stoped 用于判断当前 stop 的状态
	var stoped bool
	// 这里循环读取 done 这个 channel
	// 只要有一个退出了，我们就关闭 stop channel
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			log.Printf("server exit err: %+v", err)
		}

		if !stoped {
			stoped = true
			close(stop)
		}
	}
}

func app(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return server(mux, ":8080", stop)
}

func pprof(stop <-chan struct{}) error {
	// 注意这里主要是为了模拟服务意外退出，用于验证一个服务退出，其他服务同时退出的场景
	go func() {
		server(http.DefaultServeMux, ":8081", stop)
	}()

	time.Sleep(5 * time.Second)
	return fmt.Errorf("mock pprof exit")
}

// 启动一个服务
func server(handler http.Handler, addr string, stop <-chan struct{}) error {
	s := http.Server{
		Handler: handler,
		Addr:    addr,
	}

	// 这个 goroutine 我们可以控制退出，因为只要 stop 这个 channel close 或者是写入数据，这里就会退出
	// 同时因为调用了 s.Shutdown 调用之后，http 这个函数启动的 http server 也会优雅退出
	go func() {
		<-stop
		log.Printf("server will exiting, addr: %s", addr)
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}
Copy
```

我们看一下返回结果，这个代码启动 5s 之后就会退出程序

```
❯ go run ./01_goroutine/04
2020/12/08 21:49:43 server exit err: mock pprof exit
2020/12/08 21:49:43 server will exiting, addr: :8081
2020/12/08 21:49:43 server will exiting, addr: :8080
2020/12/08 21:49:43 server exit err: http: Server closed
```

### 不要创建一个永远都无法退出的 goroutine [goroutine 泄漏]

再来看下面一个例子，这也是常常会用到的操作

```Go 
func leak(w http.ResponseWriter, r *http.Request) {
	ch := make(chan bool, 0)
	go func() {
		fmt.Println("异步任务做一些操作")
		<-ch
	}()

	w.Write([]byte("will leak"))
}
```

复用一下上面的 server 代码，我们经常会写出这种类似的代码

- http 请求来了，我们启动一个 goroutine 去做一些耗时一点的工作
- 然后返回了
- 然后之前创建的那个 **goroutine 阻塞**了
- 然后就泄漏了

### 确保创建出的 goroutine 的工作已经完成

这个其实就是优雅退出的问题，我们可能启动了很多的 goroutine 去处理一些问题，但是服务退出的时候我们并没有考虑到就直接退出了。例如退出前日志没有 flush 到磁盘，我们的请求还没完全关闭，异步 worker 中还有 job 在执行等等。
我们也来看一个例子，假设现在有一个埋点服务，每次请求我们都会上报一些信息到埋点服务上

## Go 内存模型（The Go Memory Model 阅读笔记）

> 软件(*编译器*)或硬件(*CPU*)系统可以根据其对代码的分析结果，一定程度上打乱代码的执行顺序，以达到其不可告人的目的(*提高 CPU 利用率)*

所以我们在编写并发程序的时候一定要小心，然后回到我们本次的主题 Go 内存模型，就是要解决两个问题，一个是要了解谁先谁后，有个专有名词叫 `Happens Before` ，另外一个就是了解可见性的问题，我这次读取能不能看到另外一个线程的写入
接下来我们官方的文档[《The Go Memory Model》](https://golang.org/ref/mem)的思路一步一步的了解这些问题，因为官方的文档写的相对比较精炼，所以会比较难懂，我会尝试加入一些我的理解补充说明。

### 忠告

> Programs that modify data being simultaneously accessed by multiple goroutines must serialize such access.
> To serialize access, protect the data with channel operations or other synchronization primitives such as those in the [`sync`](https://golang.org/pkg/sync/) and [`sync/atomic`](https://golang.org/pkg/sync/atomic/) packages.

这个是说如果你的程序存在多个 goroutine 去访问数据的时候，**必须序列化的访问，**如何保证序列化呢？我们可以采用 channel 或者是 sync 以及 sync/atomic 下面提供的同步语义来保证

### Happens Before

#### 序

> Within a single goroutine, reads and writes must behave as if they executed in the order specified by the program. That is, compilers and processors may reorder the reads and writes executed within a single goroutine only when the reordering does not change the behavior within that goroutine as defined by the language specification. Because of this reordering, the execution order observed by one goroutine may differ from the order perceived by another. For example, if one goroutine executes a = 1; b = 2;, another might observe the updated value of b before the updated value of a.

这段话就解释了上面我们示例当中为什么会出现 `2 0` 这种情况。
这段话就是说我们在单个 goroutine 当中的编写的代码会总是按照我们编写代码的顺序来执行

- 当然这个也是符合我们的习惯的
- 当然这并不表示编译器在编译的时候不会对我们的程序进行指令重排
- 而是说只会在不影响语言规范对 goroutine 的行为定义的时候，编译器和 CPU 才会对读取和写入的顺序进行重新排序。

但是正是因为存在这种重排的情况，所以一个 goroutine 监测到的执行顺序和另外一个 goroutine 监测到的有可能不一样。就像我们最上面的这个例子一样，可能我们在 f 执行的顺序是先执行 `a = 1` 后执行 `b = 2` 但是在 g 中我们只看到了 b = 2 具体什么情况可能会导致这个呢？不要着急，我们后面还会说到

#### 编译器重排

我们来看参考文章中的一个编译器重排例子

```
X = 0
for i in range(100):
    X = 1
    print X
Copy
```

在这段代码中，X = 1 在 for 循环内部被重复赋值了 100 次，这完全没有必要，于是聪明的编译器就会帮助我们优化成下面的样子

```
X = 1
for i in range(100):
    print X
Copy
```

完美，在单个 goroutine 中并不会改变程序的执行，这时候同样会输出 100 次 1 ，并且减少了 100 次赋值操作。
但是，如果与此同时我们存在一个另外一个 goroutine 干了另外一个事情 X = 0 那么，这个输出就变的不可预知了，就有可能是 1001111101… 这种，所以回到刚开始的忠告：**这个是说如果你的程序存在多个 goroutine 去访问数据的时候，必须序列化的访问**

#### happens before 定义

> To specify the requirements of reads and writes, we define happens before, a partial order on the execution of memory operations in a Go program. If event `e1` happens before event `e2`, then we say that `e2` happens after `e1`. Also, if `e1` does not happen before `e2` and does not happen after `e2`, then we say that `e1` and `e2` happen concurrently.

这是 Happens Before 的定义，如果 `e1` 发生在 `e2` 之前，那么我们就说 `e2` 发生在 `e1` 之后，如果 `e1` 既不在 `e2` 前，也不在 `e2` 之后，那我们就说这俩是并发的

> Within a single goroutine, the happens-before order is the order expressed by the program.

这就是我们前面提到的，在单个 goroutine 当中，事件发生的顺序，就是程序所表达的顺序

> A read r of a variable `v` is allowed to observe a write `w` to `v` if both of the following hold:
>
> 1. r does not happen before `w`.
> 2. There is no other write `w'` to `v` that happens after w but before `r`.

假设我们现在有一个变量 `v`，然后只要满足下面的两个条件，那么读取操作 `r` 就可以对这个变量 `v` 的写入操作 `w` 进行监测

1. 读取操作 `r` 发生在写入操作 `w` 之后
2. 并且在 `w` 之后，`r` 之前没有其他对 `v` 的写入操作 `w'`

注意这里说的只是读取操作 r 可以对 w 进行监测，但是能不能读到呢，可能可以也可能不行

> To guarantee that a read `r` of a variable `v` observes a particular write `w` to `v`, ensure that `w` is the only write `r` is allowed to observe. That is, `r` is guaranteed to observe `w` if both of the following hold:
>
> 1. `w` happens before `r`.
> 2. Any other write to the shared variable `v` either happens before `w` or after `r`.

为确保对变量 `v` 的读取操作 `r` 能够监测到特定的对 `v` 进行写入的操作 `w`，需确保 `w` 是唯一允许被 `r` 监测的写入操作。也就是说，若以下条件均成立，则 `r` 能保证监测到 `w`：

1. `w` 发生在 `r` 之前。
2. 对共享变量 `v` 的其它任何写入操作都只能发生在 `w` 之前或 `r` 之后。

这对条件的要求比第一个条件更强，它需要确保没有其它写入操作与 `w` 或 `r` 并发。
在单个 goroutine 当中这两个条件是等价的，因为单个 goroutine 中不存在并发，在多个 goroutine 中就必须使用同步语义来确保顺序，这样才能到保证能够监测到预期的写入
**单个 goroutine 的情况**：
我们可以发现在单个 goroutine 当中，读取操作 r 总是可以读取到上一次 w 写入的值的
![image.png](https://img.lailin.xyz/image/1608372439492-359ad5bf-1b06-4f4f-ae77-84e96d9f6a7f.png)


**多个 goroutine 的情况**:
但是存在多个 goroutine 的时候这个就不一定了，r0 读到的是 哪一次写入的值呢？如果看图的话像是 w4 的，但其实不一定，因为图中的两个 goroutine 所表达的时间维度可能是不一致的，所以 r0 可能读到的是 w0 w3 w4 甚至是 w5 的结果，当然按照我们前面说的理论，读到的不可能是 w1 的结果的
![image.png](https://img.lailin.xyz/image/1608372753766-f3b66fe5-ac34-4f5e-a1b2-c74e7d3dfbc9.png)
**添加一些同步点**
如下图所示我们通过 sync 包中的一些同步语义或者是 channel 为多个 goroutine 加入了 同步点，那么这个时候对于 r1 而言，他就是晚于 w4 并且早于 w1 和 w5 执行的，所以它读取到的是写入操作是可以确定的，是 w4
![image.png](https://img.lailin.xyz/image/1608373281116-271c756e-386e-490b-aa7b-0fb2b741ed40.png)



> The initialization of variable `v` with the zero value for `v`‘s type behaves as a write in the memory model.

以变量 `v` 所属类型的零值来对 `v` 进行初始化，其表现如同在内存模型中进行的写入操作。

#### 机器字

> Reads and writes of values larger than a single machine word behave as multiple machine-word-sized operations in an unspecified order.

对大于单个机器字的值进行读取和写入，其表现如同以不确定的顺序对多个机器字大小的值进行操作。要理解这个我们首先要理解什么是机器字。
我们现在常见的还有 32 位系统和 64 位的系统，cpu 在执行一条指令的时候对于单个机器字长的的数据的写入可以保证是原子的，对于 32 位的就是 4 个字节，对于 64 位的就是 8 个字节，对于在 32 位情况下去写入一个 8 字节的数据时就需要执行两次写入操作，这两次操作之间就没有原子性，那就可能出现先写入后半部分的数据再写入前半部分，或者是写入了一半数据然后写入失败的情况。
也就是说虽然有时候我们看着仅仅只做了一次写入但是还是会有并发问题，因为它本身不是原子的

### 同步

#### 初始化

> Program initialization runs in a single goroutine, but that goroutine may create other goroutines, which run concurrently.
> If a package `p` imports package `q`, the completion of `q`‘s `init` functions happens before the start of any of `p`‘s.
> The start of the function `main.main` happens after all `init` functions have finished.

- 程序的初始化运行在单个 goroutine 中，但该 goroutine 可能会创建其它并发运行的 goroutine

- 若包 p 导入了包 q，则 q 的 init 函数会在 p 的任何函数启动前完成。

- 函数 main.main 会在所有的 init 函数结束后启动。

  注意: 在实际的应用代码中不要隐式的依赖这个启动顺序

#### goroutine 的创建

> The `go` statement that starts a new goroutine happens before the goroutine’s execution begins.

`go` 语句会在 goroutine 开始执行前启动它

#### goroutine 的销毁

> The exit of a goroutine is not guaranteed to happen before any event in the program。

goroutine 无法确保在程序中的任何事件发生之前退出

注意 [《The Go Memory Model》](https://golang.org/ref/mem)原文中还有关于 channel， 锁相关的阐述，因为篇幅原因在本文中就不多讲了，后面我们还有单独的文章详细讲 channel 和 锁 相关的使用，在强调一遍，原文一定要多看几遍

## 数据竞争(data race)

之前我们提到了很多次在多个 goroutine 对同一个变量的数据进行修改的时候会出现很多奇奇怪怪的问题，那我们有没有什么办法检测它呢，除了通过我们聪明的脑袋？

答案就是 data race tag，go 官方早在 1.1 版本就引入了数据竞争的检测工具，我们只需要在执行测试或者是编译的时候加上 `-race` 的 flag 就可以开启数据竞争的检测

```shell l
go test -race ./...
go build -race
```

不建议在生产环境 build 的时候开启数据竞争检测，因为这会带来一定的性能损失(一般内存5-10倍，执行时间2-20倍)，当然 必须要 debug 的时候除外。
建议在执行单元测试时始终开启数据竞争的检测。

### 案例一

我们来直接看一下下面的这个例子，这是来自课上的一个例子，但是我稍稍做了一些改造，源代码没有跑 10w 次这个操作，会导致看起来每次跑的结果都是差不多的，我们只需要把这个次数放大就可以发现每次结果都会不一样

#### 正常执行

```go
package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var counter int

func main() {
	// 多跑几次来看结果
	for i := 0; i < 100000; i++ {
		run()
	}
}

func run() {
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go routine(i)
	}
	wg.Wait()
	fmt.Printf("Final Counter: %d\n", counter)
}

func routine(id int) {
	for i := 0; i < 2; i++ {
		value := counter
		value++
		counter = value
	}
	wg.Done()
}
```

我执行了三次每次的结果都不一致，分别是:

```
Final Counter: 399996
Final Counter: 399989
Final Counter: 399988
Copy
```

为什么会导致这样的结果呢，是因为每一次执行的时候，我们都使用 `go routine(i)` 启动了两个 goroutine，但是我们并没有控制它的执行顺序，那就有好几种可能了，我这里描述两种情况

1. 执行一次 `run()` , `counter + 4` 这种情况下，第二个 goroutine 开始执行时，拿到了第一个 goroutine 的执行结果，也就是 `value := counter` 这一步时，value = 2
2. 执行一次 `run()` , `counter + 2` 这种情况下，第二个 goroutine 开始执行时，没有拿到了第一个 goroutine 的执行结果，也就是 `value := counter` 这一步时，counter 还是零值，这时候 value = 0

当然由于种种不确定性，所有肯定不止这两种情况，但是这个不是本文讨论的重点，具体的原因可以结合上一篇文章 [Week03: Go 并发编程(二) Go 内存模型](https://lailin.xyz/post/go-training-week3-go-memory-model.html) 进行思考

#### data race 执行

可以发现，写出这种代码时上线后如果出现 bug 会非常难定位，因为你不知道到底是哪里出现了问题，所以我们就要在测试阶段就结合 data race 工具提前发现问题。
我们执行以下命令

```shell
go run -race ./main.go
```

会发现结果会所有的都输出， `data race` 的报错信息，我们已经看不到了，因为终端的打印的太长了，可以发现的是，最后打印出发现了一处 data race 并且推出码为 `66`

```sh
Final Counter: 399956
Final Counter: 399960
Found 1 data race(s)
exit status 66
```

#### data race 配置

问题来了，我们有没有什么办法可以立即知道 data race 的报错呢？
答案就在官方的文档当中，我们可以通过设置 `GORACE` 环境变量，来控制 data race 的行为， 格式如下:

```sh
GORACE="option1=val1 option2=val2"
```

可选配置:

| **配置**          | **默认值** | **说明**                                                     |
| ----------------- | ---------- | ------------------------------------------------------------ |
| log_path          | stderr     | 日志文件的路径，除了文件路径外支持 stderr, stdout 这两个特殊值 |
| exitcode          | 66         | 退出码                                                       |
| strip_path_prefix | “”         | 从日志中的文件信息里面去除相关的前缀，可以去除本地信息，同时会更好看 |
| history_size      | 1          | per-goroutine 内存访问历史记录为 32K * 2 ** history_size，增加这个可以避免出现堆栈还原失败的错误，但是增加这个会导致使用的内存也跟着增加 |
| halt_on_error     | 0          | 用来控制第一个数据竞争错误出现后是否立即退出                 |
| atexit_sleep_ms   | 100        | 用来控制 main 退出之前 sleep 的时间                          |

有了这个背景知识后就很简单了，在我们这个场景我们可以控制发现数据竞争后直接退出

```sh
GORACE="halt_on_error=1 strip_path_prefix=/home/ll/project/Go-000/Week03/blog/03_sync/01_data_race" go run -race ./main.go
```

重新执行后我们的结果

```shell
==================
WARNING: DATA RACE
Read at 0x00000064a9c0 by goroutine 7:
  main.routine()
      /main.go:29 +0x47

Previous write at 0x00000064a9c0 by goroutine 8:
  main.routine()
      /main.go:31 +0x64

Goroutine 7 (running) created at:
  main.run()
      /main.go:21 +0x75
  main.main()
      /main.go:14 +0x38

Goroutine 8 (finished) created at:
  main.run()
      /main.go:21 +0x75
  main.main()
      /main.go:14 +0x38
==================
exit status 66

```

这个结果非常清晰的告诉了我们在 29 行这个地方我们有一个 goroutine 在读取数据，但是呢，在 31 行这个地方又有一个 goroutine 在写入，所以产生了数据竞争。
然后下面分别说明这两个 goroutine 是什么时候创建的，已经当前是否在运行当中。

### 典型案例

接来下我们再来看一些典型案例，这些案例都来自 go 官方的文档 [Data Race Detector](https://golang.org/doc/articles/race_detector.html)，这些也是初学者很容易犯的错误

#### 案例二 在循环中启动 goroutine 引用临时变量

```go
func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i) // Not the 'i' you are looking for.
			wg.Done()
		}()
	}
	wg.Wait()
}
```

如果你去找一些 go 的面试题，肯定能找到类似的例子，然后会问你这里会输出什么？
常见的答案就是会输出 5 个 5，因为在 for 循环的 i++ 会执行的快一些，所以在最后打印的结果都是 5
这个答案不能说不对，因为真的执行的话大概率也是这个结果，但是不全
因为这里本质上是有数据竞争，在新启动的 goroutine 当中读取 i 的值，在 main 中写入，导致出现了 data race，这个结果应该是不可预知的，因为我们不能假定 goroutine 中 print 就一定比外面的 i++ 慢，习惯性的做这种假设在并发编程中是很有可能会出问题的

```go
func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```

这个要修改也很简单，只需要将 i 作为参数传入即可，这样每个 goroutine 拿到的都是拷贝后的数据

#### 案例三 一不小心就把变量共享了

```go
package main

import "os"

func main() {
	ParallelWrite([]byte("xxx"))
}

// ParallelWrite writes data to file1 and file2, returns the errors.
func ParallelWrite(data []byte) chan error {
	res := make(chan error, 2)
	f1, err := os.Create("/tmp/file1")
	if err != nil {
		res <- err
	} else {
		go func() {
			// This err is shared with the main goroutine,
			// so the write races with the write below.
			_, err = f1.Write(data)
			res <- err
			f1.Close()
		}()
	}
	f2, err := os.Create("/tmp/file2") // The second conflicting write to err.
	if err != nil {
		res <- err
	} else {
		go func() {
			_, err = f2.Write(data)
			res <- err
			f2.Close()
		}()
	}
	return res
}
```

我们使用 `go run -race main.go` 执行，可以发现这里报错的地方是，19 行和 24 行，有 data race，这里主要是因为共享了 err 这个变量

```shell
==================
WARNING: DATA RACE
Write at 0x00c0000a01a0 by goroutine 7:
  main.ParallelWrite.func1()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:19 +0x94

Previous write at 0x00c0000a01a0 by main goroutine:
  main.ParallelWrite()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:24 +0x1dd
  main.main()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:6 +0x84

Goroutine 7 (running) created at:
  main.ParallelWrite()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:16 +0x336
  main.main()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:6 +0x84
==================
Found 1 data race(s)
exit status 66
```

修改的话只需要在两个 goroutine 中使用新的临时变量就行了

```go
...
_, err := f1.Write(data)
...
_, err := f2.Write(data)
...
```

细心的同学可能会有这个疑问，在 24 行不也是重新赋值了么，为什么在这里会和 19 行产生 data race 呢？
这是由于 go 的语法规则导致的，我们在初始化变量的时候如果在同一个作用域下，如下方代码，这里使用的 err 其实是同一个变量，只是 f1 f2 不同，具体可以看 [effective go 当中 Redeclaration and reassignment](https://golang.org/doc/effective_go.html#redeclaration) 的内容

```go
f1, err := os.Create("a")
f2, err := os.Create("b")
```

