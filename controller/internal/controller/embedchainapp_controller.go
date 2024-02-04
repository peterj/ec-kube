/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	aiappsv1 "learncloudnative.com/aiapps/api/v1"
)

// EmbedchainAppReconciler reconciles a EmbedchainApp object
type EmbedchainAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Image name constant
const (
	EC_APP_IMAGE_NAME = "ec-image:latest"
	FINALIZER_NAME    = "finalizer.aiapps.learncloudnative.com"
)

func GetIntPointer(value int32) *int32 {
	return &value
}

//+kubebuilder:rbac:groups=aiapps.learncloudnative.com,resources=embedchainapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aiapps.learncloudnative.com,resources=embedchainapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aiapps.learncloudnative.com,resources=embedchainapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EmbedchainApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *EmbedchainAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Reconciling EmbedchainApp")

	appInstance := &aiappsv1.EmbedchainApp{}
	err := r.Get(ctx, req.NamespacedName, appInstance)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			l.Info("EmbedchainApp instance not found - returning nil", "Name", req.Name, "Namespace", req.Namespace)
			// Request object not found, could have been deleted after reconcile request
			// Object was deleted, nothing to do
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request
		l.Error(err, "Failed to get EmbedchainApp instance")
		return ctrl.Result{}, err
	}

	// Check if the object is being deleted and if so, handle it
	if appInstance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(appInstance, FINALIZER_NAME) {
			controllerutil.AddFinalizer(appInstance, FINALIZER_NAME)
			err = r.Update(ctx, appInstance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(appInstance, FINALIZER_NAME) {
			// Run finalization logic for EmbedchainApp
			// Once finalization is done, remove the finalizer
			if err := r.finalizeEmbedchainApp(ctx, appInstance); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(appInstance, FINALIZER_NAME)
			err = r.Update(ctx, appInstance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Create a Deployment
	_, err = r.createDeployment(ctx, req, appInstance)
	if err != nil {
		l.Error(err, "Failed to create Deployment")
		return ctrl.Result{}, err
	}

	_, err = r.createService(ctx, req, appInstance)
	if err != nil {
		l.Error(err, "Failed to create Service")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *EmbedchainAppReconciler) finalizeEmbedchainApp(ctx context.Context, appInstance *aiappsv1.EmbedchainApp) error {
	l := log.FromContext(ctx)
	l.Info("Finalizing EmbedchainApp", "Name", appInstance.Name, "Namespace", appInstance.Namespace)

	// Delete the Deployment
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: appInstance.Namespace, Name: appInstance.Name}, deployment)
	if err != nil {
		l.Error(err, "Failed to get Deployment")
		return nil
	}

	err = r.Delete(ctx, deployment)
	if err != nil {
		l.Error(err, "Failed to delete Deployment")
		return nil
	}

	l.Info("Deleted Deployment", "Name", deployment.Name, "Namespace", deployment.Namespace)

	// Delete the Service
	service := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Namespace: appInstance.Namespace, Name: appInstance.Name}, service)
	if err != nil {
		l.Error(err, "Failed to get Service")
		return nil
	}

	err = r.Delete(ctx, service)
	if err != nil {
		l.Error(err, "Failed to delete Service")
		return nil
	}

	l.Info("Deleted Service", "Name", service.Name, "Namespace", service.Namespace)

	return nil
}

func (r *EmbedchainAppReconciler) createService(ctx context.Context, req ctrl.Request, appInstance *aiappsv1.EmbedchainApp) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Creating Service", "Name", req.Name, "Namespace", req.Namespace)

	k8sService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app": req.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": req.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.FromInt(8000),
				},
			},
		},
	}

	// Set EmbedchainApp instance as the owner and controller
	if err := controllerutil.SetControllerReference(appInstance, k8sService, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	found := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, found)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			l.Info("Creating a new Service", "Service.Namespace", k8sService.Namespace, "Service.Name", k8sService.Name)
			err = r.Create(ctx, k8sService)
			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		l.Info("Failed to get Service", "Service.Namespace", k8sService.Namespace, "Service.Name", k8sService.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *EmbedchainAppReconciler) createDeployment(ctx context.Context, req ctrl.Request, appInstance *aiappsv1.EmbedchainApp) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Creating Deployment", "Name", req.Name, "Namespace", req.Namespace)

	k8sDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			// TODO: this should be a random string
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app": req.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: GetIntPointer(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": req.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": req.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "ec-api",
							Image:           EC_APP_IMAGE_NAME,
							ImagePullPolicy: corev1.PullNever,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "OPENAI_API_KEY",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: appInstance.Spec.SecretRef.Name,
											},
											Key: "OPENAI_API_KEY",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config-volume",
									MountPath: "/app/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: appInstance.Spec.ConfigRef.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Set EmbedchainApp instance as the owner and controller
	if err := controllerutil.SetControllerReference(appInstance, k8sDeployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	found := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, found)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			l.Info("Creating a new Deployment", "Deployment.Namespace", k8sDeployment.Namespace, "Deployment.Name", k8sDeployment.Name)
			err = r.Create(ctx, k8sDeployment)
			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		l.Info("Failed to get Deployment", "Deployment.Namespace", k8sDeployment.Namespace, "Deployment.Name", k8sDeployment.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *EmbedchainAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiappsv1.EmbedchainApp{}).
		Complete(r)
}
