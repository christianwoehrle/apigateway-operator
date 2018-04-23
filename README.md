# apigateway-operator

**!!WIP!!**

kubernetes operator that dynamically adds services to an ingress resource



This project defines a CustomResourceDefinition ApiGateway and an Operator that brings this CRD to life.

The Apigateway Operator creates an Ingress for every apigateway.
The ingress handles traffic for every service, that has a label with then same name "XY" as specified
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

If a service with the label ```labelOfServicesThatHAveTobeAddedToIngress``` is created
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


then this service is adde automatically to the ingress

```
Name:             christians-apigateway2-ingress
Namespace:        default
Address:
Default backend:  lumpensammler:8080 (<none>)
Rules:
  Host  Path  Backends
  ----  ----  --------
  *
        /pathforthiservice   servicetest:10001 (10.244.0.36:6379,10.244.0.37:6379)
Annotations:
Events:  <none>
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

Just to try it: create the first api-gateway, check that it is there and delete it again
```
kubectl apply -f api-gateway.yaml

kubectl get ApiGateway

kubectl describe ApiGateway apigateway

# kubectl delete apigateway apigateway
```

