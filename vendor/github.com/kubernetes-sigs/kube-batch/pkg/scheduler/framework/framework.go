/*
Copyright 2018 The Kubernetes Authors.

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

package framework

import (
	"github.com/golang/glog"

	"github.com/kubernetes-sigs/kube-batch/pkg/scheduler/cache"
	"github.com/kubernetes-sigs/kube-batch/pkg/scheduler/conf"
)

func OpenSession(cache cache.Cache, tiers []conf.Tier) *Session {
	ssn := openSession(cache)
	ssn.Tiers = tiers

	for _, tier := range tiers {
		for _, plugin := range tier.Plugins {
			if pb, found := GetPluginBuilder(plugin.Name); !found {
				glog.Errorf("Failed to get plugin %s.", plugin.Name)
			} else {
				plugin := pb()
				ssn.plugins[plugin.Name()] = plugin
			}
		}
	}

	for _, plugin := range ssn.plugins {
		plugin.OnSessionOpen(ssn)
	}

	return ssn
}

func CloseSession(ssn *Session) {
	for _, plugin := range ssn.plugins {
		plugin.OnSessionClose(ssn)
	}

	closeSession(ssn)
}