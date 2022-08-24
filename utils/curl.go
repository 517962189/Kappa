package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Curl struct {
	Headers map[string]string
}

//headers Content-Length | Content-Type | Host | Cookies ...
func (c *Curl) SetHeader(headers, value string) {
	c.Headers = map[string]string{
		headers: value,
	}
}

func (c *Curl) Post(url string, params []byte) ([]byte, error) {
	reqBody := strings.NewReader(string(params))
	httpReq, err := http.NewRequest("post", url, reqBody)
	if err != nil {
		log.Printf("NewRequest fail, url: %s, reqBody: %s, err: %v\n", url, reqBody, err)
		return nil, err
	}

	if len(c.Headers) > 0 {
		for headerKey, headerVal := range c.Headers {
			httpReq.Header.Add(headerKey, headerVal)
		}
	} else {
		httpReq.Header.Add("Content-Type", "application/json")
	}

	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("do http fail, url: %s, reqBody: %s, err:%v\n", url, reqBody, err)
		return nil, err
	}
	defer httpRsp.Body.Close()

	// Read: HTTP结果
	rspBody, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		log.Printf("ReadAll failed, url: %s, reqBody: %s, err: %v\n", url, reqBody, err)
		return nil, err
	}
	return rspBody, nil
}

func (c *Curl) Get(url string,params []byte) {

}
