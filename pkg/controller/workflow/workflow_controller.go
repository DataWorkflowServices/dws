package workflow

import (
	"context"
	"strings"
	myerror "errors"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	dwsv1alpha1 "stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
)

// Define condtion values
const (
    ConditionTrue bool = true
    ConditionFalse bool = false
)

// Define valid state transition values
const (
	StateProposal string = "proposal"
	StateSetup string = "setup"
	StatePreRun string = "pre_run"
	StatePostRun string = "post_run"
	StateDataIn string = "data_in"
	StateDataOut string = "data_out"
	StateTearDown string = "teardown"
)

var log = logf.Log.WithName("controller_workflow")

// Add creates a new Workflow Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileWorkflow{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("workflow-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Workflow
	err = c.Watch(&source.Kind{Type: &dwsv1alpha1.Workflow{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileWorkflow implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileWorkflow{}

// ReconcileWorkflow reconciles a Workflow object
type ReconcileWorkflow struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// checkDriverStatus returns true if all registered drivers for the current state completed successfully
func checkDriverStatus(instance *dwsv1alpha1.Workflow) (bool, error) {
	for _, d := range instance.Status.Drivers {
		if d.WatchState == instance.Spec.DesiredState {
			if (strings.ToLower(d.Reason) == "error") {
				// Return errors
				return ConditionTrue, myerror.New(d.Message)
			}
			if d.Completed == ConditionFalse {
				// Return not ready
				return ConditionFalse, nil
			} 
		} 
	}
	return ConditionTrue, nil
}

// Reconcile reads the state of the cluster for a Workflow object and makes changes based on the state read
// and what is in the Workflow.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileWorkflow) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Workflow")

	// Fetch the Workflow instance
	instance := &dwsv1alpha1.Workflow{}

	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Error(err, "Workflow instance not found")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Could not get instance Workflow")
		return reconcile.Result{}, err
	}

	existing := &dwsv1alpha1.Workflow{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, existing)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Error(err, "Workflow existing not found")
		return reconcile.Result{}, err
	} else if err != nil {
		reqLogger.Error(err, "Could not get existing Workflow")
		return reconcile.Result{}, err
	}

	driverDone, err := checkDriverStatus(instance)
	if err != nil {
		reqLogger.Info("Workflow state transitioning to " + "ERROR")
	    instance.Status.State = instance.Spec.DesiredState
	    instance.Status.Ready = ConditionFalse
	    instance.Status.Reason = "ERROR" 
	    instance.Status.Message = err.Error()
	} else {
		reqLogger.Info("Workflow state transitioning to " + instance.Spec.DesiredState)
	    instance.Status.State = instance.Spec.DesiredState
		if (driverDone == ConditionTrue) {
			instance.Status.Ready = ConditionTrue
		    instance.Status.Reason = "Completed" 
		    instance.Status.Message = "Workflow " + instance.Status.State + " completed successfully"
		} else {
			instance.Status.Ready = ConditionFalse
		    instance.Status.Reason = "DriverWait" 
		    instance.Status.Message = "Workflow " + instance.Status.State + " waiting for driver completion"
		}
	}
	err = r.client.Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Failed to update Workflow state")
		return reconcile.Result{}, err
	}
	reqLogger.Info("Status was updated", "State", instance.Status.State)
	reqLogger.Info("State to/from", "Existing", existing.Status.State, "Desired", instance.Status.State)

	return reconcile.Result{}, nil
}
