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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	apicodegenv1beta1 "github.com/apimatic/apimatic-kubernetes-operator/api/v1beta1"
)

var _ = Describe("APIMatic Controller", func() {
	const (
		APIMaticName      = "test-apimatic"
		APIMaticNamespace = "default"
		timeout           = time.Second * 5
		duration          = time.Second * 5
		interval          = time.Millisecond * 250

		ConfigMapName      = "test-apimatic-license-configmap"
		ConfigMapNamespace = "default"

		ServicePort = 8080
	)

	Context("When uploading APIMatic instance to k8s API server with minimal specifications", func() {
		It("Should be created successfully with default values for replicas to 1, terminating grace period to 30s, apimatic license path default to /usr/local/apimatic, default apimatic container name to apimatic", func() {
			By("Creating a new APIMatic instance with minimal specifications")
			ctx := context.Background()
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
						Image: "obaidkhattak/apimatic-codegen",
					},
					ServiceSpec: apicodegenv1beta1.APIMaticServiceSpec{
						APIMaticServicePort: &apicodegenv1beta1.APIMaticServicePort{
							Port: 8080,
						},
					},
					PodVolumeSpec: apicodegenv1beta1.APIMaticPodVolumeSpec{
						APIMaticLicenseVolumeName: "licensevolume",
						APIMaticLicenseVolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "apimaticlicense",
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, apimatic)).Should(Succeed())

			apimaticLookupKey := types.NamespacedName{Name: APIMaticName, Namespace: APIMaticNamespace}
			createdAPIMaticCR := &apicodegenv1beta1.APIMatic{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, apimaticLookupKey, createdAPIMaticCR)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdAPIMaticCR.Spec.Replicas).Should(Equal(1))
			Expect(createdAPIMaticCR.Spec.PodSpec.TerminationGracePeriodSeconds).Should(Equal(30))
			Expect(createdAPIMaticCR.Spec.PodVolumeSpec.APIMaticLicensePath).Should(Equal("/usr/local/apimatic"))
			Expect(createdAPIMaticCR.Spec.PodSpec.Name).Should(Equal("apimatic"))
		})
	})
})
