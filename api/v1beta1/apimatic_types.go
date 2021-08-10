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

package v1beta1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// APIMaticSpec defines the desired state of APIMatic
type APIMaticSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// replicas is the desired number of instances of APIMatic. Minimum is 0. Defaults to 1 if not provided
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	PodSpec APIMaticPodSpec `json:"podspec"`

	// +kubebuilder:validation:Required
	PodVolumeSpec APIMaticPodVolumeSpec `json:"volumespec"`

	// +kubebuilder:validation:Required
	ServiceSpec APIMaticServiceSpec `json:"servicespec"`

	// +kubebuilder:validation:Optional
	APIMaticPodPlacementSpec *APIMaticPodPlacementSpec `json:"podplacementspec,omitempty"`

	// volumeClaimTemplates is a list of claims that pods are allowed to reference.
	// The StatefulSet controller is responsible for mapping network identities to
	// claims in a way that maintains the identity of a pod. Every claim in
	// this list must have at least one matching (by name) volumeMount in one
	// container in the template. A claim in this list takes precedence over
	// any volumes in the template, with the same name.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
}

// APIMaticStatus defines the observed state of APIMatic
type APIMaticStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Optional
	// statefulsetStatus displays the status of the owned service resource which exposes the APIMatic pods for communication
	StatefulSetStatus appsv1.StatefulSetStatus `json:"statefulsetStatus,omitempty"`

	// +kubebuilder:validation:Optional
	// serviceStatus displays the status of the owned service resource which exposes the APIMatic pods for communication
	ServiceStatus corev1.ServiceStatus `json:"serviceStatus,omitempty"`
}

// APIMaticPodSpec contains configuration for created APIMatic pods
type APIMaticPodSpec struct {
	// APIMatic container image
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// APIMatic container name. If none given, name will be set as apimatic
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	Name *string `json:"name,omitempty"`

	// PullPolicy describes a policy for if/when to pull a container image. Valid values are Always, Never and IfNotPresent. Defaults to IfNotPresent if not provided
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Always;Never;IfNotPresent
	ImagePullPolicy *corev1.PullPolicy `json:"imagepullpolicy,omitempty"`

	// sidecars are the collection of sidecar containers in addition to the main APIMatic container
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	SideCars []corev1.Container `json:"sidecars,omitempty"`

	// More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	InitContainers []corev1.Container `json:"initcontainers,omitempty"`

	// Resource Requirements represents the compute resource requirements of the APIMatic container
	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request.
	// Value must be non-negative integer. The value zero indicates delete immediately.
	// If this value is nil, the default grace period will be used instead.
	// The grace period is the duration in seconds after the processes running in the pod are sent
	// a termination signal and the time when the processes are forcibly halted with a kill signal.
	// Set this value longer than the expected cleanup time for your process.
	// Defaults to 30 seconds.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`

	// Optional duration in seconds the pod may be active on the node relative to
	// StartTime before the system will actively try to mark it failed and kill associated containers.
	// Value must be a positive integer.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:ExclusiveMinimum=true
	// +kubebuilder:validation:Minimum=0
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty"`

	// Set DNS policy for the pod.
	// Defaults to "ClusterFirst".
	// Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.
	// DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy.
	// To have DNS options set along with hostNetwork, you have to specify DNS policy
	// explicitly to 'ClusterFirstWithHostNet'.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=ClusterFirstWithHostNet;ClusterFirst;Default;None
	DNSPolicy *corev1.DNSPolicy `json:"dnsPolicy,omitempty"`

	// Specifies the DNS parameters of a pod.
	// Parameters specified here will be merged to the generated DNS
	// configuration based on DNSPolicy.
	// +kubebuilder:validation:Optional
	DNSConfig *corev1.PodDNSConfig `json:"dnsConfig,omitempty"`

	// Host networking requested for this pod. Use the host's network namespace.
	// If this option is set, the ports that will be used must be specified.
	// Default to false.
	// +kubebuilder:validation:Optional
	HostNetwork *bool `json:"hostNetwork,omitempty"`

	// Restart policy for all containers within the pod.
	// One of Always, OnFailure, Never.
	// Default to Always.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Always;OnFailure;Never
	RestartPolicy *corev1.RestartPolicy `json:"restartPolicy,omitempty"`

	// ServiceAccountName is the name of the ServiceAccount to use to run this pod.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	ServiceAccountName *string `json:"serviceAccountName,omitempty"`

	// AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.
	// +kubebuilder:validation:Optional
	AutomountServiceAccountToken *bool `json:"automountServiceAccountToken,omitempty"`

	// Use the host's pid namespace.
	// Optional: Default to false.
	// +kubebuilder:validation:Optional
	HostPID *bool `json:"hostPID,omitempty"`

	// Use the host's ipc namespace.
	// Optional: Default to false.
	// +kubebuilder:validation:Optional
	HostIPC *bool `json:"hostIPC,omitempty"`

	// Share a single process namespace between all of the containers in a pod.
	// When this is set containers will be able to view and signal processes from other containers
	// in the same pod, and the first process in each container will not be assigned PID 1.
	// HostPID and ShareProcessNamespace cannot both be set.
	// Optional: Default to false.
	// +kubebuilder:validation:Optional
	ShareProcessNamespace *bool `json:"shareProcessNamespace,omitempty"`

	// SecurityContext holds pod-level security attributes and common container settings.
	// Optional: Defaults to empty.  See type description for default values of each field.
	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Specifies the hostname of the Pod
	// If not specified, the pod's hostname will be set to a system-defined value.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	Hostname *string `json:"hostname,omitempty"`

	// If specified, the fully qualified Pod hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>".
	// If not specified, the pod will not have a domainname at all.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	Subdomain *string `json:"subdomain,omitempty"`

	// If specified, the pod will be dispatched by specified scheduler.
	// If not specified, the pod will be dispatched by default scheduler.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	SchedulerName *string `json:"schedulerName,omitempty"`

	// HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts
	// file if specified. This is only valid for non-hostNetwork pods.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`

	// If specified, indicates the pod's priority. "system-node-critical" and
	// "system-cluster-critical" are two special keywords which indicate the
	// highest priorities with the former being the highest priority. Any other
	// name must be defined by creating a PriorityClass object with that name.
	// If not specified, the pod priority will be default or zero if there is no
	// default.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	PriorityClassName *string `json:"priorityClassName,omitempty"`

	// The priority value. Various system components use this field to find the
	// priority of the pod. When Priority Admission Controller is enabled, it
	// prevents users from setting this field. The admission controller populates
	// this field from PriorityClassName.
	// The higher the value, the higher the priority.
	// +kubebuilder:validation:Optional
	Priority *int32 `json:"priority,omitempty"`

	// If specified, all readiness gates will be evaluated for pod readiness.
	// A pod is ready when all its containers are ready AND
	// all conditions specified in the readiness gates have status equal to "True"
	// More info: https://git.k8s.io/enhancements/keps/sig-network/0007-pod-ready%2B%2B.md
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	ReadinessGates []corev1.PodReadinessGate `json:"readinessGates,omitempty"`

	// EnableServiceLinks indicates whether information about services should be injected into pod's
	// environment variables, matching the syntax of Docker links.
	// Optional: Defaults to true.
	// +kubebuilder:validation:Optional
	EnableServiceLinks *bool `json:"enableServiceLinks,omitempty"`

	// If a pod does not have FQDN, this has no effect.
	// Default to false.
	// +kubebuilder:validation:Optional
	SetHostnameAsFQDN *bool `json:"setHostnameAsFQDN,omitempty"`
}

// APIMaticVolumeSpec contains configuration for volumes associated with created APIMatic pods
type APIMaticPodVolumeSpec struct {

	// The license path which will be used to volume mount the license file. If not provided, the license path is set as /usr/local/apimatic
	// kubebuilder:validation:Optional
	// kubebuilder:validation:MinLength=1
	APIMaticLicensePath *string `json:"licensevolumemountpath,omitempty"`

	// The name of volume from where the APIMatic license file is to be retrieved
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	APIMaticLicenseVolumeName string `json:"licensevolumename"`

	// The volume source from where the APIMatic license file is to be retrieved
	// +kubebuilder:validation:Required
	APIMaticLicenseVolumeSource corev1.VolumeSource `json:"licensevolumesource"`

	// Additional volumes if required in case of sidecar/init containers
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	AdditionalVolumes []corev1.Volume `json:"additionalvolumes,omitempty"`
}

// APIMaticServiceSpec contains configuration for the service that exposes the APIMatic pods
type APIMaticServiceSpec struct {

	// Type string describes ingress methods for a service. Valid values are ClusterIP, NodePort, LoadBalancer, ExternalName, None. Defaults to ClusterIP
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	Type *corev1.ServiceType `json:"servicetype,omitempty"`

	// clusterIP is the IP address of the service and is usually assigned
	// randomly. If an address is specified manually, is in-range (as per
	// system configuration), and is not in use, it will be allocated to the
	// service; otherwise creation of the service will fail. This field may not
	// be changed through updates unless the type field is also being changed
	// to ExternalName (which requires this field to be blank) or the type
	// field is being changed from ExternalName (in which case this field may
	// optionally be specified, as describe above).  Valid values are "None",
	// empty string (""), or a valid IP address. Setting this to "None" makes a
	// "headless service" (no virtual IP), which is useful when direct endpoint
	// connections are preferred and proxying is not required.  Only applies to
	// types ClusterIP, NodePort, and LoadBalancer. If this field is specified
	// when creating a Service of type ExternalName, creation will fail. This
	// field will be wiped when updating a Service to type ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +kubebuilder:validation:Optional
	ClusterIP *string `json:"clusterIP,omitempty"`

	// externalIPs is a list of IP addresses for which nodes in the cluster
	// will also accept traffic for this service.  These IPs are not managed by
	// Kubernetes.  The user is responsible for ensuring that traffic arrives
	// at a node with this IP.  A common example is external load-balancers
	// that are not part of the Kubernetes system.
	// +kubebuilder:validation:Optional
	ExternalIPs []string `json:"externalIPs,omitempty"`

	// Supports "ClientIP" and "None". Used to maintain session affinity.
	// Enable client IP based session affinity.
	// Must be ClientIP or None.
	// Defaults to None.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=ClientIP;None
	SessionAffinity *corev1.ServiceAffinity `json:"sessionAffinity,omitempty"`

	// externalName is the external reference that discovery mechanisms will
	// return as an alias for this service (e.g. a DNS CNAME record). No
	// proxying will be involved.  Must be a lowercase RFC-1123 hostname
	// (https://tools.ietf.org/html/rfc1123) and requires Type to be ExternalName
	// +kubebuilder:validation:Optional
	ExternalName *string `json:"externalName,omitempty"`

	// externalTrafficPolicy denotes if this Service desires to route external
	// traffic to node-local or cluster-wide endpoints. "Local" preserves the
	// client source IP and avoids a second hop for LoadBalancer and Nodeport
	// type services, but risks potentially imbalanced traffic spreading.
	// "Cluster" obscures the client source IP and may cause a second hop to
	// another node, but should have good overall load-spreading.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Local;Cluster
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`

	// healthCheckNodePort specifies the healthcheck nodePort for the service.
	// This only applies when type is set to LoadBalancer and
	// externalTrafficPolicy is set to Local. If a value is specified, is
	// in-range, and is not in use, it will be used.  If not specified, a value
	// will be automatically allocated.  External systems (e.g. load-balancers)
	// can use this port to determine if a given node holds endpoints for this
	// service or not.  If this field is specified when creating a Service
	// which does not need it, creation will fail. This field will be wiped
	// when updating a Service to no longer need it (e.g. changing type).
	// +kubebuilder:validation:Optional
	HealthCheckNodePort *int32 `json:"healthCheckNodePort,omitempty"`

	// publishNotReadyAddresses indicates that any agent which deals with endpoints for this
	// Service should disregard any indications of ready/not-ready.
	// The primary use case for setting this field is for a StatefulSet's Headless Service to
	// propagate SRV DNS records for its Pods for the purpose of peer discovery.
	// The Kubernetes controllers that generate Endpoints and EndpointSlice resources for
	// Services interpret this to mean that all endpoints are considered "ready" even if the
	// Pods themselves are not. Agents which consume only Kubernetes generated endpoints
	// through the Endpoints or EndpointSlice resources can safely assume this behavior. Defaults to false
	// +kubebuilder:validation:Optional
	PublishNotReadyAddresses *bool `json:"publishNotReadyAddresses,omitempty"`

	// sessionAffinityConfig contains the configurations of session affinity.
	// +kubebuilder:validation:Optional
	SessionAffinityConfig *corev1.SessionAffinityConfig `json:"sessionAffinityConfig,omitempty"`

	// topologyKeys is a preference-order list of topology keys which
	// implementations of services should use to preferentially sort endpoints
	// when accessing this Service, it can not be used at the same time as
	// externalTrafficPolicy=Local.
	// Topology keys must be valid label keys and at most 16 keys may be specified.
	// Endpoints are chosen based on the first topology key with available backends.
	// If this field is specified and all entries have no backends that match
	// the topology of the client, the service has no backends for that client
	// and connections should fail.
	// The special value "*" may be used to mean "any topology". This catch-all
	// value, if used, only makes sense as the last value in the list.
	// If this is not specified or empty, no topology constraints will be applied.
	// This field is alpha-level and is only honored by servers that enable the ServiceTopology feature.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	TopologyKeys []string `json:"topologyKeys,omitempty"`

	// IPFamilies is a list of IP families (e.g. IPv4, IPv6) assigned to this
	// service, and is gated by the "IPv6DualStack" feature gate.  This field
	// is usually assigned automatically based on cluster configuration and the
	// ipFamilyPolicy field. If this field is specified manually, the requested
	// family is available in the cluster, and ipFamilyPolicy allows it, it
	// will be used; otherwise creation of the service will fail.  This field
	// is conditionally mutable: it allows for adding or removing a secondary
	// IP family, but it does not allow changing the primary IP family of the
	// Service.  Valid values are "IPv4" and "IPv6".  This field only applies
	// to Services of types ClusterIP, NodePort, and LoadBalancer, and does
	// apply to "headless" services.  This field will be wiped when updating a
	// Service to type ExternalName.
	//
	// This field may hold a maximum of two entries (dual-stack families, in
	// either order).  These families must correspond to the values of the
	// clusterIPs field, if specified. Both clusterIPs and ipFamilies are
	// governed by the ipFamilyPolicy field.
	// +listType=atomic
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=2
	IPFamilies []corev1.IPFamily `json:"ipFamilies,omitempty"`

	// IPFamilyPolicy represents the dual-stack-ness requested or required by
	// this Service, and is gated by the "IPv6DualStack" feature gate.  If
	// there is no value provided, then this field will be set to SingleStack.
	// Services can be "SingleStack" (a single IP family), "PreferDualStack"
	// (two IP families on dual-stack configured clusters or a single IP family
	// on single-stack clusters), or "RequireDualStack" (two IP families on
	// dual-stack configured clusters, otherwise fail). The ipFamilies and
	// clusterIPs fields depend on the value of this field.  This field will be
	// wiped when updating a service to type ExternalName.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=SingeStack;PreferDualStack;RequireDualStack
	IPFamilyPolicy *corev1.IPFamilyPolicyType `json:"ipFamilyPolicy,omitempty"`

	// allocateLoadBalancerNodePorts defines if NodePorts will be automatically
	// allocated for services with type LoadBalancer.  Default is "true". It may be
	// set to "false" if the cluster load-balancer does not rely on NodePorts.
	// allocateLoadBalancerNodePorts may only be set for services with type LoadBalancer
	// and will be cleared if the type is changed to any other type.
	// This field is alpha-level and is only honored by servers that enable the ServiceLBNodePortControl feature.
	// +kubebuilder:validation:Optional
	AllocateLoadBalancerNodePorts *bool `json:"allocateLoadBalancerNodePorts,omitempty"`

	//APIMatic Service Port specifies how the APIMatic service is exposed within the pod
	// +kubebuilder:validation:Required
	APIMaticServicePort *APIMaticServicePort `json:"apimaticserviceport"`

	// Additional volumes if required in case of sidecar/init containers
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	AdditionalServicePorts []corev1.ServicePort `json:"additionalserviceports,omitempty"`

	// Only applies to Service Type: LoadBalancer
	// LoadBalancer will get created with the IP specified in this field.
	// This feature depends on whether the underlying cloud-provider supports specifying
	// the loadBalancerIP when a load balancer is created.
	// This field will be ignored if the cloud-provider does not support the feature.
	// +optional
	LoadBalancerIP *string `json:"loadBalancerIP,omitempty" protobuf:"bytes,8,opt,name=loadBalancerIP"`
}

// APIMaticServicePort configures the APIMatic container ports exposed by the service
type APIMaticServicePort struct {
	// The name of the APIMatic service port within the service. This must be a DNS_LABEL.
	// All ports within a ServiceSpec must have unique names. When considering
	// the endpoints for a Service, this must match the 'name' field in the
	// EndpointPort.
	// Optional if only one ServicePort is defined on this service.
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`

	// The port on each node on which this service is exposed when type is
	// NodePort or LoadBalancer.  Usually assigned by the system. If a value is
	// specified, in-range, and not in use it will be used, otherwise the
	// operation will fail.  If not specified, a port will be allocated if this
	// Service requires one.  If this field is specified when creating a
	// Service which does not need it, creation will fail. This field will be
	// wiped when updating a Service to no longer need it (e.g. changing type
	// from NodePort to ClusterIP).
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=30000
	// +kubebuilder:validation:Maximum=32767
	NodePort *int32 `json:"nodePort,omitempty"`

	// The port that will be exposed by this service.
	// +kubebuilder:validation:Required
	Port int32 `json:"port"`
}

// APIMaticPodPlacementSpec configures the APIMatic pod scheduling policy
// +kubebuilder:validation:MinProperties=1
type APIMaticPodPlacementSpec struct {
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinProperties=1
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// NodeName is a request to schedule this pod onto a specific node. If it is non-empty,
	// the scheduler simply schedules this pod onto that node, assuming that it fits resource
	// requirements.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	NodeName *string `json:"nodeName,omitempty"`

	// If specified, the pod's tolerations.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`

	// Describes node affinity scheduling rules for the pod.
	// +kubebuilder:validation:Optional
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty"`
	// Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)).
	// +kubebuilder:validation:Optional
	PodAffinity *corev1.PodAffinity `json:"podAffinity,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=apm
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.statefulsetStatus.replicas

// APIMatic is the Schema for the apimatics API
type APIMatic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIMaticSpec   `json:"spec,omitempty"`
	Status APIMaticStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// APIMaticList contains a list of APIMatic
type APIMaticList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIMatic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&APIMatic{}, &APIMaticList{})
}
