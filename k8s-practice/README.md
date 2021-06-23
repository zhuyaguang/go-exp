## K8S二次开发实践

### 1.通过 Clientset 来获取资源对象(clientset-demo)

- 获取某个命名空间deployment信息
- 实时监控pod的数量
- 通过接口创建deploymnet

* [使用k8s.io/client-go的dynamic client的示例](https://zhuanlan.zhihu.com/p/165970638)

### 2.通过informer来监控资源对象(informer-demo)

> 为什么要用informer来监控资源对象
>
> 前面我们在使用 Clientset 的时候了解到我们可以使用 Clientset 来获取所有的原生资源对象，那么如果我们想要去一直获取集群的资源对象数据呢？岂不是需要用一个轮询去不断执行 List() 操作？这显然是不合理的，实际上除了常用的 CRUD 操作之外，我们还可以进行 Watch 操作，可以监听资源对象的增、删、改、查操作，这样我们就可以根据自己的业务逻辑去处理这些数据了。

* 监控某个命名空间内的deployment资源的创建、更新、删除

### 3.自定义一个Indexer, 存储pod资源(indexer-demo)

* 根据namespace和nodename过滤pod信息


### 4.实现一个pod controller(pod-controller-demo)

* 监控某个namespace下面的pod

### 5.编写一个标准的crd资源文件(crd-demo)



### 6.给自定义的crd资源编写一个控制器(crd-controller-demo)



### 7.利用kubebuilder写一个operator(kubebuilder-demo)



### 8.利用operator SDK 开发一个operator(opdemo)



### 9.静态搭建etcd集群(etcd-cluster-demo)



### 10.K8S上搭建etcd集群(etcd-cluster-demo)



11.从0到1开发一个etcd operator

12.校验准入控制器实现

13.Mutate 准入控制器实现

14.自动生成证书和自动注入

15.实现一个一个自定义的ingress 控制器

16.自定义一个调度器(打印日志和GPU）