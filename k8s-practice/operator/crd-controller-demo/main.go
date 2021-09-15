package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog"

	clientset "crd-controller-demo/pkg/client/clientset/versioned"
	informers "crd-controller-demo/pkg/client/informers/externalversions"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

// SetupSignalHandler 注册 SIGTERM 和 SIGINT 信号
// 返回一个 stop channel，该通道在捕获到第一个信号时被关闭
// 如果捕捉到第二个信号，程序将直接退出
func setupSignalHandler() (stopCh <-chan struct{}) {
	// 当调用两次的时候 panics
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	// Notify 函数让 signal 包将输入信号转发到 c
	// 如果没有列出要传递的信号，会将所有输入信号传递到 c；否则只传递列出的输入信号
	// 参考文档：https://cloud.tencent.com/developer/article/1645996
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // 第二个信号，直接退出
	}()

	return stop
}

func initClient() (*kubernetes.Clientset, *rest.Config, error) {
	var err error
	var config *rest.Config
	// inCluster（Pod）、KubeConfig（kubectl）
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "kind-config-my-cluster"), "(可选) kubeconfig 文件的绝对路径")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig 文件的绝对路径")
	}
	flag.Parse()

	// 首先使用 inCluster 模式(需要去配置对应的 RBAC 权限，默认的sa是default->是没有获取deployments的List权限)
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 KubeConfig 文件创建集群配置 Config 对象
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	// 已经获得了 rest.Config 对象
	// 创建 Clientset 对象
	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, config, err
	}
	return kubeclient, config, nil
}

func main() {
	flag.Parse()

	//设置一个信号处理，应用于优雅关闭
	stopCh := setupSignalHandler()

	_, cfg, err := initClient()
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	// 实例化一个 CronTab 的 ClientSet
	crontabClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building crontab clientset: %s", err.Error())
	}

	// informerFactory 工厂类， 这里注入我们通过代码生成的 client
	// clent 主要用于和 API Server 进行通信，实现 ListAndWatch
	crontabInformerFactory := informers.NewSharedInformerFactory(crontabClient, time.Second*30)

	// 实例化自定义控制器
	controller := NewController(crontabInformerFactory.Crd().V1().CronTabs())

	// 启动 informer，开始List & Watch
	go crontabInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
