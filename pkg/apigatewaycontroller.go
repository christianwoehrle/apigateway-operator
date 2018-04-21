package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"log"
	"github.com/christianwoehrle/apigateway-operator/pkg/v1"
	"io/ioutil"
	"strings"
)

func main() {
	log.Println("apiGateway-controller started.------------------------")




	for {
		// Url "cw.com" must match the config spec.group in api-gateway-crd.yaml
		// URL "apigateways" must match the config spec.names.plural in api-gateway-crd.yaml
		resp, err := http.Get("http://localhost:8001/apis/cw.com/v1/apigateways?watch=true")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var event v1.ApiGatewayWatchEvent
			if err := decoder.Decode(&event); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			log.Printf("Received watch event: %s: %s: \n", event.Type, event.Object.Metadata.Name)

			if event.Type == "ADDED" {
				fmt.Println(event.Object.Spec.ServiceLabel)
				fmt.Println("===============================================")

				//createApiGateway(event.Object)
			} else if event.Type == "DELETED" {
				//deleteApiGateway(event.Object)
			}
		}
	}

}


func listenForServiceChange() {
	log.Println("apiGateway-controller started.------------------------")




	for {
		// Url "cw.com" must match the config spec.group in api-gateway-crd.yaml
		// URL "apigateways" must match the config spec.names.plural in api-gateway-crd.yaml
		resp, err := http.Get("http://localhost:8001/api/v1/services?watch=true")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var event v1.ApiGatewayWatchEvent
			if err := decoder.Decode(&event); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			log.Printf("Received watch event: %s: %s: \n", event.Type, event.Object.Metadata.Name)

			if event.Type == "ADDED" {
				fmt.Println(event.Object.Spec.ServiceLabel)
				fmt.Println("===============================================")

				//createApiGateway(event.Object)
			} else if event.Type == "DELETED" {
				//deleteApiGateway(event.Object)
			}
		}
	}

}

func createApiGateway(apiGateway v1.ApiGateway) {
	createResource(apiGateway, "apis/extensions/v1beta1", "ingresses", "../ingress-template.json")
}

func deleteApiGateway(apiGateway v1.ApiGateway) {
	deleteResource(apiGateway, "api/v1", "Ingress", getName(apiGateway));

}

func createResource(apiGateway v1.ApiGateway, apiGroup string, kind string, filename string) {
	log.Printf("Creating %s with name %s in namespace %s with group %s", kind, getName(apiGateway), apiGateway.Metadata.Namespace, apiGroup)
	templateBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	template := strings.Replace(string(templateBytes), "[NAME]", getName(apiGateway), -1)
fmt.Println(template)
	fmt.Printf("http://localhost:8001/%s/namespaces/%s/%s/\n\n", apiGroup, apiGateway.Metadata.Namespace, kind)
	resp, err := http.Post(fmt.Sprintf("http://localhost:8001/%s/namespaces/%s/%s/", apiGroup, apiGateway.Metadata.Namespace, kind), "application/json", strings.NewReader(template))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("response Status:", resp.Status)
}

func deleteResource(apiGateway v1.ApiGateway, apiGroup string, kind string, name string) {
	log.Printf("Deleting %s with name %s in namespace %s", kind, name, apiGateway.Metadata.Namespace)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8001/%s/namespaces/%s/%s/%s", apiGroup, apiGateway.Metadata.Namespace, kind, name), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("response Status:", resp.Status)

}

func getName(apiGateway v1.ApiGateway) string {
	return apiGateway.Metadata.Name + "-apiGateway";
}