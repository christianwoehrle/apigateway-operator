package v1

import "k8s.io/apimachinery/pkg/util/intstr"

type Backend struct {
	ServiceName string
	ServicePort intstr.IntOrString
}

type Metadata struct {
	Name      string
	Namespace string
}

type ApiGatewaySpec struct {
	ServiceLabel string
	Backend      Backend
}

type ApiGateway struct {
	Metadata Metadata
	Spec     ApiGatewaySpec
}

type ApiGatewayWatchEvent struct {
	Type   string
	Object ApiGateway
}
