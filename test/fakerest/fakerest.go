package fakerest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ReplyFunc defines how a particular matching expected request should be replied
type ReplyFunc func(*http.Request, *http.Response) (*http.Response, error)

// ExpectedRequest declares an expected HTTP request and its corresponding reply
type ExpectedRequest struct {
	Method    string
	URIPrefix string
	Reply     ReplyFunc
}

// ReplyToGet is a shortcut to create an expected GET request reply
func ReplyToGet(prefix string, reply ReplyFunc) ExpectedRequest {
	return newReply(http.MethodGet, prefix, reply)
}

// ReplyToPost is a shortcut to create an expected POST request reply
func ReplyToPost(prefix string, reply ReplyFunc) ExpectedRequest {
	return newReply(http.MethodPost, prefix, reply)
}

// ReplyToPatch is a shortcut to create an expected PATCH request reply
func ReplyToPatch(prefix string, reply ReplyFunc) ExpectedRequest {
	return newReply(http.MethodPatch, prefix, reply)
}

// ReplyToDelete is a shortcut to create an expected PATCH request reply
func ReplyToDelete(prefix string, reply ReplyFunc) ExpectedRequest {
	return newReply(http.MethodDelete, prefix, reply)
}

func newReply(method, prefix string, reply ReplyFunc) ExpectedRequest {
	return ExpectedRequest{Method: method, URIPrefix: prefix, Reply: reply}
}

func (er *ExpectedRequest) Matches(req *http.Request) bool {
	return req.Method == er.Method && strings.HasPrefix(req.URL.Path, er.URIPrefix)
}

// Script is a sequence of expected request to be evaluated in that order
type Script []ExpectedRequest

type Server struct {
	script Script
}

func NewServer(script Script) *Server {
	return &Server{script: script}
}

func NewCombinedServer(servers ...*Server) *Server {
	script := Script{}
	for _, server := range servers {
		script = append(script, server.script...)
	}
	return NewServer(script)
}

func (s *Server) RoundTrip(req *http.Request) (*http.Response, error) {
	rsp := NewResponseTo(req)
	for _, expectedRequest := range s.script {
		if expectedRequest.Matches(req) {
			return expectedRequest.Reply(req, rsp)
		}
	}
	panic(fmt.Errorf("unexpected request (unimplemented): %s %s", req.Method, req.URL.Path))
}

// NewResponseTo returns a base response for the given request
func NewResponseTo(req *http.Request) *http.Response {
	return &http.Response{
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
		Header:     map[string][]string{},
		Request:    req,
	}
}

// SetBody prepares a response body
func SetBody(rsp *http.Response, body string) {
	if rsp.StatusCode == 0 {
		SetStatus(rsp, http.StatusOK)
	}
	rsp.Body = io.NopCloser(bytes.NewBufferString(body))
	rsp.ContentLength = int64(len(body))
}

// SetStatus sets the response status with default status text
func SetStatus(rsp *http.Response, code int) {
	SetStatusCustom(rsp, code, http.StatusText(code))
}

// SetStatusCustom sets the response status with a custom status text
func SetStatusCustom(rsp *http.Response, code int, status string) {
	rsp.StatusCode = code
	rsp.Status = status
}

// NotFoundJSONReply default not found JSON object reply
func NotFoundJSONReply(req *http.Request, rsp *http.Response) (*http.Response, error) {
	return StatusJSON(rsp, http.StatusNotFound, "{}")
}

// RemovedJSONReply default removed JSON reply
func RemovedJSONReply(req *http.Request, rsp *http.Response) (*http.Response, error) {
	return StatusJSON(rsp, http.StatusNoContent, "{}")
}

// EmptyJSONObject returns an OK with an empty JSON object
func EmptyJSONObject(rsp *http.Response) (*http.Response, error) {
	return OKJSON(rsp, "{}")
}

// EmptyJSONArray returns an OK with an empty JSON object
func EmptyJSONArray(rsp *http.Response) (*http.Response, error) {
	return OKJSON(rsp, "[]")
}

// OKJSON replies with a status OK and some JSON content
func OKJSON(rsp *http.Response, json string) (*http.Response, error) {
	return StatusJSON(rsp, http.StatusOK, json)
}

// StatusJSON replies with the given status and some JSON content
func StatusJSON(rsp *http.Response, status int, json string) (*http.Response, error) {
	SetStatus(rsp, status)
	SetBody(rsp, json)
	return rsp, nil
}

// StatusEmpty replies with the given status and no content
func StatusEmpty(rsp *http.Response, status int) (*http.Response, error) {
	SetStatus(rsp, status)
	return rsp, nil
}
