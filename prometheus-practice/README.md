

## 一、k8s 部署 prometheus
> 提前创建好namespaces：`kubectl create  ns kube-mon`

### 1. 用 ConfigMap 形式配置prometheus.yaml(prometheus-cm.yaml)

`kubectl apply -f prometheus-cm.yaml`



prometheus 的 ConfigMap 更新完成后，执行 reload 操作，让配置生效：

~~~shell
$ kubectl apply -f prometheus-cm.yaml
configmap/prometheus-config configured
# 隔一会儿执行reload操作
$ curl -X POST "http://10.244.3.174:9090/-/reload"
~~~



### 2. 创建 prometheus 的 Pod 资源：(prometheus-deploy.yaml)

### 3. 需要配置 rbac 相关认证: (prometheus-rbac.yaml)

### 4. 创建一个 Service 对象
Pod 创建成功后，为了能够在外部访问到 prometheus 的 webui 服务，我们还需要创建一个 Service 对象：(prometheus-svc.yaml)

### 5. 使用 exporter 监控 
   部署一个 redis 应用，并用 redis-exporter 的方式来采集监控数据供 Prometheus 使用，如下资源清单文件：（prome-redis.yaml）

### 6. 监控Kubernetes 集群本身

* Kubernetes 节点的监控：比如节点的 cpu、load、disk、memory 等指标
* 内部系统组件的状态：比如 kube-scheduler、kube-controller-manager、kubedns/coredns 等组件的详细运行状态
* 编排级的 metrics：比如 Deployment 的状态、资源请求、调度和 API 延迟等数据指标

### 7. 监控集群节点 (prome-node-exporter.yaml)

也可以用 helm安装：$ helm upgrade --install node-exporter --namespace kube-mon stable/prometheus-node-exporter

### 8.服务发现

在 Kubernetes 下，Promethues 通过与 Kubernetes API 集成，主要支持5中服务发现模式，分别是：Node、Service、Pod、Endpoints、Ingress。

~~~yaml
- job_name: 'kubernetes-nodes'
  kubernetes_sd_configs:
    - role: node
~~~

### 9. 容器监控

cAdvisor已经内置在了 kubelet 组件之中，所以我们不需要单独去安装，`cAdvisor` 的数据路径为 `/api/v1/nodes/<node>/proxy/metrics`

### 10. 监控 apiserver

~~~yaml
- job_name: 'kubernetes-apiservers'
  kubernetes_sd_configs:
  - role: endpoints
~~~

### 11. 监控pod

apiserver 实际上就是一种特殊的 Endpoints，现在我们同样来配置一个任务用来专门发现普通类型的 Endpoint，其实就是 Service 关联的 Pod 列表：

~~~yaml
- job_name: 'kubernetes-endpoints'
  kubernetes_sd_configs:
  - role: endpoints
  relabel_configs:
  - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape]
    action: keep
    regex: true
  - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
    action: replace
    target_label: __scheme__
    regex: (https?)
  - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
    action: replace
    target_label: __metrics_path__
    regex: (.+)
  - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
    action: replace
    target_label: __address__
    regex: ([^:]+)(?::\d+)?;(\d+)
    replacement: $1:$2
  - action: labelmap
    regex: __meta_kubernetes_service_label_(.+)
  - source_labels: [__meta_kubernetes_namespace]
    action: replace
    target_label: kubernetes_namespace
  - source_labels: [__meta_kubernetes_service_name]
    action: replace
    target_label: kubernetes_name
  - source_labels: [__meta_kubernetes_pod_name]
    action: replace
    target_label: kubernetes_pod_name
~~~


## 二、k8s 部署grafana  (grafana.yaml)

### 安装插件 DevOpsProdigy KubeGraf

> 它是 Grafana 官方的 Kubernetes 插件 的升级版本，该插件可以用来可视化和分析 Kubernetes 集群的性能，通过各种图形直观的展示了 Kubernetes 集群的主要服务的指标和特征，还可以用于检查应用程序的生命周期和错误日志。

~~~
$ kubectl exec -it grafana-5579769f64-7729f -n kube-mon /bin/bash
bash-5.0# grafana-cli plugins install devopsprodigy-kubegraf-app

installing devopsprodigy-kubegraf-app @ 1.3.0
from: https://grafana.com/api/plugins/devopsprodigy-kubegraf-app/versions/1.3.0/download
into: /var/lib/grafana/plugins

✔ Installed devopsprodigy-kubegraf-app successfully 

Restart grafana after installing plugins . <service grafana-server restart>

# 由于该插件依赖另外一个 Grafana-piechart-panel 插件，所以如果没有安装，同样需要先安装该插件。
bash-5.0# grafana-cli plugins install Grafana-piechart-panel
......
~~~

安装删除grafana pod 重启

### 设置集群监控插件

* 点击 Set up your first k8s-cluster 创建一个新的 Kubernetes 集群:

* URL 使用 Kubernetes Service 地址即可：https://kubernetes.default:443
* Access 访问模式使用：Server(default)
* 由于插件访问 Kubernetes 集群的各种资源对象信息，所以我们需要配置访问权限，这里我们可以简单使用 kubectl 的 kubeconfig 来进行配置即可。
勾选 Auth 下面的 TLS Client Auth 和 With CA Cert 两个选项
* 其中 TLS Auth Details 下面的值就对应 kubeconfig 里面的证书信息。比如我们这里的 kubeconfig 文件格式如下所示：

~~~
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: <certificate-authority-data>
    server: https://ydzs-master:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: 'kubernetes-admin@kubernetes'
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: <client-certificate-data>
    client-key-data: <client-key-data>
~~~

## 三、k8s 部署 Alertmanager

### 创建指定配置文件 (alertmanager-config.yaml)

`kubectl apply -f alertmanager-config.yaml`

### 创建 alertmanager deployment

`kubectl apply -f alertmanager-deploy.yaml`

### 创建service alertmanager-svc.yaml
 更新普罗米修斯配置文件
~~~yaml
alerting:
  alertmanagers:
    - static_configs:
      - targets: ["alertmanager:9093"]
~~~

### 添加报警规则

在Prometheus 的配置文件中添加如下报警规则配置：

~~~yaml
rule_files:
- /etc/prometheus/rules.yml
~~~

这里我们同样将 rules.yml 文件用 ConfigMap 的形式挂载到 /etc/prometheus 目录下面即可，比如下面的规则：（alert-rules.yml）

~~~yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
  namespace: kube-mon
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
      scrape_timeout: 15s
      evaluation_interval: 30s  # 默认情况下每分钟对告警规则进行计算
    alerting:
      alertmanagers:
      - static_configs:
        - targets: ["alertmanager:9093"]
    rule_files:
    - /etc/prometheus/rules.yml
  ...... # 省略prometheus其他部分
  rules.yml: |
    groups:
    - name: test-node-mem
      rules:
      - alert: NodeMemoryUsage
        expr: (node_memory_MemTotal_bytes - (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)) / node_memory_MemTotal_bytes * 100 > 20
        for: 2m
        labels:
          team: node
        annotations:
          summary: "{{$labels.instance}}: High Memory usage detected"
          description: "{{$labels.instance}}: Memory usage is above 20% (current value is: {{ $value }}"
~~~

 * 开箱即用的 Prometheus 告警规则集 :[Awesome Prometheus Alerts](https://github.com/samber/awesome-prometheus-alerts)

###  设置 WebHook 接收器
* 告警插件 推送钉钉： https://github.com/timonwong/prometheus-webhook-dingtalk 

## 四、Prometheus 高可用

## 五、自定义指标扩缩容

##  六、PromQL

## 七、Prometheus Operator

## 八、容器化部署 系统监控容器

1. node-exporter监控节点  gpu-exporter 监控GPU  cadvisor监控容器

~~~shell
docker run -d  --restart always   -v "/proc:/host/proc:ro"   -v "/sys:/host/sys:ro"   -v "/:/rootfs:ro"   --net="host"   --name node-exporter   prom/node-exporter
docker run -d --gpus all --restart always  --name gpu-exporter -p 9400:9400 nvidia/dcgm-exporter:2.0.13-2.1.1-ubuntu18.04
docker run  --restart always   --volume=/:/rootfs:ro   --volume=/var/run:/var/run:ro   --volume=/sys:/sys:ro   --volume=/var/lib/docker/:/var/lib/docker:ro  --volume=/dev/disk/:/dev/disk:ro  --publish=8080:8080   --detach=true   --name=cadvisor   google/cadvisor:latest
~~~

2. 普罗米修斯安装
~~~shell
docker run  -d   -p 9090:9090   -v /data/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml  --name prometheus   prom/prometheus
~~~

3. grafana安装
~~~
docker run -d   -p 3000:3000   --name=grafana -v /data/grafana-storage:/var/lib/grafana grafana/grafana
~~~

## [Prometheus 监控外部 Kubernetes 集群](https://www.qikqiak.com/post/monitor-external-k8s-on-prometheus/)

* 创建用于 Prometheus 访问 Kubernetes 资源对象的 RBAC 对象