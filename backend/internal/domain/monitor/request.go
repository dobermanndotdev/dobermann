package monitor

import (
	"fmt"
	"net/http"
	"time"
)

type RequestParameters struct {
	method                      RequestMethod
	timeout                     int
	body                        string
	followRedirects             bool
	keepCookiesWhileRedirecting bool
}

func NewRequestParameters(method, body string, timeout int, followRedirects, keepCookiesWhileRedirecting bool) (RequestParameters, error) {
	requestMethod, err := newRequestMethod(method)
	if err != nil {
		return RequestParameters{}, err
	}

	return RequestParameters{
		method:                      requestMethod,
		timeout:                     timeout,
		body:                        body,
		followRedirects:             followRedirects,
		keepCookiesWhileRedirecting: keepCookiesWhileRedirecting,
	}, nil
}

func newDefaultRequestParameters() RequestParameters {
	return RequestParameters{
		method:                      RequestMethodGET,
		timeout:                     int(time.Second * 2),
		body:                        "",
		followRedirects:             false,
		keepCookiesWhileRedirecting: false,
	}
}

func (r RequestParameters) Method() RequestMethod {
	return r.method
}

func (r RequestParameters) Timeout() int {
	return r.timeout
}

func (r RequestParameters) Body() string {
	return r.body
}

func (r RequestParameters) FollowRedirects() bool {
	return r.followRedirects
}

func (r RequestParameters) KeepCookiesWhileRedirecting() bool {
	return r.keepCookiesWhileRedirecting
}

type RequestHeaders map[string]string

type HttpAuth struct {
	username string
	password string
}

func (h HttpAuth) Username() string {
	return h.username
}

func (h HttpAuth) Password() string {
	return h.password
}

var (
	RequestMethodGET     = RequestMethod{value: http.MethodGet}
	RequestMethodPOST    = RequestMethod{value: http.MethodPost}
	RequestMethodDELETE  = RequestMethod{value: http.MethodDelete}
	RequestMethodPATCH   = RequestMethod{value: http.MethodPatch}
	RequestMethodPUT     = RequestMethod{value: http.MethodPut}
	RequestMethodHEAD    = RequestMethod{value: http.MethodHead}
	RequestMethodOPTIONS = RequestMethod{value: http.MethodOptions}
)

type RequestMethod struct {
	value string
}

func (m RequestMethod) String() string {
	return m.value
}

func newRequestMethod(method string) (RequestMethod, error) {
	switch method {
	case RequestMethodGET.value:
		return RequestMethodGET, nil
	case RequestMethodPOST.value:
		return RequestMethodPOST, nil
	case RequestMethodPUT.value:
		return RequestMethodPUT, nil
	case RequestMethodPATCH.value:
		return RequestMethodPATCH, nil
	case RequestMethodDELETE.value:
		return RequestMethodDELETE, nil
	case RequestMethodHEAD.value:
		return RequestMethodHEAD, nil
	case RequestMethodOPTIONS.value:
		return RequestMethodOPTIONS, nil
	default:
		return RequestMethod{}, fmt.Errorf("%s is not a valid request method", method)
	}
}
