package qingbo

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type QingboClient struct {
	AppKey  string
	AppId   string
	Url     string
	Service string
	Version string
	Signature
}

func NewQingboClient(appKey, appId string) *QingboClient {
	return &QingboClient{
		AppKey:  appKey,
		AppId:   appId,
		Version: "v1",
		Url:     "http://api.gsdata.cn/",
	}
}

func (q *QingboClient) SetService(service string) {
	q.Service = service
}

func (q *QingboClient) SetVersion(version string) {
	q.Version = version
}

func (q *QingboClient) get(uri string, query, service string) (string, error) {
	url := q.Url + service + "/" + q.Version + "/" + uri
	resp, err := q.send("GET", url, map[string]string{"query": query})
	return resp, err
}

func (q *QingboClient) send(method, url string, params map[string]string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		method,
		url+"?"+params["query"],
		strings.NewReader(params["body"]),
	)
	if err != nil {
		panic(err)
	}
	q.Signature.SignRequest(req, q)
	log.Printf("[Time: %s] Required URL: %s\n", time.Now().Format("2006/01/02 15:04:05"), req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}
