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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	apicodegenv1beta1 "github.com/apimatic/apimatic-kubernetes-operator/api/v1beta1"
)

var _ = Describe("APIMatic Controller", func() {
	const (
		APIMaticName      = "test-apimatic"
		APIMaticNamespace = "default"
		timeout           = time.Second * 30
		duration          = time.Second * 5
		interval          = time.Millisecond * 250

		ServicePort = 8080
	)

	AfterEach(func() {
		apimaticLookupKey := types.NamespacedName{Name: APIMaticName, Namespace: APIMaticNamespace}
		apimaticCR := &apicodegenv1beta1.APIMatic{}
		err := k8sClient.Get(context.Background(), apimaticLookupKey, apimaticCR)

		if err == nil {
			Expect(k8sClient.Delete(context.Background(), apimaticCR)).Should(Succeed())
		}
	})

	Context("When uploading APIMatic instance to k8s API server with minimal specifications", func() {
		It("Should be created successfully with default values for replicas to 1, terminating grace period to 30s, apimatic license path default to /usr/local/apimatic, default apimatic container name to apimatic", func() {
			By("Creating a new APIMatic instance with minimal specifications")
			ctx := context.Background()

			// Create a docker image secret using your username and password that was registered with APIMatic. More information at https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#create-a-secret-by-providing-credentials-on-the-command-line
			var imagePullSecret corev1.LocalObjectReference = corev1.LocalObjectReference{
				Name: "apimaticimagesecret",
			}
			apimatic := &apicodegenv1beta1.APIMatic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apicodegen.apimatic.io/v1beta1",
					Kind:       "APIMatic",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      APIMaticName,
					Namespace: APIMaticNamespace,
				},
				Spec: apicodegenv1beta1.APIMaticSpec{
					PodSpec: apicodegenv1beta1.APIMaticPodSpec{
						Image:            "apimaticio/apimatic-codegen:latest",
						ImagePullSecrets: []corev1.LocalObjectReference{},
					},
					ServiceSpec: apicodegenv1beta1.APIMaticServiceSpec{
						APIMaticServicePort: &apicodegenv1beta1.APIMaticServicePort{
							Port: ServicePort,
						},
					},
					PodVolumeSpec: apicodegenv1beta1.APIMaticPodVolumeSpec{
						APIMaticLicenseVolumeName: "licensevolume",
						APIMaticLicenseVolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: ConfigMapName,
								},
							},
						},
					},
				},
			}

			apimatic.Spec.PodSpec.ImagePullSecrets = append(apimatic.Spec.PodSpec.ImagePullSecrets, imagePullSecret)

			Expect(k8sClient.Create(ctx, apimatic)).Should(Succeed())

			apimaticLookupKey := types.NamespacedName{Name: APIMaticName, Namespace: APIMaticNamespace}
			createdAPIMaticCR := &apicodegenv1beta1.APIMatic{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, apimaticLookupKey, createdAPIMaticCR)
				return err == nil
			}, timeout, duration).Should(BeTrue())
			Expect(*createdAPIMaticCR.Spec.Replicas).Should(Equal(int32(1)))
			Expect(*createdAPIMaticCR.Spec.PodSpec.TerminationGracePeriodSeconds).Should(Equal(int64(30)))
			Expect(*createdAPIMaticCR.Spec.PodVolumeSpec.APIMaticLicensePath).Should(Equal(string("/usr/local/apimatic")))
			Expect(*createdAPIMaticCR.Spec.PodSpec.Name).Should(Equal(string("apimatic")))

			createdStatefulSet := &appsv1.StatefulSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, apimaticLookupKey, createdStatefulSet)
				return err == nil
			}, timeout, duration).Should(BeTrue())
			Expect(*createdStatefulSet.Spec.Replicas).Should(Equal(int32(1)))

			createdService := &corev1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, apimaticLookupKey, createdService)
				return err == nil
			}, timeout, duration).Should(BeTrue())
			Expect(createdService.Spec.Type).Should(Equal(corev1.ServiceTypeClusterIP))
		})
	})
})
