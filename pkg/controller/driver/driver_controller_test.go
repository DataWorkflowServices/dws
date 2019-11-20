/*
 * Copyright 2019 Cray Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package driver

import (
	"context"
	"testing"

	dwsv1 "stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

/*
 * Run the reconciler with requeueing. This mimics how the controller would
 * actually be called.
 */

func ReconcileWithRequeue(r *ReconcileDriver , request reconcile.Request) (reconcile.Result, error) {
	for {
		res, err := r.Reconcile(request)
		if err != nil {
			// Do not requeue on error
			return res, err
		}
		if !res.Requeue {
			return res, err
		}
	}
}

/*
 * Run the reconciler without any objects in the fake client. Reconcile()
 * should not return an error in this situation.
 */
func TestDriverController1(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(dwsv1.SchemeGroupVersion, &dwsv1.Driver{})

	// Create a fake client to mock API calls.
	objs := []runtime.Object{}
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileDriver object with the scheme and fake client.
	r := &ReconcileDriver{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "driver-test1",
			Namespace: "dws",
		},
	}
	_, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
}

/*
 * Run the reconciler with a new Driver  object in the fake client. Check
 * that the Driver  object was created Run the reconciler again and verify
 * no errors.
 */
func TestDriverController2(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name      = "driver-test2"
		namespace = "dws"
	)

	var statesToWatch = []string {
		"prerun",
		"postrun",
	}

	// A driver resource with metadata and spec.
	driver := &dwsv1.Driver{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: dwsv1.DriverSpec{
			DriverID:         "5f239bd8-30db-450b-8c2c-a1a7c8631a1a",
			WatchStates:      statesToWatch,
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		driver,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(dwsv1.SchemeGroupVersion, &dwsv1.Driver{})

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileDriver object with the scheme and fake client.
	r := &ReconcileDriver{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	_, err := ReconcileWithRequeue(r, req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	namespacedname := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	drivertst := &dwsv1.Driver{}
	err = cl.Get(context.TODO(), namespacedname, drivertst)
	if err != nil {
		t.Error("Driver object not created")
	}
}
