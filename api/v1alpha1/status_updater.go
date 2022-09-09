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
// deferred method to close the updater when returning from reconcile.
//
// i.e.
//
//	func (r *myController) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
//		rsrc := &MyResource{}
//		if err := r.Get(ctx, req.NamespacedName, rsrc); err != nil {
//			return err
//		}
//
//		updater := NewStatusUpdater[*MyResourceStatus](rsrc)
//		defer func() {
//			if err == nil {
//				err = updater.Close(ctx, r)
//			}
//		}()
//
//		...
func NewStatusUpdater[S status[S]](rsrc resource[S]) *statusUpdater[S] {
	return &statusUpdater[S]{
		resource: rsrc,
		status:   rsrc.GetStatus().DeepCopy(),
	}
}

// CloseWithUpdate will attempt to update the resource if any of the status fields have
// changed from the initially recorded status. CloseWithUpdate will NOT return an error
// if there is a resource conflict as it's expected the Reconcile() method will be called again.
func (updater *statusUpdater[S]) CloseWithUpdate(ctx context.Context, c client.Writer) error {
	return updater.close(ctx, c)
}

// CloseWithStatusUpdate will attempt to update the resource's status if any of the status
// fields have changed from the initially recorded status. CloseWithStatusUpdate will NOT
// return an error if there is a resource conflict as it's expected the Reconcile() method
// will be called again.
func (updater *statusUpdater[S]) CloseWithStatusUpdate(ctx context.Context, c client.StatusClient) error {
	return updater.close(ctx, c.Status())
}

type clientUpdater interface {
	Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
}

func (updater *statusUpdater[S]) close(ctx context.Context, c clientUpdater) error {
	if !reflect.DeepEqual(updater.resource.GetStatus(), updater.status) {
		err := c.Update(ctx, updater.resource)
		if !errors.IsConflict(err) {
			return err
		}
	}

	return nil
}
