# apigateway-operator
kubernetes operator that dynamically adds services to an ingress resource

!!WIP

The Apigateway Controller creates Ingress Resources that handle traffic for services with specific labels.

The operator checks if services are started, stopped and adds these services to the ingress resources, that this operater is responsible for.


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
