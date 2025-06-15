package quota

import (
	"fmt"
	v1 "github.com/Malyue/quotaguard/pkg/apis/quota/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"reflect"
	"time"
)

func (c *Controller) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.queue.Add(key)
}

func (c *Controller) onUpdate(oldObj, newObj interface{}) {
	oldPolicy := oldObj.(*v1.QuotaPolicy)
	newPolicy := newObj.(*v1.QuotaPolicy)

	// only spec change
	if reflect.DeepEqual(oldPolicy.Spec, newPolicy.Spec) {
		return
	}

	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.queue.AddAfter(key, 500*time.Microsecond)
}

func (c *Controller) onDelete(obj interface{}) {
	_, ok := obj.(*v1.QuotaPolicy)
	if !ok {
		// 处理 etcd 中的墓碑状态
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("unexpected object type: %T", obj))
			return
		}

		_, ok = tombstone.Obj.(*v1.QuotaPolicy)
		if !ok {
			runtime.HandleError(fmt.Errorf("tombstone contained unexpected object type: %T", tombstone.Obj))
			return
		}
	}

	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}

	c.queue.Add(key)
}
