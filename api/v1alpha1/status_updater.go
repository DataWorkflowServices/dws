package v1alpha1

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type resource[T any] interface {
	client.Object
	GetStatus() status[T]
}

type status[T any] interface {
	DeepCopy() T
}

type statusUpdater[T any] struct {
	resource resource[T]
	status   T
}

// NewStatusUpdater returns a status updater meant for updating the status of the
// supplied resource when the updater is closed. Typically users will want to
// create a status updater early on in a controller's Reconcile() method and add a
// deferred method to close the updater when returning reconcile.
//
// i.e.
//
//	func (r *myController) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
//		rsrc := &MyResource{}
//		if err := r.Get(ctx, req.NamespacedName, rsrc); err != nil {
//			return err
//		}
//
//		updater := NewStatusUpdater[*MyResource, *MyResourceStatus](rsrc)
//		defer func() {
//			if err == nil {
//				err = updater.Close(ctx, r)
//			}
//		}()
//
//		...
func NewStatusUpdater[T resource[S], S status[S]](rsrc T) *statusUpdater[S] {
	return &statusUpdater[S]{
		resource: rsrc,
		status:   rsrc.GetStatus().DeepCopy(),
	}
}

// Close will attempt to update the status of the updating resource if it has changed
// from the initially recorded status. Close will NOT return an error if there is a
// conflict as it's expected the Reconcile() method will be called again.
func (updater *statusUpdater[S]) Close(ctx context.Context, c client.Client) error {
	if !reflect.DeepEqual(updater.resource.GetStatus(), updater.status) {
		err := c.Status().Update(ctx, updater.resource)
		if !errors.IsConflict(err) {
			return err
		}
	}

	return nil
}
