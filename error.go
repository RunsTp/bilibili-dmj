package bilibili_dmj

import (
	"fmt"
	"net/http"
)

type HttpRequestError struct {
	Msg      string `json:"msg"`
	Response *http.Response
}

func (e *HttpRequestError) Error() string {
	return fmt.Sprintf("msg: %s", e.Msg)
}

func NewHttpRequestError(msg string, resp *http.Response) *HttpRequestError {
	return &HttpRequestError{
		Msg:      msg,
		Response: resp,
	}
}
