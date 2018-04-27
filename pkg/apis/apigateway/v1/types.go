package v1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Backend struct {
	ServiceName string
	ServicePort intstr.IntOrString
}

type ApiGatewaySpec struct {
	ServiceLabel string
	Host         string
	Backend      Backend
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiGateway is the main Type and used to create the apigateway-controller and corresponding ingress
type ApiGateway struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata,omitempty"`
	Spec          ApiGatewaySpec `json:"spec"`
}

type ApiGatewayWatchEvent struct {
	Type   string
	Object ApiGateway
}
