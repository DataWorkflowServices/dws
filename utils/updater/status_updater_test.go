/*
 * Copyright 2022 Hewlett Packard Enterprise Development LP
 * Other additional copyright holders may be indicated within.
 *
 * The entirety of this work is licensed under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 *
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
package updater

import (
	"context"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type testObject struct {
	client.Object
	updated bool

	status testStatus
}

func (obj *testObject) GetStatus() Status[*testStatus] {
	return &obj.status
}

type testStatus struct {
	changed bool
	updated bool
}

func (in *testStatus) DeepCopy() *testStatus {
	out := new(testStatus)
	*out = *in
	return out
}

type testWriter struct {
	client.Writer
}

func (*testWriter) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	object, ok := obj.(*testObject)
	if !ok {
		panic("can't convert to test object")
	}

	object.updated = true

	return nil
}

func (*testWriter) Status() client.StatusWriter { return &testStatusWriter{} }

func TestUpdate(t *testing.T)   { testUpdate(t, true) }
func TestNoUpdate(t *testing.T) { testUpdate(t, false) }

// Test that when a change occurs to the object's status, only the object is updated
// and not the status
func testUpdate(t *testing.T, changed bool) {
	obj := &testObject{updated: false}
	updater := NewStatusUpdater[*testStatus](obj)

	obj.status.changed = changed

	updater.CloseWithUpdate(context.TODO(), &testWriter{})

	if obj.updated != changed {
		t.Errorf("Test object not updated")
	}

	if obj.status.updated {
		t.Errorf("Test status incorrectly updated")
	}
}

type testStatusWriter struct {
	client.StatusWriter
}

func (*testStatusWriter) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	object, ok := obj.(*testObject)
	if !ok {
		panic("can't convert to test object")
	}

	object.status.updated = true

	return nil
}

func TestStatusUpdate(t *testing.T)   { testStatusUpdate(t, true) }
func TestNoStatusUpdate(t *testing.T) { testStatusUpdate(t, false) }

// Test when a change occurs to an object's status, only the status fields are
// updated and not the object.
func testStatusUpdate(t *testing.T, changed bool) {
	obj := &testObject{updated: false}
	updater := NewStatusUpdater[*testStatus](obj)

	obj.status.changed = changed // toggle the status changed field so the update occurs

	updater.CloseWithStatusUpdate(context.TODO(), &testWriter{})

	if obj.updated {
		t.Errorf("Test object incorrectly updated")
	}

	if obj.status.updated != changed {
		t.Errorf("Test status not updated")
	}
}
