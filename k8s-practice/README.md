## K8S二次开发实践

### 1.通过 Clientset 来获取资源对象(clientset-demo)

- 获取某个命名空间deployment信息
- 实时监控pod的数量
- 通过接口创建deploymnet

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

[kubebuilder 下载地址](https://github.com/kubernetes-sigs/kubebuilder/releases)

### 8.利用operator SDK 开发一个operator(opdemo)

* crd=AppService=deploy+svc



### 9.etcd集群搭建(etcd-cluster-demo)
* 1.静态搭建etcd集群

使用二进制搭建，[启动命令](./etcd-cluster-demo/README.md)

* 2.K8S上搭建etcd集群

etcd 集群的编排的资源清单文件我们可以使用 Kubernetes 源码中提供的，位于目录：test/e2e/testing-manifests/statefulset/etcd 下面

```shell
$ ls -la test/e2e/testing-manifests/statefulset/etcd  
total 40
drwxr-xr-x   6 ych  staff   192 Jun 18  2019 .
drwxr-xr-x  10 ych  staff   320 Oct 10  2018 ..
-rw-r--r--   1 ych  staff   173 Oct 10  2018 pdb.yaml
-rw-r--r--   1 ych  staff   242 Oct 10  2018 service.yaml
-rw-r--r--   1 ych  staff  6441 Jun 18  2019 statefulset.yaml
-rw-r--r--   1 ych  staff   550 Oct 10  2018 tester.yaml
```
* 其中 service.yaml 文件中就是一个用户 StatefulSet 使用的 headless service
* 而 pdb.yaml 文件是用来保证 etcd 的高可用的一个 PodDisruptionBudget 资源对象



### 10.从0到1开发一个etcd operator(etcd-operator-demo)
在开发 Operator 之前我们需要先提前想好我们的 CRD 资源对象
~~~yaml
apiVersion: etcd.ydzs.io/v1alpha1
kind: EtcdCluster
metadata:
  name: demo
spec:
    size: 3  # 副本数量
    image: cnych/etcd:v3.4.13  # 镜像
~~~

* 1.kubebuilder init --domain ydzs.io --owner cnych --repo github.com/cnych/etcd-operator 

* 2.kubebuilder create api --group etcd --version v1alpha1 --kind EtcdCluster

* 3.make 

* 4.修改文件 api/v1alpha1/etcdcluster_types.go 

> 注意每次修改完成后需要执行 make 命令重新生成代码

~~~go
// EtcdClusterSpec defines the desired state of EtcdCluster
type EtcdClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Size  *int32  `json:"size"`
	Image string  `json:"image"`
}
~~~

* 5.业务逻辑




### 11.校验准入控制器实现（admission-registry）

* kind K8S 集群中可以进去kind容器/etc/kubernetes/manifest目录，修改API-server配置

* 只允许使用来自白名单镜像仓库的资源创建 Pod，拒绝使用不受信任的镜像仓库中进行拉取镜像

* cfssl 生成证书

### 12.Mutate 准入控制器实现（admission-registry）

* 当我们的资源对象（Deployment 或 Service）中包含一个需要 mutate 的 annotation 注解后，通过这个 Webhook 后我们就给这个对象添加上一个执行了 mutate 操作的注解

* 管理 Admission Webhook 的 TLS 证书：使用自签名证书，然后通过使用 Init 容器来自行处理 CA。


### [13.实现一个一个自定义的ingress 控制器](https://github.com/cnych/simple-ingress)

* 通过 Kubernetes API 查询和监听 Service、Ingress 以及 Secret 这些对象
* 加载 TLS 证书用于 HTTPS 请求
* 根据加载的 Kubernetes 数据构造一个用于 HTTP 服务的路由，当然该路由需要非常高效，因为所有传入的流量都将通过该路由
* 在 80 和 443 端口上监听传入的 HTTP 请求，然后根据路由查找对应的后端服务，然后代理请求和响应。443 端口将使用 TLS 证书进行安全连接。

### 14.自定义一个调度器(打印日志和GPU）（scheduler-demo）


## 参考资料

[kubebuilder中文官网](https://cloudnative.to/kubebuilder/introduction.html)



   