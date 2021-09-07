package plugins

import (
"context"
"fmt"

v1 "k8s.io/api/core/v1"
"k8s.io/apimachinery/pkg/runtime"
"k8s.io/klog/v2"
frameworkruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

const (
// Name 定义插件名称
Name = "sample-plugin"
preFilterStateKey ="PreFilter"+ Name
)

var _ framework.PreFilterPlugin = &Sample{}
var _ framework.FilterPlugin = &Sample{}

type SampleArgs struct {
FavoriteColor string `json:"favorColor,omitempty"`
FavoriteNumber int `json:"favorNumber,omitempty"`
ThanksTo string `json:"thanksTo,omitempty"`
}

// 获取插件配置的参数
func getSampleArgs(object runtime.Object) (*SampleArgs, error) {
sa := &SampleArgs{}
if err := frameworkruntime.DecodeInto(object, sa); err != nil {
return nil, err
}
return sa, nil
}

type preFilterState struct {
framework.Resource   // requests,limits
}

func (s *preFilterState) Clone() framework.StateData {
return s
}

func getPreFilterState(state *framework.CycleState) (*preFilterState, error) {
data, err := state.Read(preFilterStateKey)
if err != nil {
return nil, err
}
s, ok := data.(*preFilterState)
if !ok {
return nil, fmt.Errorf("%+v convert to SamplePlugin preFilterState error", data)
}
return s, nil
}

type Sample struct {
args *SampleArgs
handle framework.FrameworkHandle
}

func (s *Sample) Name() string {
return Name
}

func computePodResourceLimit(pod *v1.Pod) *preFilterState {
result := &preFilterState{}
for _, container := range pod.Spec.Containers {
result.Add(container.Resources.Limits)
}
return result
}

func (s *Sample) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) *framework.Status {
if klog.V(2).Enabled() {
klog.InfoS("Start PreFilter Pod", "pod", pod.Name)
}
state.Write(preFilterStateKey, computePodResourceLimit(pod))
return nil
}

func (s *Sample) PreFilterExtensions() framework.PreFilterExtensions {
return nil
}

func (s *Sample) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
preState, err := getPreFilterState(state)
if err != nil {
return framework.NewStatus(framework.Error, err.Error())
}
if klog.V(2).Enabled() {
klog.InfoS("Start Filter Pod", "pod", pod.Name, "node", nodeInfo.Node().Name, "preFilterState", preState)
}
// logic
return framework.NewStatus(framework.Success, "")
}

//type PluginFactory = func(configuration runtime.Object, f v1alpha1.FrameworkHandle) (v1alpha1.Plugin, error)

func New(object runtime.Object, f framework.FrameworkHandle) (framework.Plugin, error) {
args, err := getSampleArgs(object)
if err != nil {
return nil, err
}
// validate args
if klog.V(2).Enabled() {
klog.InfoS("Successfully get plugin config args", "plugin", Name, "args", args)
}
return &Sample{
args: args,
handle: f,
}, nil
}