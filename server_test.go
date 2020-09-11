package src

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestNewServerFail(t *testing.T) {
	server := NewServer("4500")
	err := server.Start()
	if err == nil {
		t.Error("Start server should fail")
	}
}

func TestHandler(t *testing.T) {
	const path = "/mock"
	const method = "GET"
	const port = ":4500"
	const content = "server running on :4500"
	const targetMock = "http://example.com/foo"
	handlerMock := func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, content) }
	req := httptest.NewRequest(method, targetMock, nil)
	w := httptest.NewRecorder()
	server := NewServer(port)
	server.Handle(path, method, handlerMock)
	handlerMock(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if server.Router.rules == nil {
		t.Error("Register handler filed")
	}
	if server.Router.rules[path] == nil {
		t.Error("path /mock does not registered")
	}
	if server.Router.rules[path][method] == nil {
		t.Error("method GET does not registered")
	}
	if resp.StatusCode != 200 {
		t.Errorf("got %v, want %v", resp.StatusCode, 200)
	}
	if bodyResponse := string(body); bodyResponse != content {
		t.Errorf("got %v, want %v", bodyResponse, content)
	}
}

func TestDialServers(t *testing.T) {
	_ = os.Setenv(DNSTimeout, "5000000000")
	_ = os.Setenv(ServerPort, ":4500")
	tempServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "I am MR. Content? or Changelog? Maybe Both")
	}))
	defer tempServer.Close()
	u, err := url.Parse(tempServer.URL)
	if err != nil {
		log.Fatal(err)
	}
	dnsList := []string{u.Host, u.Host}
	response := DialServers(dnsList)
	if response != "" {
		t.Errorf("got %v, want %v", response, "")
	}
}
