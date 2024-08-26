package controller

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	v1 "github.com/vyas-git/wordpress-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Function to generate a random password
func generateRandomPassword() (string, error) {
	length := 16
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

// Function to create a Kubernetes Secret with the MySQL root password
func (r *WordpressReconciler) createMysqlPasswordSecret(cr *v1.Wordpress) (*corev1.Secret, error) {
	secretName := "mysql-root-password-secret"
	namespace := cr.Namespace

	// Check if the Secret already exists
	secret := &corev1.Secret{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      secretName,
		Namespace: namespace,
	}, secret)
	if err == nil {
		// Secret already exists, return it
		return secret, nil
	}

	// Generate a new random password if the Secret doesn't exist
	password, err := generateRandomPassword()
	if err != nil {
		return nil, err
	}

	// Create the Secret
	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"password": []byte(password),
		},
	}

	// Set the owner reference so that the Secret is cleaned up when the CR is deleted
	if err := controllerutil.SetControllerReference(cr, secret, r.Scheme); err != nil {
		return nil, err
	}

	// Create the Secret in Kubernetes
	err = r.Client.Create(context.TODO(), secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// Creates a CronJob for MySQL backups

func (r *WordpressReconciler) cronJobForMysqlBackup(cr *v1.Wordpress) (*batchv1.CronJob, error) {
	labels := map[string]string{
		"app": cr.Name,
	}

	// Ensure the MySQL Secret exists
	secret, err := r.createMysqlPasswordSecret(cr)
	if err != nil {
		return nil, err
	}

	cronJob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-backup-cronjob",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: "*/5 * * * *", // Every 5 minutes
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "mysql-backup",
									Image: "mysql:5.6",
									Command: []string{
										"sh",
										"-c",
										"mysqldump -u root -p$MYSQL_ROOT_PASSWORD wordpress > /backup/wordpress_backup.sql",
									},
									Env: []corev1.EnvVar{
										{
											Name: "MYSQL_ROOT_PASSWORD",
											ValueFrom: &corev1.EnvVarSource{
												SecretKeyRef: &corev1.SecretKeySelector{
													LocalObjectReference: corev1.LocalObjectReference{
														Name: secret.Name,
													},
													Key: "password",
												},
											},
										},
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "backup-storage",
											MountPath: "/backup",
										},
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Volumes: []corev1.Volume{
								{
									Name: "backup-storage",
									VolumeSource: corev1.VolumeSource{
										PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
											ClaimName: "backup-pv-claim",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, cronJob, r.Scheme)
	return cronJob, nil
}

// func (r *WordpressReconciler) cronJobForMysqlBackup(cr *v1.Wordpress) *batchv1.CronJob {
// 	labels := map[string]string{
// 		"app": cr.Name,
// 	}

// 	cronJob := &batchv1.CronJob{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "mysql-backup-cronjob",
// 			Namespace: cr.Namespace,
// 			Labels:    labels,
// 		},
// 		Spec: batchv1.CronJobSpec{
// 			Schedule: "*/5 * * * *", // Every 5 minute
// 			JobTemplate: batchv1.JobTemplateSpec{
// 				Spec: batchv1.JobSpec{
// 					Template: corev1.PodTemplateSpec{
// 						Spec: corev1.PodSpec{
// 							Containers: []corev1.Container{
// 								{
// 									Name:  "mysql-backup",
// 									Image: "mysql:5.6",
// 									Command: []string{
// 										"sh",
// 										"-c",
// 										"mysqldump -u root -p$MYSQL_ROOT_PASSWORD wordpress > /backup/wordpress_backup.sql",
// 									},
// 									Env: []corev1.EnvVar{
// 										{
// 											Name:  "MYSQL_ROOT_PASSWORD",
// 											Value: cr.Spec.SqlRootPassword,
// 										},
// 									},
// 									VolumeMounts: []corev1.VolumeMount{
// 										{
// 											Name:      "backup-storage",
// 											MountPath: "/backup",
// 										},
// 									},
// 								},
// 							},
// 							RestartPolicy: corev1.RestartPolicyOnFailure,
// 							Volumes: []corev1.Volume{
// 								{
// 									Name: "backup-storage",
// 									VolumeSource: corev1.VolumeSource{
// 										PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
// 											ClaimName: "backup-pv-claim",
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	controllerutil.SetControllerReference(cr, cronJob, r.Scheme)
// 	return cronJob
// }

// Creates a PersistentVolumeClaim for MySQL backups
func (r *WordpressReconciler) pvcForBackup(cr *v1.Wordpress) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app": cr.Name,
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "backup-pv-claim",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse("10Gi"),
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc
}

// Creates a Deployment for MySQL
// func (r *WordpressReconciler) deploymentForMysql(cr *v1.Wordpress) *appsv1.Deployment {
// 	labels := map[string]string{
// 		"app": cr.Name,
// 	}
// 	matchlabels := map[string]string{
// 		"app":  cr.Name,
// 		"tier": "mysql",
// 	}

// 	// Setting the replicas for the MySQL deployment
// 	replicas := int32(*cr.Spec.MysqlReplicas)

// 	dep := &appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "wordpress-mysql",
// 			Namespace: cr.Namespace,
// 			Labels:    labels,
// 		},

// 		Spec: appsv1.DeploymentSpec{
// 			Replicas: &replicas, // Set the replicas here
// 			Selector: &metav1.LabelSelector{
// 				MatchLabels: matchlabels,
// 			},
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: matchlabels,
// 				},
// 				Spec: corev1.PodSpec{
// 					Containers: []corev1.Container{{
// 						Image: "mysql",
// 						Name:  "mysql",
// 						Env: []corev1.EnvVar{
// 							{
// 								Name:  "MYSQL_ROOT_PASSWORD",
// 								Value: cr.Spec.SqlRootPassword,
// 							},
// 						},
// 						Ports: []corev1.ContainerPort{{
// 							ContainerPort: 3306,
// 							Name:          "mysql",
// 						}},
// 						VolumeMounts: []corev1.VolumeMount{
// 							{
// 								Name:      "mysql-persistent-storage",
// 								MountPath: "/var/lib/mysql",
// 							},
// 						},
// 					}},
// 					Volumes: []corev1.Volume{
// 						{
// 							Name: "mysql-persistent-storage",
// 							VolumeSource: corev1.VolumeSource{
// 								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
// 									ClaimName: "mysql-pv-claim",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	controllerutil.SetControllerReference(cr, dep, r.Scheme)
// 	return dep
// }

func (r *WordpressReconciler) deploymentForMysql(cr *v1.Wordpress) (*appsv1.Deployment, error) {
	labels := map[string]string{
		"app": cr.Name,
	}
	matchlabels := map[string]string{
		"app":  cr.Name,
		"tier": "mysql",
	}

	// Setting the replicas for the MySQL deployment
	replicas := int32(*cr.Spec.MysqlReplicas)

	// Get or create the MySQL root password Secret
	secret, err := r.createMysqlPasswordSecret(cr)
	if err != nil {
		return nil, err
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wordpress-mysql",
			Namespace: cr.Namespace,
			Labels:    labels,
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas, // Set the replicas here
			Selector: &metav1.LabelSelector{
				MatchLabels: matchlabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchlabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "mysql",
						Name:  "mysql",
						Env: []corev1.EnvVar{
							{
								Name: "MYSQL_ROOT_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secret.Name,
										},
										Key: "password",
									},
								},
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "mysql-persistent-storage",
								MountPath: "/var/lib/mysql",
							},
						},
					}},
					Volumes: []corev1.Volume{
						{
							Name: "mysql-persistent-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "mysql-pv-claim",
								},
							},
						},
					},
				},
			},
		},
	}

	// Set owner reference so that deployment is cleaned up when the CR is deleted
	controllerutil.SetControllerReference(cr, dep, r.Scheme)
	return dep, nil
}

// Creates a PersistentVolumeClaim for MySQL
func (r *WordpressReconciler) pvcForMysql(cr *v1.Wordpress) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app": cr.Name,
	}

	pvc := &corev1.PersistentVolumeClaim{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-pv-claim",
			Namespace: cr.Namespace,
			Labels:    labels,
		},

		Spec: corev1.PersistentVolumeClaimSpec{

			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},

			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse("10Gi"),
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc

}

// Creates a Service for MySQL
func (r *WordpressReconciler) serviceForMysql(cr *v1.Wordpress) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	matchlabels := map[string]string{
		"app":  cr.Name,
		"tier": "mysql",
	}

	ser := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "wordpress-mysql",
			Namespace: cr.Namespace,
			Labels:    labels,
		},

		Spec: corev1.ServiceSpec{
			Selector: matchlabels,

			Ports: []corev1.ServicePort{
				{
					Port: 3306,
					Name: cr.Name,
				},
			},
			ClusterIP: "None",
		},
	}

	controllerutil.SetControllerReference(cr, ser, r.Scheme)
	return ser

}

// Checks if the MySQL deployment is up
func (r *WordpressReconciler) isMysqlUp(v *v1.Wordpress) bool {
	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      "wordpress-mysql",
		Namespace: v.Namespace,
	}, deployment)

	if err != nil {
		r.Log.Error(err, "Deployment mysql not found")
		fmt.Println(err, "Deployment mysql not found")

		return false
	}
	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false

}
