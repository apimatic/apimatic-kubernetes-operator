/*
Copyright 2021 APIMatic.io.

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

	"reflect"

	patch "github.com/banzaicloud/k8s-objectmatcher/patch"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apicodegenv1beta1 "github.com/apimatic/apimatic-kubernetes-operator/api/v1beta1"
)

// APIMaticReconciler reconciles a APIMatic object
type APIMaticReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	LICENSEPATH   = "/usr/apimatic/license"
	LICENSEVOLUME = "licensevolume"
)

//+kubebuilder:rbac:groups=apicodegen.apimatic.io,resources=apimatics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apicodegen.apimatic.io,resources=apimatics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apicodegen.apimatic.io,resources=apimatics/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;patch;delete
//+kubebuilder:rbac:groups=core,resources=services/status,verbs=get
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the APIMatic object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *APIMaticReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// your logic here
	// Fetch the APIMatic instance
	apimatic := &apicodegenv1beta1.APIMatic{}
	err := r.Get(ctx, req.NamespacedName, apimatic)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("APIMatic resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		log.Error(err, "Failed to get APIMatic")
		return ctrl.Result{}, err
	}

	// Check if service already exists, if not create a new one
	foundService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: apimatic.Name, Namespace: apimatic.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Service
		dep := r.serviceForAPIMatic(apimatic)
		if err = patch.DefaultAnnotator.SetLastAppliedAnnotation(dep); err != nil {
			log.Error(err, "Error setting annotations on Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
			return ctrl.Result{}, err
		}
		log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
		ctrl.SetControllerReference(apimatic, dep, r.Scheme)
		err = r.Create(ctx, dep)
		if err != nil && errors.IsAlreadyExists(err) {
			log.Error(err, "Waiting to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Service created successfully- return and requeue
		log.Info("Successfully created new service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// Check if service needs to be updated according to APIMatic Service specifications and then update if this is needed
	foundService, needsUpdate, err := r.shouldUpdateServiceForAPIMatic(apimatic, foundService)
	if err != nil {
		log.Error(err, "Error in checking if Service update needed", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)
		return ctrl.Result{}, err
	} else {
		if needsUpdate {
			log.Info("Deleting outdated service for APIMatic instance", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)
			ctrl.SetControllerReference(apimatic, foundService, r.Scheme)
			err = r.Delete(ctx, foundService)
			if err != nil {
				log.Error(err, "Failure deleting outdated service for APIMatic instance", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)
				return ctrl.Result{}, err
			} else {
				log.Info("Successfully deleted service for APIMatic instance", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)
				return ctrl.Result{Requeue: true}, nil
			}
		}
	}

	// Check if deployment already exists, if not create a new one
	foundDeployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: apimatic.Name, Namespace: apimatic.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForAPIMatic(apimatic)
		if err = patch.DefaultAnnotator.SetLastAppliedAnnotation(dep); err != nil {
			log.Error(err, "Error setting annotations on deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		ctrl.SetControllerReference(apimatic, dep, r.Scheme)
		err = r.Create(ctx, dep)

		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		// Deployment created successfully- return and requeue
		log.Info("Successfully created new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment", "Deployment.Namespace", apimatic.Namespace, "Deployment.Name", apimatic.Name)
		return ctrl.Result{}, err
	}

	// Check if deployment needs to be updated according to APIMatic spec and update if needed
	foundDeployment, needsUpdate, err = r.shouldUpdateDeploymentForAPIMatic(apimatic, foundDeployment)
	if err != nil {
		log.Error(err, "Error in checking if Deployment update needed", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
		return ctrl.Result{}, err
	} else {
		if needsUpdate {
			log.Info("Updating Deployment for APIMatic instance", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			ctrl.SetControllerReference(apimatic, foundDeployment, r.Scheme)
			err = r.Update(ctx, foundDeployment)

			if err != nil {
				log.Error(err, "Failure updating Deployment for APIMatic instance", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
				return ctrl.Result{}, err
			} else {
				log.Info("Successfully updated Deployment for APIMatic instance", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
				return ctrl.Result{Requeue: true}, nil
			}
		}
	}

	var foundServiceStatus = foundService.Status
	var foundDeploymentStatus = foundDeployment.Status

	if !reflect.DeepEqual(foundServiceStatus, apimatic.Status.ServiceStatus) || !reflect.DeepEqual(foundDeploymentStatus, apimatic.Status.DeploymentStatus) {
		apimatic.Status.ServiceStatus = foundServiceStatus
		apimatic.Status.DeploymentStatus = foundDeploymentStatus

		err = r.Status().Update(ctx, apimatic)
		if err != nil {
			log.Error(err, "Failed to update APIMatic status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func labelsForAPIMatic(name string) map[string]string {
	return map[string]string{"app": "apimatic", "apimatic_cr": name}
}

func (r *APIMaticReconciler) serviceForAPIMatic(a *apicodegenv1beta1.APIMatic) *corev1.Service {
	ls := labelsForAPIMatic(a.Name)

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{{
				Port:       a.Spec.ServiceSpec.APIMaticServicePort.Port,
				TargetPort: intstr.FromInt(80),
				Name:       a.Spec.ServiceSpec.APIMaticServicePort.Name,
			}},
		},
	}

	dep.Spec.Type = a.Spec.ServiceSpec.Type

	dep.Spec.SessionAffinity = a.Spec.ServiceSpec.SessionAffinity

	dep.Spec.PublishNotReadyAddresses = a.Spec.ServiceSpec.PublishNotReadyAddresses

	if reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeLoadBalancer) || reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeNodePort) {
		if a.Spec.ServiceSpec.APIMaticServicePort.NodePort != nil {
			dep.Spec.Ports[0].NodePort = *a.Spec.ServiceSpec.APIMaticServicePort.NodePort
		}
	}

	if reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeLoadBalancer) {
		if a.Spec.ServiceSpec.LoadBalancerIP != nil {
			dep.Spec.LoadBalancerIP = *a.Spec.ServiceSpec.LoadBalancerIP
		}
	}

	if reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeLoadBalancer) || reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeNodePort) {
		if a.Spec.ServiceSpec.ExternalTrafficPolicy != nil {
			dep.Spec.ExternalTrafficPolicy = *a.Spec.ServiceSpec.ExternalTrafficPolicy
		} else {
			dep.Spec.ExternalTrafficPolicy = corev1.ServiceExternalTrafficPolicyTypeCluster
		}
	}

	if reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeLoadBalancer) && reflect.DeepEqual(dep.Spec.ExternalTrafficPolicy, corev1.ServiceExternalTrafficPolicyTypeLocal) {
		if a.Spec.ServiceSpec.HealthCheckNodePort != nil {
			dep.Spec.HealthCheckNodePort = *a.Spec.ServiceSpec.HealthCheckNodePort
		}
	}

	if a.Spec.ServiceSpec.SessionAffinityConfig != nil {
		dep.Spec.SessionAffinityConfig = a.Spec.ServiceSpec.SessionAffinityConfig
	}

	if a.Spec.ServiceSpec.TopologyKeys != nil {
		dep.Spec.TopologyKeys = []string{}
		dep.Spec.TopologyKeys = append(dep.Spec.TopologyKeys, a.Spec.ServiceSpec.TopologyKeys...)
	}

	if a.Spec.ServiceSpec.IPFamilyPolicy != nil && !reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeExternalName) {
		dep.Spec.IPFamilyPolicy = a.Spec.ServiceSpec.IPFamilyPolicy
	} else {
		dep.Spec.IPFamilyPolicy = new(corev1.IPFamilyPolicyType)
		*dep.Spec.IPFamilyPolicy = corev1.IPFamilyPolicySingleStack
	}

	if (reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeClusterIP) || reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeLoadBalancer) || reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeNodePort)) && (!reflect.DeepEqual(*dep.Spec.IPFamilyPolicy, corev1.IPFamilyPolicySingleStack)) {
		if a.Spec.ServiceSpec.IPFamilies != nil {
			dep.Spec.IPFamilies = []corev1.IPFamily{}
			dep.Spec.IPFamilies = append(dep.Spec.IPFamilies, a.Spec.ServiceSpec.IPFamilies...)
		}
	}

	if reflect.DeepEqual(dep.Spec.Type, corev1.ServiceTypeExternalName) {
		if a.Spec.ServiceSpec.ExternalName != nil {
			dep.Spec.ExternalName = *a.Spec.ServiceSpec.ExternalName
		}
	}

	return dep
}

func (r *APIMaticReconciler) shouldUpdateServiceForAPIMatic(a *apicodegenv1beta1.APIMatic, s *corev1.Service) (*corev1.Service, bool, error) {
	newService := r.serviceForAPIMatic(a)

	opts := []patch.CalculateOption{
		patch.IgnoreStatusFields(),
	}
	patchResult, err := patch.DefaultPatchMaker.Calculate(s, newService, opts...)
	if err != nil {
		return newService, false, err
	}

	if !patchResult.IsEmpty() {
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(newService); err != nil {
			return newService, true, err
		} else {
			return newService, true, nil
		}
	}
	return s, false, nil
}

func (r *APIMaticReconciler) deploymentForAPIMatic(a *apicodegenv1beta1.APIMatic) *appsv1.Deployment {
	ls := labelsForAPIMatic(a.Name)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &a.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &a.Spec.PodSpec.TerminationGracePeriodSeconds,
					Containers: []corev1.Container{{
						Image: a.Spec.PodSpec.APIMaticContainerSpec.Image,
						Name:  "apimatic",
						Env: []corev1.EnvVar{{
							Name:  "LICENSEPATH",
							Value: LICENSEPATH,
						}},
						VolumeMounts: []corev1.VolumeMount{{
							ReadOnly:  true,
							Name:      LICENSEVOLUME,
							MountPath: LICENSEPATH,
						}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
						}},
					}},
					Volumes: []corev1.Volume{
						*volumeSource(a),
					},
				},
			},
		},
	}

	dep.Spec.Template.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{{
				Weight: 100,
				PodAffinityTerm: corev1.PodAffinityTerm{
					TopologyKey: "anti-affinity",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: ls,
					},
				},
			}},
		},
	}

	dep.Spec.Template.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{{
				Weight: 100,
				PodAffinityTerm: corev1.PodAffinityTerm{
					TopologyKey: "anti-affinity",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: ls,
					},
				},
			}},
		},
	}

	if a.Spec.PodSpec.SecurityContext != nil {
		dep.Spec.Template.Spec.SecurityContext = a.Spec.PodSpec.SecurityContext
	}

	if a.Spec.PodSpec.APIMaticContainerSpec.ImagePullSecret != nil {
		dep.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{}
		imagePullSecret := corev1.LocalObjectReference{
			Name: *a.Spec.PodSpec.APIMaticContainerSpec.ImagePullSecret,
		}
		dep.Spec.Template.Spec.ImagePullSecrets = append(dep.Spec.Template.Spec.ImagePullSecrets, imagePullSecret)
	}

	if a.Spec.PodSpec.Hostname != nil {
		dep.Spec.Template.Spec.Hostname = *a.Spec.PodSpec.Hostname
	}

	if a.Spec.PodSpec.Subdomain != nil {
		dep.Spec.Template.Spec.Subdomain = *a.Spec.PodSpec.Subdomain
	}

	if a.Spec.PodSpec.SchedulerName != nil {
		dep.Spec.Template.Spec.SchedulerName = *a.Spec.PodSpec.SchedulerName
	}

	if a.Spec.PodSpec.HostAliases != nil {
		dep.Spec.Template.Spec.HostAliases = []corev1.HostAlias{}
		dep.Spec.Template.Spec.HostAliases = append(dep.Spec.Template.Spec.HostAliases, a.Spec.PodSpec.HostAliases...)
	}

	if a.Spec.PodSpec.PriorityClassName != nil {
		dep.Spec.Template.Spec.PriorityClassName = *a.Spec.PodSpec.PriorityClassName
	}

	if a.Spec.PodSpec.Priority != nil {
		dep.Spec.Template.Spec.Priority = a.Spec.PodSpec.Priority
	}

	if a.Spec.PodSpec.ReadinessGates != nil {
		dep.Spec.Template.Spec.ReadinessGates = []corev1.PodReadinessGate{}
		dep.Spec.Template.Spec.ReadinessGates = append(dep.Spec.Template.Spec.ReadinessGates, a.Spec.PodSpec.ReadinessGates...)
	}

	dep.Spec.Template.Spec.EnableServiceLinks = &a.Spec.PodSpec.EnableServiceLinks

	dep.Spec.Template.Spec.SetHostnameAsFQDN = &a.Spec.PodSpec.SetHostnameAsFQDN

	dep.Spec.Template.Spec.RestartPolicy = a.Spec.PodSpec.RestartPolicy

	dep.Spec.Template.Spec.DNSPolicy = a.Spec.PodSpec.DNSPolicy

	dep.Spec.Template.Spec.HostNetwork = a.Spec.PodSpec.HostNetwork

	dep.Spec.Template.Spec.HostPID = a.Spec.PodSpec.HostPID

	dep.Spec.Template.Spec.HostIPC = a.Spec.PodSpec.HostIPC

	if a.Spec.PodSpec.ServiceAccountName != nil {
		dep.Spec.Template.Spec.ServiceAccountName = *a.Spec.PodSpec.ServiceAccountName
	}

	if a.Spec.PodSpec.AutomountServiceAccountToken != nil {
		dep.Spec.Template.Spec.AutomountServiceAccountToken = a.Spec.PodSpec.AutomountServiceAccountToken
	}

	if a.Spec.PodSpec.DNSConfig != nil {
		dep.Spec.Template.Spec.DNSConfig = a.Spec.PodSpec.DNSConfig
	}

	if a.Spec.PodSpec.ActiveDeadlineSeconds != nil {
		dep.Spec.Template.Spec.ActiveDeadlineSeconds = a.Spec.PodSpec.ActiveDeadlineSeconds
	}

	if a.Spec.PodSpec.APIMaticContainerSpec.ImagePullPolicy != nil {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = *a.Spec.PodSpec.APIMaticContainerSpec.ImagePullPolicy
	}

	if a.Spec.PodSpec.APIMaticContainerSpec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *a.Spec.PodSpec.APIMaticContainerSpec.Resources
	}

	if a.Spec.APIMaticPodPlacementSpec != nil {
		if a.Spec.APIMaticPodPlacementSpec.NodeAffinity != nil {
			dep.Spec.Template.Spec.Affinity.NodeAffinity = a.Spec.APIMaticPodPlacementSpec.NodeAffinity
		}

		if a.Spec.APIMaticPodPlacementSpec.PodAffinity != nil {
			dep.Spec.Template.Spec.Affinity.PodAffinity = a.Spec.APIMaticPodPlacementSpec.PodAffinity
		}

		if a.Spec.APIMaticPodPlacementSpec.NodeName != nil {
			dep.Spec.Template.Spec.NodeName = *a.Spec.APIMaticPodPlacementSpec.NodeName
		}

		if a.Spec.APIMaticPodPlacementSpec.Tolerations != nil {
			dep.Spec.Template.Spec.Tolerations = []corev1.Toleration{}
			dep.Spec.Template.Spec.Tolerations = append(dep.Spec.Template.Spec.Tolerations, a.Spec.APIMaticPodPlacementSpec.Tolerations...)
		}
	}

	return dep
}

func (r *APIMaticReconciler) shouldUpdateDeploymentForAPIMatic(a *apicodegenv1beta1.APIMatic, d *appsv1.Deployment) (*appsv1.Deployment, bool, error) {
	newDeployment := r.deploymentForAPIMatic(a)
	opts := []patch.CalculateOption{
		patch.IgnoreStatusFields(),
	}

	patchResult, err := patch.DefaultPatchMaker.Calculate(d, newDeployment, opts...)

	if err != nil {
		return d, false, err
	}

	if !patchResult.IsEmpty() {
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(newDeployment); err != nil {
			return newDeployment, true, err
		} else {
			return newDeployment, true, nil
		}
	}

	return d, false, nil
}

func volumeSource(a *apicodegenv1beta1.APIMatic) *corev1.Volume {

	licenseSource := &corev1.Volume{
		Name: LICENSEVOLUME,
	}
	if reflect.DeepEqual(a.Spec.LicenseSpec.APIMaticLicenseSourceType, "ConfigMap") {
		licenseSource.VolumeSource = corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: a.Spec.LicenseSpec.APIMaticLicenseSourceName,
				},
			},
		}
	} else {
		licenseSource.VolumeSource = corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: a.Spec.LicenseSpec.APIMaticLicenseSourceName,
			},
		}
	}

	return licenseSource
}

// SetupWithManager sets up the controller with the Manager.
func (r *APIMaticReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apicodegenv1beta1.APIMatic{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
