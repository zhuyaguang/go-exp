启动3个etcd实例



```shell
/tmp/etcd/etcd --name s1 \  # etcd 节点名称  
--data-dir /tmp/etcd/s1 \  # 数据存储目录   
--listen-client-urls http://localhost:2379 \  # 本节点访问地址  
--advertise-client-urls http://localhost:2379 \  # 用于通知其他 ETCD 节点，客户端接入本节点的监听地址  
--listen-peer-urls http://localhost:2380 \  # 本节点与其他节点进行数据交换的监听地址  
--initial-advertise-peer-urls http://localhost:2380 \  # 通知其他节点与本节点进行数据交换的地址  
--initial-cluster s1=http://localhost:2380,s2=http://localhost:22380,s3=http://localhost:32380 \  
# 集群所有节点配置  
--initial-cluster-token tkn \  # 集群唯一标识  
--initial-cluster-state new  # 节点初始化方式

#启动第一个节点
/tmp/etcd/etcd --name s1 --data-dir /tmp/etcd/s1 --listen-client-urls http://localhost:2379 --advertise-client-urls http://localhost:2379 --listen-peer-urls http://localhost:2380 --initial-advertise-peer-urls http://localhost:2380 --initial-cluster s1=http://localhost:2380,s2=http://localhost:22380,s3=http://localhost:32380 --initial-cluster-token tkn --initial-cluster-state new 

#启动第二个节点
/tmp/etcd/etcd --name s2 \
  --data-dir /tmp/etcd/s2 \
  --listen-client-urls http://localhost:22379 \
  --advertise-client-urls http://localhost:22379 \
  --listen-peer-urls http://localhost:22380 \
  --initial-advertise-peer-urls http://localhost:22380 \
  --initial-cluster s1=http://localhost:2380,s2=http://localhost:22380,s3=http://localhost:32380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new
  
  
  #启动第三个节点
  /tmp/etcd/etcd --name s3 \
  --data-dir /tmp/etcd/s3 \
  --listen-client-urls http://localhost:32379 \
  --advertise-client-urls http://localhost:32379 \
  --listen-peer-urls http://localhost:32380 \
  --initial-advertise-peer-urls http://localhost:32380 \
  --initial-cluster s1=http://localhost:2380,s2=http://localhost:22380,s3=http://localhost:32380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new
  
  
  #查看集群状态
ETCDCTL_API=3 /tmp/etcd/etcdctl \
  --endpoints localhost:2379,localhost:22379,localhost:32379 \
  endpoint health
  
ETCDCTL_API=3 /tmp/etcd/etcdctl --endpoints localhost:2379,localhost:22379,localhost:32379 endpoint status --write-out=table
  
```

