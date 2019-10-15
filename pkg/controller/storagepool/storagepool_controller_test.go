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

package storagepool

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

func ReconcileWithRequeue(r *ReconcileStoragePool , request reconcile.Request) (reconcile.Result, error) {
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
func TestStoragePoolController1(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(dwsv1.SchemeGroupVersion, &dwsv1.StoragePool{})

	// Create a fake client to mock API calls.
	objs := []runtime.Object{}
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileMemcached object with the scheme and fake client.
	r := &ReconcileStoragePool{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "storagepool-test1",
			Namespace: "dws",
		},
	}
	_, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
}

/*
 * Run the reconciler with a new StoragePool  object in the fake client. Check
 * that the StoragePool  object was created Run the reconciler again and verify
 * no errors.
 */
func TestStoragePoolController2(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name      = "storagepool-test2"
		namespace = "dws"
	)

	// A Memcached resource with metadata and spec.
	storagepool := &dwsv1.StoragePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: dwsv1.StoragePoolSpec{
			PoolId:           "storagepool-test",
			Units:            "bytes",
			Granularity:      107374182400,
			Quantity:         10,
			Free:             10,
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		storagepool,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(dwsv1.SchemeGroupVersion, &dwsv1.StoragePool{})

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileMemcached object with the scheme and fake client.
	r := &ReconcileStoragePool{client: cl, scheme: s}

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

	storagepooltst := &dwsv1.StoragePool{}
	err = cl.Get(context.TODO(), namespacedname, storagepooltst)
	if err != nil {
		t.Error("StoragePool object not created")
	}
}
