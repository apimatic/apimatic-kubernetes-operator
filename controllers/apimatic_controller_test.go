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

	Context("When uploading APIMatic instance to k8s API server with specifications", func() {
		It("Should be created successfully", func() {
			By("Creating a new APIMatic instance with specifications")
			ctx := context.Background()

			// Create a docker image secret using your username and password that was registered with APIMatic. More information at https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#create-a-secret-by-providing-credentials-on-the-command-line
			var imagePullSecret string = "apimaticimagesecret"
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
					Replicas: 1,
					PodSpec: apicodegenv1beta1.APIMaticPodSpec{
						APIMaticContainerSpec: apicodegenv1beta1.APIMaticContainerSpec{
							Image:           "apimaticio/apimatic-codegen:latest",
							ImagePullSecret: &imagePullSecret,
						},
					},
					ServiceSpec: apicodegenv1beta1.APIMaticServiceSpec{
						APIMaticServicePort: &apicodegenv1beta1.APIMaticServicePort{
							Port: ServicePort,
						},
					},
					LicenseSpec: apicodegenv1beta1.APIMaticLicenseSpec{
						APIMaticLicenseSourceType: "ConfigMap",
						APIMaticLicenseSourceName: ConfigMapName,
					},
				},
			}

			Expect(k8sClient.Create(ctx, apimatic)).Should(Succeed())

			apimaticLookupKey := types.NamespacedName{Name: APIMaticName, Namespace: APIMaticNamespace}
			createdAPIMaticCR := &apicodegenv1beta1.APIMatic{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, apimaticLookupKey, createdAPIMaticCR)
				return err == nil
			}, timeout, duration).Should(BeTrue())
			Expect(createdAPIMaticCR.Spec.Replicas).Should(Equal(int32(1)))
			Expect(createdAPIMaticCR.Spec.PodSpec.TerminationGracePeriodSeconds).Should(Equal(int64(30)))

			createdDeployment := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, apimaticLookupKey, createdDeployment)
				return err == nil
			}, timeout, duration).Should(BeTrue())
			Expect(*createdDeployment.Spec.Replicas).Should(Equal(int32(1)))

			createdService := &corev1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, apimaticLookupKey, createdService)
				return err == nil
			}, timeout, duration).Should(BeTrue())
			Expect(createdService.Spec.Type).Should(Equal(corev1.ServiceTypeClusterIP))
		})
	})
})
