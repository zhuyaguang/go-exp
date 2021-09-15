package main

import (
	"fmt"
	"time"

	crdv1beta1 "crd-controller-demo/pkg/apis/stable/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	informers "crd-controller-demo/pkg/client/informers/externalversions/stable/v1"
)

type Controller struct {
	informer  informers.CronTabInformer
	workqueue workqueue.RateLimitingInterface
}

func NewController(informer informers.CronTabInformer) *Controller {
	//使用client 和前面创建的 Informer，初始化了自定义控制器
	controller := &Controller{
		informer: informer,
		// WorkQueue 的实现，负责同步 Informer 和控制循环之间的数据
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "CronTab"),
	}

	klog.Info("Setting up crontab event handlers")

	// informer 注册了三个 Handler（AddFunc、UpdateFunc 和 DeleteFunc）
	// 分别对应 API 对象的“添加”“更新”和“删除”事件。
	// 而具体的处理操作，都是将该事件对应的 API 对象加入到工作队列中
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueCronTab,
		UpdateFunc: func(old, new interface{}) {
			oldObj := old.(*crdv1beta1.CronTab)
			newObj := new.(*crdv1beta1.CronTab)
			// 如果资源版本相同则不处理
			if oldObj.ResourceVersion == newObj.ResourceVersion {
				return
			}
			controller.enqueueCronTab(new)
		},
		DeleteFunc: controller.enqueueCronTabForDelete,
	})
	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// 记录开始日志
	klog.Info("Starting CronTab control loop")
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.informer.Informer().HasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")
	return nil
}

// runWorker 是一个不断运行的方法，并且一直会调用 c.processNextWorkItem 从workqueue读取和读取消息
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// 从workqueue读取和读取消息
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}
	return true
}

// 尝试从 Informer 维护的缓存中拿到了它所对应的 CronTab 对象
func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	crontab, err := c.informer.Lister().CronTabs(namespace).Get(name)

	//从缓存中拿不到这个对象,那就意味着这个 CronTab 对象的 Key 是通过前面的“删除”事件添加进工作队列的。
	if err != nil {
		if errors.IsNotFound(err) {
			// 对应的 crontab 对象已经被删除了
			klog.Warningf("[CronTabCRD] %s/%s does not exist in local cache, will delete it from CronTab ...",
				namespace, name)
			klog.Infof("[CronTabCRD] deleting crontab: %s/%s ...", namespace, name)
			return nil
		}
		runtime.HandleError(fmt.Errorf("failed to get crontab by: %s/%s", namespace, name))
		return err
	}
	klog.Infof("[CronTabCRD] try to process crontab: %#v ...", crontab)
	return nil
}

func (c *Controller) enqueueCronTab(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) enqueueCronTabForDelete(obj interface{}) {
	var key string
	var err error
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
