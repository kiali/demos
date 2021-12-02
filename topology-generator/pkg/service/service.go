package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/leandroberetta/mimik/pkg/api"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func NewService(name, port, fileName, version string) (api.Service, error) {
	service := api.Service{Name: name, Port: port, Version: version}
	err := loadEndpoints(fileName, &service.Endpoints)
	return service, err
}

func loadEndpoints(fileName string, endpoints *[]api.Endpoint) error {
	file, err := os.Open(fileName)
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, endpoints)
	return err
}

func EndpointHandler(service api.Service, client httpClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := api.Response{Name: service.Name, Version: service.Version, StatusCode: http.StatusNotFound}
		ch := make(chan api.Response)
		headers := getHeaders(r.Header)
		for _, endpoint := range service.Endpoints {
			resp.Path = r.URL.Path
			if endpoint.Path == r.URL.Path {
				resp.StatusCode = http.StatusOK
				upstreamResponses := make([]api.Response, len(endpoint.Connections))
				for _, connection := range endpoint.Connections {
					go handleReq(makeURL(connection), connection, headers, client, ch)
				}
				for i := range endpoint.Connections {
					upstreamResponses[i] = <-ch
				}
				resp.UpstreamResponses = upstreamResponses
			}
		}
		responseJSON, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

func getHeaders(header http.Header) map[string]string {
	headers := make(map[string]string)
	headers["x-request-id"] = header.Get("x-request-id")
	headers["x-b3-traceid"] = header.Get("x-b3-traceid")
	headers["x-b3-spanid"] = header.Get("x-b3-spanid")
	headers["x-b3-parentspanid"] = header.Get("x-b3-parentspanid")
	headers["x-b3-sampled"] = header.Get("x-b3-sampled")
	headers["x-b3-flags"] = header.Get("x-b3-flags")
	headers["Authorization"] = header.Get("Authorization")
	return headers
}

func handleReq(url string, conn api.Connection, headers map[string]string, client httpClient, ch chan api.Response) {
	req, _ := http.NewRequest(conn.Method, url, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		ch <- makeErrorResponse(conn, http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- makeErrorResponse(conn, http.StatusInternalServerError)
		return
	}
	res := api.Response{}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		ch <- makeErrorResponse(conn, http.StatusInternalServerError)
		return
	}
	ch <- res
}

func makeURL(conn api.Connection) string {
	return fmt.Sprintf("http://%s:%d/%s", conn.Service, conn.Port, conn.Path)
}

func makeErrorResponse(conn api.Connection, error int) api.Response {
	return api.Response{Name: conn.Service, Version: "", StatusCode: error, Path: conn.Path, UpstreamResponses: nil}
}
