/*
Copyright 2023.

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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1 "github.com/fazekasrobert/operator-sandbox/api/v1"
)

// DeployerReconciler reconciles a Deployer object
type DeployerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operator.github.com,resources=deployers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.github.com,resources=deployers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.github.com,resources=deployers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Deployer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DeployerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling foo custom resource")

	var deployer operatorv1.Deployer
	if err := r.Get(ctx, req.NamespacedName, &deployer); err != nil {
		log.Error(err, "unable to fetch Deployer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !deployer.Status.DeploymentOK {
		deployment := appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deployer.Name,
				Namespace: "sandbox-system",
				Labels:    map[string]string{"app": deployer.Name},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: deployer.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": deployer.Name},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": deployer.Name},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "httpd-container",
								Image: deployer.Spec.Image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
			},
		}
		if err := ctrl.SetControllerReference(&deployer, &deployment, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Creating a new Deployment", "Deployment.Name", deployment.Name)
		if err := r.Client.Create(context.TODO(), &deployment); err != nil {
			return ctrl.Result{}, err
		}

		log.Info("Deployment created successfully")
		deployer.Status.DeploymentOK = true
	}

	if !deployer.Status.ServiceOK {
		service := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deployer.Name,
				Namespace: "sandbox-system",
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Port:       80,
						TargetPort: intstr.FromInt(80),
					},
				},
				Selector: map[string]string{"app": deployer.Name},
			},
		}
		log.Info("Creating a new Service", "Service.Name", service.Name)
		if err := r.Client.Create(context.TODO(), &service); err != nil {
			return ctrl.Result{}, err
		}

		log.Info("Service created successfully")
		deployer.Status.ServiceOK = true
	}

	if !deployer.Status.IngressOK {
		prefix := networkingv1.PathType("Prefix")
		nginx := "nginx"
		ingress := networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deployer.Name + "-ingress",
				Namespace: "sandbox-system",
			},
			Spec: networkingv1.IngressSpec{
				IngressClassName: &nginx,
				Rules: []networkingv1.IngressRule{
					{
						Host: deployer.Spec.Host,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     "/",
										PathType: &prefix,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: deployer.Name,
												Port: networkingv1.ServiceBackendPort{
													Number: 80,
												},
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
		log.Info("Creating a new Ingress", "Ingress.Name", ingress.Name)
		if err := r.Client.Create(context.TODO(), &ingress); err != nil {
			return ctrl.Result{}, err
		}

		log.Info("Ingress created successfully")
		deployer.Status.IngressOK = true
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeployerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.Deployer{}).
		Complete(r)
}
