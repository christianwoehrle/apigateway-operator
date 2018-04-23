# apigateway-operator
kubernetes operator that dynamically adds services to an ingress resource

!!WIP

The Apigateway Controller creates Ingress Resources that handle traffic for services with specific labels.

The ApiGateway specifies a ServiceLabel that has to be set in the srvices, that want to be accessiable via the ApiGateway/Ingress Controller

The operator checks if services are started with the sepcified Label and addes the Service to the Ingress Controller with the path being the value of the service label.




## Custom Resource Defintion
First thing that's needed for a new Operator is a custom resource definition as specified in api-gateway-crd.yaml

This can be added with
```
kubectl apply -f api-gateway-crd.yaml
```

Check the new  crd is installed:
```
kubectl get crd
```


kubectl apply -f api-gateway.yaml

kubectl get ApiGateway

kubectl describe ApiGateway apigateway

kubectl delete apigateway apigateway


