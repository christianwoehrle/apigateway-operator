# apigateway-operator
kubernetes operator that dynamically adds services to an ingress resource

!!WIP

The Apigateway Controller creates an Ingress for every apigateway.
The ingress controller handles traffic for every service, that has a label with then same name "XY" as specified
in the ApiGateways attribute ServiceLabel.


Every Service with this Label is added to the Ingress. The value of the label serves as the path in the ingress.

For an ApiGateway
```
apiVersion: cw.com/v1
kind: ApiGateway
metadata:
  name: christians-apigateway2
spec:
  serviceLabel: labelOfServicesThatHAveTobeAddedToIngress
```

an Ingress is created like that:

```
apiVersion: cw.com/v1
kind: ApiGateway
metadata:
  name: christians-apigateway2
spec:
  serviceLabel: labelOfServicesThatHAveTobeAddedToIngress
  host: cw.de
  backend:
    serviceName: lumpensammler
    servicePort: 8080
```




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


