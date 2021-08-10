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
	"path/filepath"
	"testing"

	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	apicodegenv1beta1 "github.com/apimatic/apimatic-kubernetes-operator/api/v1beta1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var testEnv *envtest.Environment

const (
	ConfigMapName      = "test-apimatic-license-configmap"
	ConfigMapNamespace = "default"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = apicodegenv1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&APIMaticReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sClient.Scheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName,
			Namespace: ConfigMapNamespace,
		},
		Data: map[string]string{
			"License.lic": "\u003cLicense\u003e\r\n  \u003cId\u003e2c5a837e-ca01-447a-a393-d509899c7662\u003c/Id\u003e\r\n  \u003cType\u003eStandard\u003c/Type\u003e\r\n  \u003cExpiration\u003eMon, 04 Aug 2025 00:00:00 GMT\u003c/Expiration\u003e\r\n  \u003cQuantity\u003e2\u003c/Quantity\u003e\r\n  \u003cProductFeatures\u003e\r\n    \u003cFeature name=\"importFormats\"\u003eAPIBluePrint,APIElements,Postman2,Wsdl,WADL2009,WADL2006,Raml,Raml10,Swagger20,SwaggerYaml,Swagger10,OpenApi3Json,OpenApi3Yaml,IODocs,IODocsV0314,Insomnia,InsomniaYaml,RAMLAsJson,GoogleDiscovery,Postman\u003c/Feature\u003e\r\n    \u003cFeature name=\"exportFormats\"\u003eSwagger20,SwaggerYaml,Swagger10,OpenApi3Json,OpenApi3Yaml,Wsdl,APIBluePrint,APIElements,Postman2,WADL2009,WADL2006,WSDL,IODocs,IODocsV0314,RAML,RAML10,GoogleDiscovery,Postman20,Insomnia,InsomniaYaml\u003c/Feature\u003e\r\n    \u003cFeature name=\"platforms\"\u003edocs,java,dotnet,ios,node,ruby,python,go,android,php\u003c/Feature\u003e\r\n  \u003c/ProductFeatures\u003e\r\n  \u003cCustomer\u003e\r\n    \u003cName\u003eDeveloper\u003c/Name\u003e\r\n    \u003cEmail\u003edeveloper@apimatic.io\u003c/Email\u003e\r\n  \u003c/Customer\u003e\r\n  \u003cSignature\u003eMEUCIQDrxYeTDrouY4FWz0mgHu+Um6d5UzDNBkwpnZr3MHDkBQIgCgodCwLQVWQC7T7DVRAWWMSMfH6LU3t0XERBQv7FyfI=\u003c/Signature\u003e\r\n\u003c/License\u003e",
		},
	}

	Expect(k8sClient.Create(context.Background(), configMap)).Should(Succeed())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
