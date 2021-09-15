## [容器化部署EFK&ELK](https://www.jianshu.com/p/cabb33fba1b4)

## K8S部署EFK&ELK

### [yaml部署](https://www.qikqiak.com/k8strain/logging/efk/)
1. 创建ns
`kubectl create -f kube-logging.yaml`
2. 创建服务elasticsearch-svc.yaml
`kubectl create -f elasticsearch-svc.yaml`
3. 创建elasticsearch-statefulset.yaml 
` kubectl create -f elasticsearch-statefulset.yaml`

* 对于 volumeClaimTemplates 手动创建一个PV

* 或者用loacl path

`kubectl apply -f pv-hostpath.yaml`

`kubectl create -f pvc-hostpath.yaml`
4. 创建Kibana服务 
`kubectl create -f kibana.yaml`
5. 部署Fluentd

* 新建 fluentd-configmap.yaml

* 新建一个 fluentd-daemonset.yaml 的文件

6. 基于日志的报警

部署 elastalert.yaml

###  [Elastic Cloud on Kubernetes(ECK) 部署](https://www.qikqiak.com/post/elastic-cloud-on-k8s/)

Elastic Cloud on Kubernetes(ECK)是一个 Elasticsearch Operator，但远不止于此。 ECK 使用 Kubernetes Operator 模式构建而成，需要安装在您的 Kubernetes 集群内，其功能绝不仅限于简化 Kubernetes 上 Elasticsearch 和 Kibana 的部署工作这一项任务。ECK 专注于简化所有后期运行工作，例如：

    管理和监测多个集群
    轻松升级至新的版本
    扩大或缩小集群容量
    更改集群配置
    动态调整本地存储的规模（包括 Elastic Local Volume（一款本地存储驱动器））
    备份

ECK 不仅能自动完成所有运行和集群管理任务，还专注于简化在 Kubernetes 上使用 Elasticsearch 的完整体验。ECK 的愿景是为 Kubernetes 上的 Elastic 产品和解决方案提供 SaaS 般的体验。 

* 安装ECK
`kubectl apply -f https://download.elastic.co/downloads/eck/0.8.1/all-in-one.yaml`

* 利用 CRD 对象来创建一个非常简单的单个 Elasticsearch 集群：(elastic.yaml)

`kubectl create -f elastic.yaml`

>  ECK 添加了一个 validation webhook 的 Admission，我们可以临时将这个对象删除

* 用 CRD 对象 Kibana 来部署 kibana 应用：(kibana.yaml)

`kubectl create -f kibana.yaml`


* 访问 kibana 来验证我们的集群，比如我们可以再添加一个 Ingress 对象：(ingress.yaml)

`kubectl create -f ingress.yaml`

## [容器化部署PLG](https://blog.csdn.net/qq_30442207/article/details/114583870)

## [K8S部署Loki](https://www.qikqiak.com/post/grafana-loki-usage/)


##  [使用 Elastic 技术栈构建 K8S 全栈监控](https://www.qikqiak.com/post/k8s-monitor-use-elastic-stack-1/)

* 监控指标提供系统各个组件的时间序列数据，比如 CPU、内存、磁盘、网络等信息，通常可以用来显示系统的整体状况以及检测某个时间的异常行为
* 日志为运维人员提供了一个数据来分析系统的一些错误行为，通常将系统、服务和应用的日志集中收集在同一个数据库中
* 追踪或者 APM（应用性能监控）提供了一个更加详细的应用视图，可以将服务执行的每一个请求和步骤都记录下来（比如 HTTP 调用、数据库查询等），通过追踪这些数据，我们可以检测到服务的性能，并相应地改进或修复我们的系统。
