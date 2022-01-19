/*
Copyright 2022 Hewlett Packard Enterprise Development LP
*/

package controllers

import (
	"context"
	"os/exec"
	"reflect"
	"strings"

	"github.com/go-logr/logr"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	dwsv1alpha1 "github.hpe.com/hpe/hpc-dpm-dws-operator/api/v1alpha1"
)

// ClientMountReconciler reconciles a ClientMount object
type ClientMountReconciler struct {
	client.Client
	Mock   bool
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const (
	// finalizerClientMount defines the key used for the finalizer
	finalizerClientMount = "nnf.cray.hpe.com/client_mount"
)

//+kubebuilder:rbac:groups=dws.cray.hpe.com,resources=clientmounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dws.cray.hpe.com,resources=clientmounts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dws.cray.hpe.com,resources=clientmounts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ClientMountReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
	log := r.Log.WithValues("ClientMount", req.NamespacedName)
	clientMount := &dwsv1alpha1.ClientMount{}
	if err := r.Get(ctx, req.NamespacedName, clientMount); err != nil {
		// ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Create a status updater that handles the call to status().Update() if any of the fields
	// in clientMount.Status change
	statusUpdater := newClientMountStatusUpdater(clientMount)
	defer func() {
		if err == nil {
			err = statusUpdater.close(ctx, r)
		}
	}()

	// Handle cleanup if the resource is being deleted
	if !clientMount.GetDeletionTimestamp().IsZero() {
		if !controllerutil.ContainsFinalizer(clientMount, finalizerClientMount) {
			return ctrl.Result{}, nil
		}

		// Unmount everything before removing the finalizer
		log.Info("Unmounting all file systems due to resource deletion")
		if err := r.unmountAll(ctx, clientMount); err != nil {
			return ctrl.Result{}, err
		}

		controllerutil.RemoveFinalizer(clientMount, finalizerClientMount)
		if err := r.Update(ctx, clientMount); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// Create the status section if it doesn't exist yet
	if len(clientMount.Status.Mounts) != len(clientMount.Spec.Mounts) {
		clientMount.Status.Mounts = make([]dwsv1alpha1.ClientMountInfoStatus, len(clientMount.Spec.Mounts))
	}

	// Initialize the status section if the desired state doesn't match the status state
	if clientMount.Status.Mounts[0].State != clientMount.Spec.DesiredState {
		for i := 0; i < len(clientMount.Status.Mounts); i++ {
			clientMount.Status.Mounts[i].State = clientMount.Spec.DesiredState
			clientMount.Status.Mounts[i].Ready = false
			clientMount.Status.Mounts[i].Message = ""
		}

		return ctrl.Result{}, nil
	}

	// Add finalizer if it doesn't exist
	if !controllerutil.ContainsFinalizer(clientMount, finalizerClientMount) {
		controllerutil.AddFinalizer(clientMount, finalizerClientMount)
		if err := r.Update(ctx, clientMount); err != nil {
			return ctrl.Result{Requeue: true}, nil
		}

		return ctrl.Result{}, nil
	}

	if clientMount.Spec.DesiredState == "mounted" {
		err := r.mountAll(ctx, clientMount)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if clientMount.Spec.DesiredState == "unmounted" {
		err := r.unmountAll(ctx, clientMount)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// unmountAll unmounts all the file systems listed in the spec.Mounts list
func (r *ClientMountReconciler) unmountAll(ctx context.Context, clientMount *dwsv1alpha1.ClientMount) error {
	log := r.Log.WithValues("ClientMount", types.NamespacedName{Name: clientMount.Name, Namespace: clientMount.Namespace})

	var firstError error = nil
	for i, mount := range clientMount.Spec.Mounts {
		err := r.unmount(ctx, mount, log)
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			clientMount.Status.Mounts[i].Message = err.Error()
		} else {
			clientMount.Status.Mounts[i].Message = ""
			clientMount.Status.Mounts[i].Ready = true
		}
	}

	return firstError
}

// unmount unmounts a single mount point described in the ClientMountInfo object
func (r *ClientMountReconciler) unmount(ctx context.Context, clientMountInfo dwsv1alpha1.ClientMountInfo, log logr.Logger) error {
	mounted, err := r.checkMount(clientMountInfo.MountPath)
	if err != nil {
		return err
	}

	if mounted == false {
		return nil
	}

	output, err := r.run("umount " + clientMountInfo.MountPath)

	if err != nil {
		log.Info("Could not unmount file system", "mount path", clientMountInfo.MountPath, "Error output", output)
		return err
	}

	log.Info("Unmounted file system", "mount path", clientMountInfo.MountPath)
	return nil
}

// mountAll mounts all the file systems listed in the spec.Mounts list
func (r *ClientMountReconciler) mountAll(ctx context.Context, clientMount *dwsv1alpha1.ClientMount) error {
	log := r.Log.WithValues("ClientMount", types.NamespacedName{Name: clientMount.Name, Namespace: clientMount.Namespace})

	var firstError error = nil
	for i, mount := range clientMount.Spec.Mounts {
		err := r.mount(ctx, mount, log)
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			clientMount.Status.Mounts[i].Message = err.Error()
		} else {
			clientMount.Status.Mounts[i].Message = ""
			clientMount.Status.Mounts[i].Ready = true
		}
	}

	return firstError
}

// mount mounts a single mount point described in the ClientMountInfo object
func (r *ClientMountReconciler) mount(ctx context.Context, clientMountInfo dwsv1alpha1.ClientMountInfo, log logr.Logger) error {

	// Check whether the file system is already mounted
	mounted, err := r.checkMount(clientMountInfo.MountPath)
	if err != nil {
		return err
	}

	if mounted == true {
		log.Info("Already mounted")
		return nil
	}

	device, err := r.getDevice(clientMountInfo)
	if err != nil {
		return err
	}

	// Create the mount directory
	output, err := r.run("mkdir -p " + clientMountInfo.MountPath)
	if err != nil {
		log.Error(err, "Could not create mount directory", "mount path", clientMountInfo.MountPath, "device", device, "Error output", output)
		return err
	}

	// Run the mount command
	mountCmd := "mount -t " + clientMountInfo.Type + " " + device + " " + clientMountInfo.MountPath
	output, err = r.run(mountCmd)
	if err != nil {
		log.Info("Could not unmount file system", "mount path", clientMountInfo.MountPath, "device", device, "Error output", output)
		return err
	}

	log.Info("Mounted file system", "Mount path", clientMountInfo.MountPath, "device", device)

	return nil
}

// getDevice builds the device string for the mount command. This is dependent on the type of file
func (r *ClientMountReconciler) getDevice(clientMountInfo dwsv1alpha1.ClientMountInfo) (string, error) {
	switch clientMountInfo.Device.Type {
	case "lustre":
		var device string = ""

		for _, mgsAddress := range clientMountInfo.Device.Lustre.MgsAddresses {
			device = device + mgsAddress + ":"
		}

		device = device + "/" + clientMountInfo.Device.Lustre.FileSystemName

		return device, nil
	}

	return "", nil
}

// checkMount checks whether a file system is mounted at the path specified= in "mountPath"
func (r *ClientMountReconciler) checkMount(mountPath string) (bool, error) {
	output, err := r.run("mount")
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			if fields[2] == mountPath {
				return true, nil
			}
		}
	}

	return false, nil
}

// run runs a command on the host OS and returns the output as a string.
func (r *ClientMountReconciler) run(c string) (string, error) {
	if r.Mock {
		r.Log.Info("Run", "Command", c)
		return "", nil
	}

	output, err := exec.Command("bash", "-c", c).Output()

	return string(output), err
}

type clientMountStatusUpdater struct {
	clientMount    *dwsv1alpha1.ClientMount
	existingStatus dwsv1alpha1.ClientMountStatus
}

func newClientMountStatusUpdater(c *dwsv1alpha1.ClientMount) *clientMountStatusUpdater {
	return &clientMountStatusUpdater{
		clientMount:    c,
		existingStatus: (*c.DeepCopy()).Status,
	}
}

func (c *clientMountStatusUpdater) close(ctx context.Context, r *ClientMountReconciler) error {
	if !reflect.DeepEqual(c.clientMount.Status, c.existingStatus) {
		err := r.Status().Update(ctx, c.clientMount)
		if !apierrors.IsConflict(err) {
			return err
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClientMountReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dwsv1alpha1.ClientMount{}).
		Complete(r)
}
