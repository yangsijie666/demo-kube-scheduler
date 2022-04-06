package scheduler

import (
	"context"
	"errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

const Name = "demo-scheduler"

var scheme = runtime.NewScheme()

type DemoScheduler struct {
	handle v1alpha1.FrameworkHandle
	cache  cache.Cache
}

func (d *DemoScheduler) Name() string {
	return Name
}

func New(_ runtime.Object, h v1alpha1.FrameworkHandle) (v1alpha1.Plugin, error) {
	mgrConfig := ctrl.GetConfigOrDie()
	mgrConfig.QPS = 1000
	mgrConfig.Burst = 1000

	mgr, err := ctrl.NewManager(mgrConfig, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "",
		LeaderElection:     false,
		Port:               9443,
	})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	go func() {
		if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			klog.Error(err)
			panic(err)
		}
	}()

	scvCache := mgr.GetCache()
	if scvCache.WaitForCacheSync(context.TODO()) {
		return &DemoScheduler{
			handle: h,
			cache:  scvCache,
		}, nil
	} else {
		return nil, errors.New("Cache Not Sync! ")
	}
}

func (y *DemoScheduler) Reserve(ctx context.Context, state *v1alpha1.CycleState, pod *v1.Pod, nodeName string) *v1alpha1.Status {
	klog.Infof("[Reserve] Node: %v, in Plugin: demo-scheduler When scheduling Pod: %v/%v", nodeName, pod.GetNamespace(), pod.GetName())
	return v1alpha1.NewStatus(v1alpha1.Success, "")
}

func (y *DemoScheduler) Unreserve(ctx context.Context, _ *v1alpha1.CycleState, pod *v1.Pod, nodeName string) {
	klog.Infof("[Unreserve] Node: %v, in Plugin: demo-scheduler When scheduling Pod: %v/%v", nodeName, pod.GetNamespace(), pod.GetName())
}
