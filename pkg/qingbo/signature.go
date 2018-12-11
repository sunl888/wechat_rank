package qingbo

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type signature struct{}

const (
	VERSION       = "1.0.2"
	ISO8601_BASIC = "20060102T150405Z"
)

func (s *signature) signRequest(request *http.Request, qClient *client) {
	// 必须使用 GMT Time
	hour, _ := time.ParseDuration("-8h")
	now := time.Now().Add(hour).Format(ISO8601_BASIC)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
		return
	}
	h := sha256.New()
	// 这里必须要写入body 不能使用Sum(body) 否则post请求会报错 暂时还没搞懂为啥 下次看到这里再研究(希望不要再看到了)
	h.Write(body)
	sumBody := fmt.Sprintf("%x", h.Sum(nil))
	qClientUrl, err := url.Parse(qClient.Url)
	request.Header.Set("x-gsdata-date", now)
	request.Header.Set("Host", qClientUrl.Host)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "GSDATA-v"+VERSION+"-SDK")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	creq, headers := s.createContext(request, sumBody)
	toSign := s.createStringToSign(now, creq)

	signingKey := s.getSigningKey(now[0:8], request.URL.EscapedPath(), qClient.AppKey)

	mac := hmac.New(sha256.New, signingKey)
	mac.Write([]byte(toSign))
	signature := fmt.Sprintf("%x", mac.Sum(nil))
	request.Header.Set("Authorization", "GSDATA-HMAC-SHA256 AppKey="+qClient.AppId+", "+"SignedHeaders="+headers+", Signature="+signature)
}

func (s *signature) createContext(r *http.Request, sumBody string) (creq, headers string) {
	blacklist := hashset.New()
	blacklist.Add(
		"cache-control",
		"content-type",
		"content-length",
		"expect",
		"max-forwards",
		"pragma",
		"range",
		"te",
		"if-match",
		"if-none-match",
		"if-modified-since",
		"if-unmodified-since",
		"if-range",
		"accept",
		"authorization",
		"proxy-authorization",
		"from",
		"referer",
		"x-gsdagta-trace-id",
	)

	canon := r.Method + "\n" + r.URL.EscapedPath() + "\n" + r.URL.Query().Encode() + "\n"
	aggregat := make(map[string]string, 3)
	for k, v := range r.Header {
		lk := strings.ToLower(k)
		if !blacklist.Contains(lk) {
			aggregat[lk] = string(v[0])
		}
	}
	var (
		signedHeadersString string
		signedString        string
	)
	for k, v := range aggregat {
		signedHeadersString += k + ";"
		signedString += k + ":" + v + "\n"
	}
	signedHeadersString = strings.TrimRight(signedHeadersString, ";")
	canon += signedString + signedHeadersString + "\n" + sumBody
	return canon, signedHeadersString
}

func (s *signature) createStringToSign(lTime, creq string) string {
	h := sha256.New()
	h.Write([]byte(creq))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return "GSDATA-HMAC-SHA256\n" + lTime + "\n" + hash
}

func (s *signature) getSigningKey(shortDate, service, secretKey string) []byte {
	//hmac, use sha1
	key := []byte("GSDATA" + secretKey)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(shortDate))
	dateKey := mac.Sum(nil)

	mac1 := hmac.New(sha256.New, []byte(dateKey))
	mac1.Write([]byte(service))
	serviceKey := mac1.Sum(nil)

	mac2 := hmac.New(sha256.New, []byte(serviceKey))
	mac2.Write([]byte("gsdata_request"))
	finalKey := mac2.Sum(nil)

	return finalKey
}
