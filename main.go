package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"encoding/json"
	"io"
	"log"
	"net/http"

	cwv1 "github.com/christianwoehrle/apigateway-operator/pkg/v1"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ApiGateways struct {
	apigateways map[string]*ApiGateway
}

type ApiGateway struct {
	apigatewayCRD cwv1.ApiGateway
	services      map[string]Service
	ingress       v1beta1.Ingress
}

type Service struct {
	path        string
	serviceName string
	servicePort int
}

var apigateways = ApiGateways{
	apigateways: make(map[string]*ApiGateway),
}

func (*ApiGateways) AddGateway(gw cwv1.ApiGateway, clientset *kubernetes.Clientset) {
	fmt.Println("*ApiGateways --> AddGateway: ", gw)

	apiGateway := ApiGateway{
		apigatewayCRD: gw,
	}

	apigateways.apigateways[gw.Metadata.Name] = &apiGateway

	createIngress(gw, clientset)

}
func (*ApiGateways) DeleteGateway(gw cwv1.ApiGateway, clientset *kubernetes.Clientset) {
	fmt.Println("*ApiGateways --> DeleteGateway, needs implementation: ", gw)

}
func (*ApiGateways) ModifyGateway(gw cwv1.ApiGateway, clientset *kubernetes.Clientset) {
	fmt.Println("*ApiGateways --> ModifyGateway, needs implementation: ", gw)

}

/** AddNewService adds the service to the ingress controllers, that want to handle request for that service
  i.e. these ingresses, for which the ApiGateway has a matching serviceLabel
*/
func (*ApiGateways) AddNewService(service *v1.Service) {
	fmt.Println("Service added (mock): ", service.Name)

}

/** DeleteService deletes the service to the ingress controllers, that want to handle request for that service
  i.e. these ingresses, for which the ApiGateway has a matching serviceLabel
*/
func (*ApiGateways) DeleteService(service *v1.Service) {
	fmt.Println("Service deleted (mock): ", service.Name)

}

/** createIngress adds the Ingress for an ApiGateway
 */
func createIngress(gw cwv1.ApiGateway, clientset *kubernetes.Clientset) {

	ingressName := gw.Metadata.Name + "-ingress"

	// Check if INgress already exists
	ingress, err := clientset.ExtensionsV1beta1().Ingresses("default").Get(ingressName, metav1.GetOptions{})
	if ingress != nil {
		fmt.Println("Ingress already exists, do nothing now (Implementation needed)")
		return
	}

	fmt.Println("add Ingress: ", ingressName)
	ingresses, err := clientset.ExtensionsV1beta1().Ingresses("default").List(metav1.ListOptions{})
	handleErr("Couldn't read ingresses", err)

	for i, ingress := range ingresses.Items {
		fmt.Printf("Ingress %d: %s\n", i, ingress.Name)
	}

	newIngress := v1beta1.Ingress{Spec: v1beta1.IngressSpec{}}
	newIngress.SetName(ingressName)

	//TODO: Add Services
	newIngress.Spec.Backend = &v1beta1.IngressBackend{ServiceName: gw.Spec.Backend.ServiceName, ServicePort: gw.Spec.Backend.ServicePort}
	createdIngress, err := clientset.ExtensionsV1beta1().Ingresses("default").Create(&newIngress)

	if err != nil {
		handleErr("Ingress not created:"+createdIngress.String(), err)
	}
	fmt.Println("Ingress added")
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

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	go handleServiceEvents(clientset)
	handleApiGatewayEvents(clientset)
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
				if event.Type == "ADDED" {
					apigateways.AddNewService(service)
				}
				if event.Type == "DELETED" {
					apigateways.DeleteService(service)
				}
			}

		}
	}
}

func handleApiGatewayEvents(clientset *kubernetes.Clientset) {

	resp, err := http.Get("http://localhost:8001/apis/cw.com/v1/apigateways?watch=true")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	handleApiGatewayEvent(resp, clientset)
}

func handleApiGatewayEvent(resp *http.Response, clientset *kubernetes.Clientset) {
	decoder := json.NewDecoder(resp.Body)
	for {
		var event cwv1.ApiGatewayWatchEvent
		if err := decoder.Decode(&event); err == io.EOF {
			fmt.Println("handleNewApiGateways EOF")
			continue
		} else if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received watch event: %s: %s: \n", event.Type, event.Object.Metadata.Name)

		switch event.Type {
		case "ADDED":
			apigateways.AddGateway(event.Object, clientset)
		case "DELETED":
			apigateways.DeleteGateway(event.Object, clientset)
		case "MODIFIED":
			apigateways.ModifyGateway(event.Object, clientset)
		}
	}
}

func handleErr(text string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", text, err)
		os.Exit(1)
	}
}
