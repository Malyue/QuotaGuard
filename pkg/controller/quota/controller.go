package quota

import (
	clientset "github.com/Malyue/quotaguard/pkg/generated/clientset/versioned"
	informers "github.com/Malyue/quotaguard/pkg/generated/informers/externalversions"
	listers "github.com/Malyue/quotaguard/pkg/generated/listers/quota/v1"
	"github.com/Malyue/quotaguard/pkg/quota"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sync"
)

type Controller struct {
	crdClient    clientset.Interface
	lister       listers.QuotaPolicyLister
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	quotaManager *quota.QuotaManager
	workers      int
	stopCh       chan struct{}
	wg           sync.WaitGroup
}

func NewController(
	crdClient clientset.Interface,
	informerFactory informers.SharedInformerFactory,
	quotaManager *quota.QuotaManager,
	workers int) *Controller {

	quotaInformer := informerFactory.Quota().V1().QuotaPolicies()

	c := &Controller{
		crdClient:    crdClient,
		lister:       quotaInformer.Lister(),
		queue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "quotapolicies"),
		informer:     quotaInformer.Informer(),
		quotaManager: quotaManager,
		workers:      workers,
		stopCh:       make(chan struct{}),
	}

	quotaInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	})

	return c
}

func (c *Controller) Start() {

	go c.informer.Run(c.stopCh)

	if !cache.WaitForCacheSync(c.stopCh, c.informer.HasSynced) {
		panic("failed to wait for caches to sync")
	}

	for i := 0; i < c.workers; i++ {
		c.wg.Add(1)
		go c.runWorker()
	}

	klog.Info("QuotaPolicy controller started")

}

func (c *Controller) Stop() {
	close(c.stopCh)    // 1. 先停止 Informer
	c.queue.ShutDown() // 2. 关闭工作队列
	c.wg.Wait()        // 3. 等待所有 worker 退出
	klog.Info("QuotaPolicy controller stopped")
}

func (c *Controller) runWorker() {
	defer c.wg.Done()
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncHandler(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	c.handleError(err, key)
	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	policy, err := c.lister.QuotaPolicies(namespace).Get(name)
	switch {
	case errors.IsNotFound(err):
		err = c.quotaManager.DeleteQuotaPolicy(policy)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		err = c.quotaManager.UpdateQuotaPolicy(policy)
		if err != nil {
			return err
		}
	}

	return nil

}
