/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package interpodaffinity

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulerlisters "k8s.io/kubernetes/pkg/scheduler/listers"
)

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "InterPodAffinity"

	defaultHardPodAffinityWeight int32 = 1
	minHardPodAffinityWeight     int32 = 0
	maxHardPodAffinityWeight     int32 = 100
)

// Args holds the args that are used to configure the plugin.
type Args struct {
	HardPodAffinityWeight *int32 `json:"hardPodAffinityWeight,omitempty"`
}

var _ framework.PreFilterPlugin = &InterPodAffinity{}
var _ framework.FilterPlugin = &InterPodAffinity{}
var _ framework.PreScorePlugin = &InterPodAffinity{}
var _ framework.ScorePlugin = &InterPodAffinity{}

// InterPodAffinity is a plugin that checks inter pod affinity
type InterPodAffinity struct {
	sharedLister          schedulerlisters.SharedLister
	hardPodAffinityWeight int32
	sync.Mutex
}

// Name returns name of the plugin. It is used in logs, etc.
func (pl *InterPodAffinity) Name() string {
	return Name
}

// New initializes a new plugin and returns it.
func New(plArgs *runtime.Unknown, h framework.FrameworkHandle) (framework.Plugin, error) {
	if h.SnapshotSharedLister() == nil {
		return nil, fmt.Errorf("SnapshotSharedlister is nil")
	}
	args := &Args{}
	if err := framework.DecodeInto(plArgs, args); err != nil {
		return nil, err
	}
	if err := validateArgs(args); err != nil {
		return nil, err
	}
	pl := &InterPodAffinity{
		sharedLister:          h.SnapshotSharedLister(),
		hardPodAffinityWeight: defaultHardPodAffinityWeight,
	}
	if args.HardPodAffinityWeight != nil {
		pl.hardPodAffinityWeight = *args.HardPodAffinityWeight
	}
	return pl, nil
}

func validateArgs(args *Args) error {
	if args.HardPodAffinityWeight == nil {
		return nil
	}

	weight := *args.HardPodAffinityWeight
	if weight < minHardPodAffinityWeight || weight > maxHardPodAffinityWeight {
		return fmt.Errorf("invalid args.hardPodAffinityWeight: %d, must be in the range %d-%d", weight, minHardPodAffinityWeight, maxHardPodAffinityWeight)
	}
	return nil
}
