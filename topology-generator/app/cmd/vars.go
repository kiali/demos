package cmd

import "github.com/kiali/demos/topology-generator/pkg/api"

const (
	docPath        = "./doc/commands/"
	ReleaseVersion = "v0.0.1"
)

var generatorConfig api.Generator

var path string

var name, istioProxyRequestCPU, istioProxyRequestMemory, mimikRequestCPU,
	mimikRequestMemory, mimikLimitCPU, mimikLimitMemory, image, enableInjection, injectionlabel, version string

var namespaces, services, connections, randomconnections int

var replicas, serverPort int
