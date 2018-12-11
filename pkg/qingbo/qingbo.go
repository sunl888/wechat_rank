package qingbo

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type client struct {
	AppKey  string
	AppId   string
	Url     string
	Service string
	Version string
	signature
}

func NewQingboClient(appKey, appId string) *client {
	return &client{
		AppKey:  appKey,
		AppId:   appId,
		Version: "v1",
		Url:     "http://api.gsdata.cn/",
	}
}

func (q *client) SetService(service string) {
	q.Service = service
}

func (q *client) SetVersion(version string) {
	q.Version = version
}

func (q *client) get(uri string, query map[string]string, service string) (string, error) {
	url := q.Url + service + "/" + q.Version + "/" + uri
	resp, err := q.send("GET", url, query)
	return resp, err
}

func (q *client) post(uri string, query map[string]string, service string) (string, error) {
	url := q.Url + service + "/" + q.Version + "/" + uri
	resp, err := q.send("POST", url, query)
	return resp, err
}

func (q *client) send(method, url string, params map[string]string) (string, error) {
	var (
		err    error
		client *http.Client
		req    *http.Request
	)
	client = &http.Client{}
	var r http.Request
	_ = r.ParseForm()
	for k, v := range params {
		r.Form.Add(k, v)
	}
	switch method {
	case http.MethodPost:
		req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(r.Form.Encode()))
		if err != nil {
			return "", err
		}
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, url+"?"+r.Form.Encode(), strings.NewReader(""))
		if err != nil {
			return "", err
		}
	default:
		return "", errors.QingboError("方法不允许", "方法不允许", 401, 401)
	}
	q.signature.SignRequest(req, q)
	log.Printf("[Time: %s] Method:%s,URL: %s\n", time.Now().Format("2006/01/02 15:04:05"), req.Method, req.URL)
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
