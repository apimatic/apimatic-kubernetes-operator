## Table of contents

* [Introduction](#introduction)
* [Features](#features)
* [Running the Sample](#running-the-sample)
- * [Prerequisites](#prerequisites)
- * [Steps](#steps)
* [Technical Support](#technical-support)
* [Copyrights](#copyrights)

## Introduction

APIMatic Operator simplifies the configuration and lifecycle management of the APIMatic code and docs generation on different Kubernetes distributions and OpenShift. The Operator encapsulates key operational knowledge on how to configure and upgrade the APIMatic CodeGen application, making the use of it for APIMatic API management features easy to set and use.


More information about the underlying APIMatic CodeGen API that is exposed
by this operator can be found [here](https://apimatic-core-v3-docs.netlify.app/#/http/getting-started/overview-apimatic-core)

## Features

APIMatic Operator provides the following features:
- Deploys the APIMatic CodeGen Web API service within the Kubernetes or OpenShift cluster.
- Expose the APIMatic CodeGen API external to the cluster using Service type as [NodePort](https://kubernetes.io/docs/concepts/services-networking/service/#nodeport),[LoadBalancer](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) or [ExternalName](https://kubernetes.io/docs/concepts/services-networking/service/#externalname).
- For exposing the service through an ingress resource, create an ingress resource in the namespace of your APIMatic CR and set owned APIMatic service created by the operator as a backed service. More information can be found [here](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- Manual horizontal scaling of pods
  ```sh
  kubectl scale apm apimatic-sample--replicas=2
  ```

## Running the Sample

### Prerequisites

Please contact APIMatic at [support@apimatic.io](mailto:support@apimatic.io) to register with the APIMatic CodeGen dockerhub registry and acquire a valid license to run the APIMatic CodeGen API.

Further prerequisites for running the sample include:

- [go](https://golang.org/) v1.16.*
- [git](https://git-scm.com/)
- [make](https://www.gnu.org/software/make/)
- A running Kubernetes cluster with [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) on client. For testing purposes, you can use [Minikube](https://minikube.sigs.k8s.io/docs/) or [kind](https://kind.sigs.k8s.io/)
- For checking the service created by the APIMatic operator on-prem, you can use [MetalLB](https://metallb.org/)

### Steps

Following are the steps to run the sample for checking the APIMatic operator.

- Clone the APIMatic repository into your working directory
  ```sh  
  git clone https://github.com/apimatic/apimatic-kubernetes-operator.git  
  ```
- Run 'make deploy' to set up the APIMatic operator resources. This will deploy the 'apimatic-system' namespace as well as the CRD and the RBAC manifests.

- Create a secret named 'apimaticimagesecret' to allow pulling the APIMatic CodeGen image using the dockerhub username registered with APIMatic. If you have not done so, please contact APIMatic.io at support@apimatic.io for the steps required.
  ```sh
  kubectl create secret docker-registry apimaticimagesecret --docker-server=https://index.docker.io/v1/ --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
  ```
- Create a [configmap](https://kubernetes.io/docs/concepts/configuration/configmap/) resource named 'apimaticlicense' that will serve as the volume storing the APIMatic license. An example of this is given below which will create the configmap using the License.lic file located at /usr/local/apimatic/license/  
  ```sh  
  kubectl create configmap apimaticlicense --from-file /usr/local/apimatic/license/License.lic
  ```
- This will deploy a ConfigMap resource with the following definition
  ```sh
  apiVersion: v1  
  data:
    License.lic: \"<License file contents here>\"  
  kind: ConfigMap
  metadata:    
    name: apimaticlicense    
    namespace: default  
  ```
- Now run the sample using the following command
  ```sh  
     kubectl apply -f config/samples/apicodegen_v1beta1_apimatic.yaml
  ```
- You will now see a new [StatefulSet](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) with replica count of 3 and [Service](https://kubernetes.io/docs/concepts/services-networking/service/) created, both named apimatic-sample. The service is of type 'NodePort' with nodeport set to 32000. Accessing http://localhost:32000 from your browser(or using curl from within the Minikube or Kind cluster, if using that) should now show the APIMatic Titan page.
- You can now use the exposed APIMatic CodeGen API to generate API SDKs and Docs using [curl](https://curl.se/), [Postman](https://www.postman.com/) or your own custom Web application that consumes the APIMatic CodeGen API service.

## Technical Support

- To request additional features in the future, or if you notice any discrepency regarding this document, please drop an email to [support@apimatic.io](mailto:support@apimatic.io)

### Copyrights

&copy; 2021 APIMatic.io