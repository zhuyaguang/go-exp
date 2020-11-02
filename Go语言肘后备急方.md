1. 下面代码输出什么？请简要说明。

![image-20200914161053194](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200914161053194.png)

参考答案及解析：都输出 nil。知识点：变量的作用域。if作用域内的err变量只能作用到if语句内，外面的err无法控制。

2. 下面说法正确的是

- A. Go 语言中，声明的常量未使用会报错；
- B. cap() 函数适用于 array、slice、map 和 channel;
- C. 空指针解析会触发异常；
- D. 从一个已经关闭的 channel 接收数据，如果缓冲区中为空，则返回一个零值；

参考答案及解析：CD。A.声明的常量未使用不会报错；B.cap() 函数不适用 map。



3.下面的代码输出什么？

![image-20200914161617793](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200914161617793.png)

参考答案及解析：D。iota 的使用。



4. Go 语言中中大多数数据类型都可以转化为有效的 JSON 文本，下面几种类型除外。

- A. 指针
- B. channel
- C. complex
- D. 函数

参考答案及解析：BCD。

5. 下面代码输出什么？如果想要代码输出 10，应该如何修改？

![image-20200914161924977](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200914161924977.png)

参考答案及解析：输出 1。知识点：并发、引用。修复代码如下：加传入参数i

6.下面的代码输出什么？

![image-20200914165637830](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200914165637830.png)

参考答案及解析：2 1 3 4。知识点：可变函数。切片作为参数传入可变函数时不会创建新的切片.

7.下面代码输出什么？

```go
func main() {
  ns := []int{010: 200, 005: 100}
  print(len(ns))
}
```

参考答案及解析：9。Go 语言中，0x 开头表示 十六进制；0 开头表示八进制。

8.关于 const 常量定义，下面正确的使用方式是？

![image-20200915100407955](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200915100407955.png)

参考答案及解析：ABD。

9.下面的代码可以随机输出大小写字母，尝试在 A 处添加一行代码使得字母先按大写再按小写的顺序输出。

![image-20200915105621407](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200915105621407.png)

10.下面代码有什么问题？

![image-20200915110525794](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200915110525794.png)

参考答案及解析：panic。协程开启还未来得及执行，chan 就已经 close() ，往已经关闭的 chan 写数据会 panic。

11.下面的代码输出什么？

![image-20200915152242098](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200915152242098.png)

参考答案及解析：people:{}。按照 go 的语法，小写开头的方法、属性或 struct 是私有的，同样，在 json 解码或转码的时候也无法实现私有属性的转换。

这段代码是无法正常得到 People 的 name 值的。而且，私有属性 name 也不应该加 json 的标签。

12.下面代码输出什么？请简要说明。

![image-20200916145153777](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200916145153777.png)

参考答案及解析：C。append() 并不是并发安全的，有兴趣的同学可以尝试用锁去解决这个问题。

13.下面代码输出什么？请简要说明。

![image-20200916150203799](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200916150203799.png)

参考答案即解析：C。程序执行到第二个 groutine 时，ch 还未初始化，导致第二个 goroutine 阻塞。需要注意的是第一个 goroutine 不会阻塞。

14.下面代码输出什么？请简要说明。

![image-20200916150615830](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200916150615830.png)

参考答案及解析：B。能正确输出，不过主协程会阻塞 f() 函数的执行。

15.关于同步锁，下面说法正确的是？

- A. 当一个 goroutine 获得了 Mutex 后，其他 goroutine 就只能乖乖的等待，除非该 goroutine 释放这个 Mutex；

- B. RWMutex 在读锁占用的情况下，会阻止写，但不阻止读；

- C. RWMutex 在写锁占用情况下，会阻止任何其他 goroutine（无论读和写）进来，整个锁相当于由该 goroutine 独占；

- D. Lock() 操作需要保证有 Unlock() 或 RUnlock() 调用与之对应；

  ABD

16.

![image-20200916161645974](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200916161645974.png)

参考答案及解析：输出 0 40 50。知识点：defer 语句与返回值。函数的 return value 不是原子操作，而是在编译器中分解为两部分：返回值赋值 和 return。可以细看《[5 年 Gopher 都不知道的 defer 细节，你别再掉进坑里！](http://mp.weixin.qq.com/s?__biz=MzI2MDA1MTcxMg==&mid=2648466918&idx=2&sn=151a8135f22563b7b97bf01ff480497b&chksm=f2474389c530ca9f3dc2ae1124e4e5ed3db4c45096924265bccfcb8908a829b9207b0dd26047&scene=21#wechat_redirect)》

17.

![image-20200916164401832](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200916164401832.png)



18.

![image-20200917145343740](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200917145343740.png)

参考答案及解析：1。知识点：变量的作用域。注意 for 语句的变量 a 是重新声明，它的作用范围只在 for 语句范围内。

19.

![image-20200917150238182](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200917150238182.png)

![](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200917150259495.png)

20.

![image-20200921141952146](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200921141952146.png)

![image-20200921142014208](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200921142014208.png)

21.

![image-20200921142117652](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200921142117652.png)

![image-20200921142148425](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200921142148425.png)

22.

![image-20200921143549861](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200921143549861.png)

参考答案及解析：编译不通过。当使用 type 声明一个新类型，它不会继承原有类型的方法集。

23.

![image-20200922140338873](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200922140338873.png)

输出：

[0 11 12]
[21 12 13]



24.

![image-20200922141319080](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200922141319080.png)

![image-20200922141340253](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200922141340253.png)

25.

![image-20200922142245322](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200922142245322.png)

参考答案及解析：A。Go语言的内存回收机制规定，只要有一个指针指向引用一个变量，那么这个变量就不会被释放（内存逃逸），因此在 Go 语言中返回函数参数或临时变量是安全的。

26.

![image-20200922142437012](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200922142437012.png)

参考答案及解析：AD。全局变量要定义在函数之外，而在函数之外定义的变量只能用 var 定义。短变量声明 := 只能用于函数之内。

27.

![image-20200923092839735](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200923092839735.png)

参考答案及解析：可以编译通过，输出：true。知识点：Go 代码断行规则。注意

28.

![image-20200923111751244](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200923111751244.png)

参考答案及解析：132。这一题有两点需要注意：1.Add() 方法的返回值依然是指针类型 *Slice，所以可以循环调用方法 Add()；2.defer 函数的参数（包括接收者）是在 defer 语句出现的位置做计算的，而不是在函数执行的时候计算的，所以 s.Add(1) 会先于 s.Add(3) 执行。

29.

![image-20200924093824427](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200924093824427.png)

参考答案及解析：21。recover() 必须在 defer() 函数中调用才有效，所以第 9 行代码捕获是无效的。在调用 defer() 时，便会计算函数的参数并压入栈中，所以执行第 6 行代码时，此时便会捕获 panic(2)；此后的 panic(1)，会被上一层的 recover() 捕获。所以输出 21。

30.

![image-20200924094112246](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200924094112246.png)

AD

31.

![image-20200924095552807](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200924095552807.png)

参考答案及解析：4。知识点：多重赋值。

多重赋值分为两个步骤，有先后顺序：

- 计算等号左边的索引表达式和取址表达式，接着计算等号右边的表达式；
- 赋值；

所以本例，会先计算 s[k]，等号右边是两个表达式是常量，所以赋值运算等同于`k, s[1] = 0, 3`。

![image-20200924101856512](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200924101856512.png)

参考答案及解析：第 8 行。因为两个比较值的动态类型为同一个不可比较类型。

![image-20200924103105190](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200924103105190.png)

参考答案及解析：321。第一次循环，写操作已经准备好，执行 o(3)，输出 3；第二次，读操作准备好，执行 o(2)，输出 2 并将 c 赋值为 nil；第三次，由于 c 为 nil，走的是 default 分支，输出 1。

![image-20200924110249937](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200924110249937.png)

参考答案及解析：编译错误

```
invalid operation: fn1 != fn2 (func can only be compared to nil)
```

![image-20200924165331569](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200924165331569.png)

![image-20200924165355087](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200924165355087.png)

![image-20200925160229779](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200925160229779.png)

参考答案及解析：ABC

![image-20200925160257078](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200925160257078.png)

参考答案及解析：有方向的 channel 不可以被关闭。

![image-20200925162231198](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20200925162231198.png)

参考答案及解析：编译报错`cannot take the address of i`。知识点：常量。常量不同于变量的在运行期分配内存，常量通常会被编译器在预处理阶段直接展开，作为指令数据使用，所以常量无法寻址。

![image-20201012102138285](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20201012102138285.png)

```
r =  [1 12 13 4 5]
a =  [1 12 13 4 5]
```

 a 是一个切片，那切片是怎么实现的呢？切片在 go 的内部结构有一个指向底层数组的指针，当 range 表达式发生复制时，副本的指针依旧指向原底层数组，所以对切片的修改都会反应到底层数组上，所以通过 v 可以获得修改后的数组元素。

![image-20201012102518848](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20201012102518848.png)

可变函数、append()操作。Go 提供的语法糖`...`，可以将 slice 传进可变函数，不会创建新的切片。第一次调用 change() 时，append() 操作使切片底层数组发生了扩容，原 slice 的底层数组不会改变；第二次调用change() 函数时，使用了操作符`[i,j]`获得一个新的切片，假定为 slice1，它的底层数组和原切片底层数组是重合的，不过 slice1 的长度、容量分别是 2、5，所以在 change() 函数中对 slice1 底层数组的修改会影响到原切片。

![image-20201012103302304](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20201012103302304.png)

for range 使用短变量声明(:=)的形式迭代变量，需要注意的是，变量 i、v 在每次循环体中都会被重用，而不是重新声明。

各个 goroutine 中输出的 i、v 值都是 for range 循环结束后的 i、v 最终值，而不是各个goroutine启动时的i, v值。可以理解为闭包引用，使用的是上下文环境的值。3 3 3

**使用函数传递**

```
for i, v := range m {
    go func(i,v int) {
        fmt.Println(i, v)
    }(i,v)
}
```

![image-20201019140157702](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20201019140157702.png)

![image-20201019140403567](C:\Users\WIN10\AppData\Roaming\Typora\typora-user-images\image-20201019140403567.png)