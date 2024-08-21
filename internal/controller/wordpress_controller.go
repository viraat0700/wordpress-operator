package controller

import (
	"context"
	"fmt"
	"time"

	// appsv1 "k8s.io/api/apps/v1"

	"github.com/go-logr/logr"
	v1 "github.com/vyas-git/wordpress-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// corev1 "k8s.io/api/core/v1"
	// batchv1 "k8s.io/api/batch/v1"
)

// WordpressReconciler reconciles a Wordpress object
type WordpressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile function where you manage the WordPress and MySQL resources
func (r *WordpressReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("wordpress", request.NamespacedName)

	r.Log.Info("Reconciling Wordpress")

	wordpress := &v1.Wordpress{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, wordpress)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var result *ctrl.Result

	// MySQL resources
	result, err = r.ensurePVC(request, wordpress, r.pvcForMysql(wordpress))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(request, wordpress, r.deploymentForMysql(wordpress))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureService(request, wordpress, r.serviceForMysql(wordpress))
	if result != nil {
		return *result, err
	}

	mysqlRunning := r.isMysqlUp(wordpress)
	if !mysqlRunning {
		delay := time.Second * 5
		r.Log.Info(fmt.Sprintf("MySQL isn't running, waiting for %s", delay))
		return ctrl.Result{RequeueAfter: delay}, nil
	}

	// WordPress resources
	result, err = r.ensurePVC(request, wordpress, r.pvcForWordpress(wordpress))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(request, wordpress, r.deploymentForWordpress(wordpress))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureService(request, wordpress, r.serviceForWordpress(wordpress))
	if result != nil {
		return *result, err
	}

	// Backup logic
	result, err = r.ensureBackupPVC(request, wordpress, r.pvcForBackup(wordpress))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureCronJob(request, wordpress, r.cronJobForMysqlBackup(wordpress))
	if result != nil {
		return *result, err
	}

	return ctrl.Result{}, nil
}
func (r *WordpressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Wordpress{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&batchv1.CronJob{}).              // Watches for CronJob resources
		Owns(&corev1.PersistentVolumeClaim{}). // Watches for PVCs, including backup PVCs
		Complete(r)
}
