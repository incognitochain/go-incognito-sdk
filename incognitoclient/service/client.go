package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type IncogClientInterface interface {
	Post() error
	buildParameter() error
}

type IncogClient struct {
	Client        *http.Client
	ChainEndpoint string
}

func (i *IncogClient) PostAndReceiveInterface(method string, params interface{}) (interface{}, []byte, error) {
	body, err := i.Post(method, params)
	if err != nil {
		return nil, nil, errors.Wrap(err, "post")
	}

	var v interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, nil, errors.Wrap(err, "json.Unmarshal")
	}

	return v, body, nil
}

func (i *IncogClient) Post(method string, params interface{}) ([]byte, error) {
	args := i.buildParameter(method, params)

	data, err := json.Marshal(args)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}

	//fmt.Println("chain" , b.config.Incognito.ChainEndpoint)
	if strings.Contains(i.ChainEndpoint, "https") {
		//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		i.Client.Transport = tr
	}

	req, err := http.NewRequest(http.MethodPost, i.ChainEndpoint, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Close = true

	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", i.BearerToken)

	// if b.config.Env == "localhost" {
	// 	db, _ := httputil.DumpRequest(req, true)
	// 	b.logger.With(zap.String("request", string(db))).Debug("post")
	// 	fmt.Printf("string(db) = %+v\n", string(db))
	// 	fmt.Printf("Call %+v\n", args)
	// }

	resp, err := i.Client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "b.c.Do: %q", req.URL.String())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// if b.config.Env == "localhost" {
	// 	fmt.Println("________________________________________________")
	// 	fmt.Println("Response from BlockChain:", string(body))
	// 	fmt.Println("________________________________________________")
	// }
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	return body, nil
}

func (i *IncogClient) buildParameter(method string, params interface{}) *Parameter {
	payload := &Parameter{
		JsonRpc: "1.0",
		Method:  method,
		Params:  params,
		Id:      1,
	}

	return payload
}
