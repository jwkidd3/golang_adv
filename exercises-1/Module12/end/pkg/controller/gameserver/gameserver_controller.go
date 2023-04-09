package gameserver

import (
	"context"
	"fmt"
	"time"

	gameserverv1 "github.com/motiso/gameserveroperator/pkg/apis/gameserver/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_gameserver")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Gameserver Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileGameserver{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("gameserver-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Gameserver
	err = c.Watch(&source.Kind{Type: &gameserverv1.Gameserver{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner VisitorsApp
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &gameserverv1.Gameserver{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &gameserverv1.Gameserver{},
	})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Gameserver
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &gameserverv1.Gameserver{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileGameserver implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileGameserver{}

// ReconcileGameserver reconciles a Gameserver object
type ReconcileGameserver struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Gameserver object and makes changes based on the state read
// and what is in the Gameserver.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileGameserver) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Gameserver")

	// Fetch the Gameserver instance
	instance := &gameserverv1.Gameserver{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var result *reconcile.Result

	// == MySQL ==========
	result, err = r.ensureSecret(request, instance, r.mysqlAuthSecret(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureDeployment(request, instance, r.mysqlDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(request, instance, r.mysqlService(instance))
	if result != nil {
		return *result, err
	}

	mysqlRunning := r.isMysqlUp(instance)

	//we need mysql to be running first
	if !mysqlRunning {
		// If MySQL isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		log.Info(fmt.Sprintf("MySQL isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == gameserver Backend  ==========
	result, err = r.ensureLatestConfigMap(request, instance, r.newConfigMap(instance))
	if err != nil {
		return *result, err
	}

	result, err = r.ensureDeployment(request, instance, r.backendDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(request, instance, r.backendService(instance))
	if result != nil {
		return *result, err
	}

	err = r.updateBackendStatus(instance)
	if err != nil {
		// Requeue the request if the status could not be updated
		return reconcile.Result{}, err
	}

	result, err = r.handleBackendChanges(instance)
	if result != nil {
		return *result, err
	}

	// == Finish ==========
	// Everything went fine, don't requeue
	return reconcile.Result{}, nil
}
