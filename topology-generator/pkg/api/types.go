package api

const (
	MODE_SERVER = "s"
	MODE_LOCAL  = "l"
)

var GlobalConfig Configurations

type Endpoint struct {
	Name        string       `json:"name"`
	Path        string       `json:"path"`
	Method      string       `json:"method"`
	Connections []Connection `json:"connections"`
}

type Connection struct {
	Service string `json:"service"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
	Method  string `json:"method"`
}

type Service struct {
	Name      string
	Port      string
	Version   string
	Endpoints []Endpoint
	C         Configurations
}

type Configurations struct {
	Name                    string // "mimik"
	IstioProxyRequestCPU    string
	IstioProxyRequestMemory string
	MimikRequestCPU         string
	MimikRequestMemory      string
	MimikLimitCPU           string
	MimikLimitMemory        string
	EnableInjection         string //"true"
	ImageTag                string //"quay.io/leandroberetta/mimik"
	ImageVersion            string //"v0.0.2"
	InjectionLabel          string
	Replicas                int // 1
}

func NewConfigurations(name string, istioProxyRequestCPU string, istioProxyRequestMemory string, mimikRequestCPU string,
	mimikRequestMemory string, mimikLimitCPU string, mimikLimitMemory string, version string, image string, injection string, replicas int, injectionLabel string) Configurations {

	return Configurations{
		Name:                    name,
		IstioProxyRequestCPU:    istioProxyRequestCPU,
		IstioProxyRequestMemory: istioProxyRequestMemory,
		MimikRequestCPU:         mimikRequestCPU,
		MimikRequestMemory:      mimikRequestMemory,
		MimikLimitCPU:           mimikLimitCPU,
		MimikLimitMemory:        mimikLimitMemory,
		EnableInjection:         injection,
		ImageTag:                image,
		ImageVersion:            version,
		Replicas:                replicas,
		InjectionLabel:          injectionLabel,
	}
}

func NewDefaultConfigurations() Configurations {
	return Configurations{
		Name:                    "mimik",
		IstioProxyRequestCPU:    "50m",
		IstioProxyRequestMemory: "128Mi",
		MimikRequestCPU:         "25m",
		MimikRequestMemory:      "64Mi",
		MimikLimitCPU:           "200m",
		MimikLimitMemory:        "256Mi",
		EnableInjection:         "true",
		ImageTag:                "quay.io/leandroberetta/mimik",
		ImageVersion:            "v0.0.2",
		Replicas:                1,
		InjectionLabel:          "istio-injection:enabled",
	}
}

type Response struct {
	Name              string     `json:"name"`
	Version           string     `json:"version"`
	Path              string     `json:"path"`
	StatusCode        int        `json:"statusCode"`
	UpstreamResponses []Response `json:"upstreamResponses"`
}

type Generator struct {
	Namespaces        int `json:"namespaces"`
	Services          int `json:"services"`
	Connections       int `json:"connections"`
	RandomConnections int `json:"randomConnections"`
}
