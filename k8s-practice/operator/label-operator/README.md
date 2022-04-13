### 实现协调逻辑

下面是我们想让 Reconcile 方法做的：

    * 在 ctrl.Request 中使用 Pod 的名称和名称空间从 Kubernetes API 获取 Pod。
    * 如果 Pod 有一个 add-pod-name-label 注释，添加一个 pod-name 标签到 Pod；如果注释缺失，不要添加标签。
    * 在 Kubernetes API 中更新 Pod 以保持所做的更改。