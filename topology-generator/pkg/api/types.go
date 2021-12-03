package api

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
