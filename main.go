package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/christianwoehrle/apigateway-operator/pkg/apis/apigateway/v1alpha"
	cwclientset "github.com/christianwoehrle/apigateway-operator/pkg/generated/clientset/versioned/typed/apigateway/v1alpha"

	"time"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ApiGateways struct {
	apigateways map[string]*ApiGateway
}

type ApiGateway struct {
	apigatewayCRD v1alpha.ApiGateway
	services      map[string]Service
	ingress       *v1beta1.Ingress
}

type Service struct {
	path        string
	serviceName string
	servicePort int
}

var apigateways = ApiGateways{
	apigateways: make(map[string]*ApiGateway),
}

func (*ApiGateways) AddGateway(gw *v1alpha.ApiGateway, clientset *kubernetes.Clientset) {
	fmt.Println("*ApiGateways --> AddGateway: ", gw)

	apiGateway := ApiGateway{
		apigatewayCRD: *gw,
	}

	apigateways.apigateways[gw.Name] = &apiGateway

	ingress := createIngress(*gw, clientset)
	if ingress != nil {
		apiGateway.ingress = ingress
	}
}
func (*ApiGateways) DeleteGateway(gw *v1alpha.ApiGateway, clientset *kubernetes.Clientset) {
	fmt.Println("*ApiGateways --> DeleteGateway, needs implementation: ", gw)

}
func (*ApiGateways) ModifyGateway(gw *v1alpha.ApiGateway, clientset *kubernetes.Clientset) {
	fmt.Println("*ApiGateways --> ModifyGateway, needs implementation: ", gw)

}

/** AddNewService adds the service to the ingress controllers, that want to handle request for that service
  i.e. these ingresses, for which the ApiGateway has a matching serviceLabel
*/
func (*ApiGateways) AddNewService(service *v1.Service, clientset *kubernetes.Clientset) {
	defer fmt.Println("Ended AddNewService")

	fmt.Println("Service added (mock): ", service.Name)

	for _, apigateway := range apigateways.apigateways {
		servicelabel := apigateway.apigatewayCRD.Spec.ServiceLabel
		for key, value := range service.Labels {
			if key == servicelabel {
				// add new service to ingress with path value
				ingress := apigateway.ingress.DeepCopy()
				servicepath := "/" + value
				if ingress == nil {

					fmt.Println("ingress nicht vorhanden, nichts anlegen")
					return
				}
				backend := v1beta1.IngressBackend{ServiceName: service.Name, ServicePort: intstr.IntOrString{IntVal: service.Spec.Ports[0].Port}}
				path := v1beta1.HTTPIngressPath{Backend: backend, Path: servicepath}
				paths := []v1beta1.HTTPIngressPath{path}

				httpIngressRuleValue := v1beta1.HTTPIngressRuleValue{Paths: paths}

				ingressRuleValue := v1beta1.IngressRuleValue{}
				ingressRuleValue.HTTP = &httpIngressRuleValue
				ingressRule := v1beta1.IngressRule{}
				ingressRule.IngressRuleValue = ingressRuleValue

				ingressRules := []v1beta1.IngressRule{ingressRule}
				ingress.Spec.Rules = ingressRules
				updatedIngress, err := clientset.ExtensionsV1beta1().Ingresses("default").Update(ingress)
				fmt.Println("Updated ingress")
				handleErr("Ingress couldn´t be updates with new Service", err)
				fmt.Println(updatedIngress)

			}
		}
	}

	fmt.Printf("Labels %V \n", service.Labels)
	for key, value := range service.Labels {
		fmt.Printf("Key, Value: %s %s\n", key, value)
	}

}

/** DeleteService deletes the service to the ingress controllers, that want to handle request for that service
  i.e. these ingresses, for which the ApiGateway has a matching serviceLabel
*/
func (*ApiGateways) DeleteService(service *v1.Service, clientset *kubernetes.Clientset) {
	fmt.Println("Service deleted (mock): ", service.Name)

}

func (*ApiGateways) ModifyService(service *v1.Service, clientset *kubernetes.Clientset) {
	fmt.Println("Service deleted (mock): ", service.Name)

}

/** createIngress adds the Ingress for an ApiGateway
 */
func createIngress(gw v1alpha.ApiGateway, clientset *kubernetes.Clientset) *v1beta1.Ingress {

	ingressName := gw.Name + "-ingress"

	// Check if INgress already exists
	ingress, err := clientset.ExtensionsV1beta1().Ingresses("default").Get(ingressName, metav1.GetOptions{})

	if err == nil {
		fmt.Println("Ingress already exists, do nothing now (Implementation needed)", err)
		return ingress
	}

	fmt.Println("add Ingress: ", ingressName)

	newIngress := v1beta1.Ingress{Spec: v1beta1.IngressSpec{}}
	newIngress.SetName(ingressName)

	//TODO: Add Services
	newIngress.Spec.Backend = &v1beta1.IngressBackend{ServiceName: gw.Spec.Backend.ServiceName, ServicePort: gw.Spec.Backend.ServicePort}

	createdIngress, err := clientset.ExtensionsV1beta1().Ingresses("default").Create(&newIngress)

	if err != nil {
		handleErr("Ingress not created:"+createdIngress.String(), err)
	}
	fmt.Println("Ingress added")
	return createdIngress
}

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	apigatewayClientset, err := cwclientset.NewForConfig(config)

	//newApiGateway := v1alpha.ApiGateway{}
	//newApiGateway.Name = "chrissisapigateway"
	//newApiGateway.Spec.ServiceLabel = "chrississervicelabel"
	//createdapiGateway, err := apigatewayClientset.ApiGateways("default").Create(&newApiGateway)
	//fmt.Println("createdapiGateway: ", createdapiGateway)
	//handleErr("error creating apigateway", err)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	go handleApiGatewayEvents(clientset, apigatewayClientset)
	time.Sleep(30 * time.Second)

	go handleServiceEvents(clientset)

	time.Sleep(3600 * time.Second)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func handleServiceEvents(clientset *kubernetes.Clientset) {
	for {
		serviceStreamWatcher, err := clientset.CoreV1().Services("").Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		//fmt.Printf("%T\n", serviceStreamWatcher)
		for {
			select {
			case event := <-serviceStreamWatcher.ResultChan():
				service := event.Object.(*v1.Service)
				fmt.Printf("Labels %V \n", service.Labels)
				for key, value := range service.Labels {
					fmt.Printf("Key, Value: %s %s\n", key, value)
				}

				switch event.Type {
				case "ADDED":
					apigateways.AddNewService(service, clientset)
				case "DELETED":
					apigateways.DeleteService(service, clientset)
				case "MODIFIED":
					apigateways.ModifyService(service, clientset)
				default:
					fmt.Println("UNEXPECTED event.Type in handleServiceEvents")
				}
			}

		}
	}
}

func handleApiGatewayEvents(clientset *kubernetes.Clientset, apigatewayClientSet *cwclientset.CwV1alphaClient) {
	for {
		apigatewayWatcher, err := apigatewayClientSet.ApiGateways("default").Watch(metav1.ListOptions{})
		if err != nil {

			fmt.Println("%V", err)
			fmt.Printf("%T\n", err)
			fmt.Printf("%V", err.(*errors.StatusError).ErrStatus)
			panic(err.Error())
		}
		//fmt.Printf("%T\n", serviceStreamWatcher)
		for {
			select {
			case event := <-apigatewayWatcher.ResultChan():
				apigateway := event.Object.(*v1alpha.ApiGateway)

				switch event.Type {
				case "ADDED":
					apigateways.AddGateway(apigateway, clientset)
				case "DELETED":
					apigateways.DeleteGateway(apigateway, clientset)
				case "MODIFIED":
					apigateways.ModifyGateway(apigateway, clientset)
				default:
					fmt.Println("UNEXPECTED event.Type in handleServiceEvents")
				}
			}

		}
	}

}

func handleErr(text string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", text, err)
		os.Exit(1)
	}
}
