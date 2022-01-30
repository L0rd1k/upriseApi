package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var ErrJsonPayloadEmpty = errors.New("JSON payload is empty")

type Request struct {
	*http.Request
	PathParams  map[string]string
	environment map[string]interface{}
}

func (rqst *Request) PathParam(name string) string {
	fmt.Println("<Request:PathParam>:", name)
	return rqst.PathParams[name]
}

func (rqst *Request) DecodeJsonPayload(v interface{}) error {
	fmt.Println("<Request:DecodeJsonPayload>")
	content, err := ioutil.ReadAll(rqst.Body)
	rqst.Body.Close()
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return ErrJsonPayloadEmpty
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return err
	}
	return nil
}

func (rqst *Request) BaseUrl() *url.URL {
	fmt.Println("<Request:BaseUrl>")
	scheme := rqst.URL.Scheme //> Define current Scheme
	fmt.Println("	<Request:BaseUrl>:scheme", scheme)
	if scheme == "" {
		scheme = "http"
	} else if scheme == "http" && rqst.TLS != nil {
		scheme = "https"
	}
	host := rqst.Host
	fmt.Println("	<Request:BaseUrl>:host", host)

	if len(host) > 0 && host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}
	fmt.Println("	<Request:BaseUrl>:host", host)

	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}

func (rqst *Request) UrlFor(path string, queryParams map[string][]string) *url.URL {
	fmt.Println("<Request:UrlFor>")
	baseUrl := rqst.BaseUrl()
	baseUrl.Path = path
	if queryParams != nil {
		query := url.Values{}
		for k, v := range queryParams {
			for _, vv := range v {
				query.Add(k, vv)
			}
		}
		baseUrl.RawQuery = query.Encode()
	}
	return baseUrl
}

// ==============================================================================

type CorsInfo struct {
	IsCors                      bool
	IsPreflight                 bool
	Origin                      string
	OriginUrl                   *url.URL
	AccessControlRequestMethod  string
	AccessControlRequestHeaders []string
}

func (rqst *Request) GetCorsInfo() *CorsInfo {
	fmt.Println("<Request:GetCorsInfo>")
	origin := rqst.Header.Get("Origin")
	var originUrl *url.URL
	var isCors bool
	if origin == "" {
		isCors = false
	} else if origin == "null" {
		isCors = true
	} else {
		var err error
		originUrl, err = url.ParseRequestURI(origin)
		isCors = err == nil && rqst.Host != originUrl.Host
	}
	reqMethod := rqst.Header.Get("Access-Control-Request-Method")
	reqHeaders := []string{}
	rawReqHeaders := rqst.Header[http.CanonicalHeaderKey("Access-Control-Request-Headers")]
	for _, rawReqHeader := range rawReqHeaders {
		if len(rawReqHeader) == 0 {
			continue
		}
		for _, reqHeader := range strings.Split(rawReqHeader, ",") {
			reqHeaders = append(reqHeaders, http.CanonicalHeaderKey(strings.TrimSpace(reqHeader)))
		}
	}
	isPreflight := isCors && rqst.Method == "OPTIONS" && reqMethod != ""
	return &CorsInfo{
		IsCors:                      isCors,
		IsPreflight:                 isPreflight,
		Origin:                      origin,
		OriginUrl:                   originUrl,
		AccessControlRequestMethod:  strings.ToUpper(reqMethod),
		AccessControlRequestHeaders: reqHeaders,
	}
}
