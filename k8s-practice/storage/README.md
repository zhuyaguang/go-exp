### 环境准备

*  安装 lvm2 软件包
~~~
# Centos
sudo yum install -y lvm2

# Ubuntu
sudo apt-get install -y lvm2
~~~

### 部署rook-ceph-operator

> https://github.com/rook/rook/tree/release-1.2/cluster/examples/kubernetes/ceph  下载 common.yaml 与 operator.yaml 

~~~
# 会安装crd、rbac相关资源对象
$ kubectl apply -f common.yaml
# 安装 rook operator
$ kubectl apply -f operator.yaml
~~~

### 创建 Ceph 集群

* 确保设置的 dataDirHostPath 属性值为有效得主机路径

`kubectl apply -f cluster.yaml `


### 部署 Rook 工具箱 来验证集群

`kubectl apply -f toolbox.yaml`
`kubectl exec rook-ceph-tools-f54bf64c6-dlc4n -i -t  bash  -n rook-ceph`

### 部署Ceph  Dashboard 工具

`kubectl apply -f dashboard-external.yaml`
用户名 admin
获取密码:
`kubectl -n rook-ceph get secret rook-ceph-dashboard-password -o jsonpath="{['data']['password']}" | base64 --decode && echo xxxx（登录密码）`
### 使用Ceph 集群
* 创建 RBD 类型的存储池 

`kubectl apply -f pool.yaml `

* 创建 StorageClass

`kubectl apply -f storageclass.yaml `

* 创建一个 PVC 来使用上面的 StorageClass 对象

`kubectl apply -f pvc.yaml `

