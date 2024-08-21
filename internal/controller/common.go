package controller

import (
	"context"

	v1 "github.com/vyas-git/wordpress-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *WordpressReconciler) ensureDeployment(_ reconcile.Request,
	instance *v1.Wordpress,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	found := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)

	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		r.Log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			r.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		}
		// Deployment was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the deployment not existing
		r.Log.Error(err, "Failed to get Deployment")
		return &ctrl.Result{}, err
	}

	return nil, nil

}

func (r *WordpressReconciler) ensureService(_ reconcile.Request,
	instance *v1.Wordpress,
	s *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		r.Log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.Client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the service not existing
		r.Log.Error(err, "Failed to get Service")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *WordpressReconciler) ensurePVC(_ reconcile.Request,
	instance *v1.Wordpress,
	s *corev1.PersistentVolumeClaim,
) (*reconcile.Result, error) {

	found := &corev1.PersistentVolumeClaim{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the PVC
		r.Log.Info("Creating a new PVC", "PVC.Namespace", s.Namespace, "PVC.Name", s.Name)
		err = r.Client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new PVC", "PVC.Namespace", s.Namespace, "PVC.Name", s.Name)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the pvc not existing
		r.Log.Error(err, "Failed to get PVC")
		return &ctrl.Result{}, err
	}

	return nil, nil

}

func (r *WordpressReconciler) ensureCronJob(_ reconcile.Request,
	instance *v1.Wordpress,
	cj *batchv1.CronJob,
) (*reconcile.Result, error) {

	found := &batchv1.CronJob{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cj.Name,
		Namespace: instance.Namespace,
	}, found)

	if err != nil && errors.IsNotFound(err) {

		// Create the CronJob
		r.Log.Info("Creating a new CronJob", "CronJob.Namespace", cj.Namespace, "CronJob.Name", cj.Name)
		err = r.Client.Create(context.TODO(), cj)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new CronJob", "CronJob.Namespace", cj.Namespace, "CronJob.Name", cj.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the CronJob not existing
		r.Log.Error(err, "Failed to get CronJob")
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *WordpressReconciler) ensureBackupPVC(_ reconcile.Request,
	instance *v1.Wordpress,
	pvc *corev1.PersistentVolumeClaim,
) (*reconcile.Result, error) {

	found := &corev1.PersistentVolumeClaim{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      pvc.Name,
		Namespace: instance.Namespace,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the PVC
		r.Log.Info("Creating a new Backup PVC", "PVC.Namespace", pvc.Namespace, "PVC.Name", pvc.Name)
		err = r.Client.Create(context.TODO(), pvc)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new Backup PVC", "PVC.Namespace", pvc.Namespace, "PVC.Name", pvc.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the PVC not existing
		r.Log.Error(err, "Failed to get Backup PVC")
		return &ctrl.Result{}, err
	}

	return nil, nil
}
