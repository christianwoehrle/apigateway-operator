package v1

type Metadata struct {
	Name string
	Namespace string
}

type ApiGatewaySpec struct {
	ServiceLabel string
}

type ApiGateway struct {
	Metadata Metadata
	Spec ApiGatewaySpec
}

type ApiGatewayWatchEvent struct {
	Type string
	Object ApiGateway
}