/*
Copyright 2020 The Knative Authors

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

package v1alpha1

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	apistest "knative.dev/pkg/apis/testing"
)

func TestVSphereSourceDuckTypes(t *testing.T) {
	tests := []struct {
		name string
		t    duck.Implementable
	}{{
		name: "conditions",
		t:    &duckv1.Conditions{},
	}, {
		name: "source",
		t:    &duckv1.Source{},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := duck.VerifyType(&VSphereSource{}, test.t)
			if err != nil {
				t.Errorf("VerifyType(VSphereSource, %T) = %v", test.t, err)
			}
		})
	}
}

func TestVSphereSourceGetGroupVersionKind(t *testing.T) {
	r := &VSphereSource{}
	want := schema.GroupVersionKind{
		Group:   "sources.knative.dev",
		Version: "v1alpha1",
		Kind:    "VSphereSource",
	}
	if got := r.GetGroupVersionKind(); got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestTypicalSourceFlow(t *testing.T) {
	r := &VSphereSourceStatus{}
	r.InitializeConditions()
	apistest.CheckConditionOngoing(r, VSphereSourceConditionReady, t)

	// Check the progression of the SourceReady condition.
	r.PropagateSourceStatus(duckv1.SourceStatus{})
	apistest.CheckConditionOngoing(r, VSphereSourceConditionSourceReady, t)
	r.PropagateSourceStatus(duckv1.SourceStatus{
		Status: duckv1.Status{
			Conditions: []apis.Condition{{
				Type:   apis.ConditionReady,
				Status: corev1.ConditionUnknown,
			}},
		},
	})
	apistest.CheckConditionOngoing(r, VSphereSourceConditionSourceReady, t)
	r.PropagateSourceStatus(duckv1.SourceStatus{
		Status: duckv1.Status{
			Conditions: []apis.Condition{{
				Type:   apis.ConditionReady,
				Status: corev1.ConditionFalse,
			}},
		},
	})
	apistest.CheckConditionFailed(r, VSphereSourceConditionSourceReady, t)
	apistest.CheckConditionFailed(r, VSphereSourceConditionReady, t)
	r.PropagateSourceStatus(duckv1.SourceStatus{
		Status: duckv1.Status{
			Conditions: []apis.Condition{{
				Type:   apis.ConditionReady,
				Status: corev1.ConditionTrue,
			}},
		},
	})
	apistest.CheckConditionSucceeded(r, VSphereSourceConditionSourceReady, t)
	apistest.CheckConditionOngoing(r, VSphereSourceConditionReady, t)

	// Check the progression of the AuthReady condition.
	r.PropagateAuthStatus(duckv1.Status{})
	apistest.CheckConditionOngoing(r, VSphereSourceConditionAuthReady, t)
	r.PropagateAuthStatus(duckv1.Status{
		Conditions: []apis.Condition{{
			Type:   apis.ConditionReady,
			Status: corev1.ConditionUnknown,
		}},
	})
	apistest.CheckConditionOngoing(r, VSphereSourceConditionAuthReady, t)
	r.PropagateAuthStatus(duckv1.Status{
		Conditions: []apis.Condition{{
			Type:   apis.ConditionReady,
			Status: corev1.ConditionFalse,
		}},
	})
	apistest.CheckConditionFailed(r, VSphereSourceConditionAuthReady, t)
	apistest.CheckConditionFailed(r, VSphereSourceConditionReady, t)
	r.PropagateAuthStatus(duckv1.Status{
		Conditions: []apis.Condition{{
			Type:   apis.ConditionReady,
			Status: corev1.ConditionTrue,
		}},
	})
	apistest.CheckConditionSucceeded(r, VSphereSourceConditionAuthReady, t)
	apistest.CheckConditionOngoing(r, VSphereSourceConditionReady, t)

	// Check the progression of the AdapterReady condition.
	r.PropagateAdapterStatus(appsv1.DeploymentStatus{})
	apistest.CheckConditionOngoing(r, VSphereSourceConditionAdapterReady, t)
	r.PropagateAdapterStatus(appsv1.DeploymentStatus{
		Conditions: []appsv1.DeploymentCondition{{
			Type:   appsv1.DeploymentAvailable,
			Status: corev1.ConditionUnknown,
		}},
	})
	apistest.CheckConditionOngoing(r, VSphereSourceConditionAdapterReady, t)
	r.PropagateAdapterStatus(appsv1.DeploymentStatus{
		Conditions: []appsv1.DeploymentCondition{{
			Type:   appsv1.DeploymentAvailable,
			Status: corev1.ConditionFalse,
		}},
	})
	apistest.CheckConditionFailed(r, VSphereSourceConditionAdapterReady, t)
	apistest.CheckConditionFailed(r, VSphereSourceConditionReady, t)
	r.PropagateAdapterStatus(appsv1.DeploymentStatus{
		Conditions: []appsv1.DeploymentCondition{{
			Type:   appsv1.DeploymentAvailable,
			Status: corev1.ConditionTrue,
		}},
	})
	apistest.CheckConditionSucceeded(r, VSphereSourceConditionAdapterReady, t)

	// After all of that, we're finally ready!
	apistest.CheckConditionSucceeded(r, VSphereSourceConditionReady, t)
}
