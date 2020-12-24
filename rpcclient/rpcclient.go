package rpcclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	*http.Client
	url      string
	protocol string
	host     string
	port     uint
}

// NewHttpClient to get http client instance
func NewHttpClient(url string, protocol string, host string, port uint) *HttpClient {
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}
	return &HttpClient{
		Client:   httpClient,
		url:      url,
		protocol: protocol,
		host:     host,
		port:     port,
	}
}

func buildHttpServerAddress(url string, protocol string, host string, port uint) string {
	if url != "" {
		return url
	}
	return fmt.Sprintf("%s://%s:%d", protocol, host, port)
}

func (client *HttpClient) RPCCall(
	method string,
	params interface{},
	rpcResponse interface{},
) (err error) {
	rpcEndpoint := buildHttpServerAddress(
		client.url, client.protocol, client.host, client.port,
	)

	payload := map[string]interface{}{
		"method": method,
		"params": params,
		"id":     0,
	}
	payloadInBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(rpcEndpoint, "application/json", bytes.NewBuffer(payloadInBytes))

	if err != nil {
		return err
	}

	respBody := resp.Body
	defer respBody.Close()

	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, rpcResponse)
	if err != nil {
		return err
	}
	return nil
}
