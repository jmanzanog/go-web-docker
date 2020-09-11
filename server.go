package src

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type server struct {
	Port   string
	Router *Router
}

//Create new server instance
func NewServer(port string) *server {
	return &server{
		Port:   port,
		Router: NewRouter(),
	}
}

// Start server utility
func (s *server) Start() (err error) {
	http.Handle(RootPath, s.Router)
	err = http.ListenAndServe(s.Port, nil)
	return err
}

// Handle function map endpoints
func (s *server) Handle(path string, method string, handler http.HandlerFunc) {
	if !s.Router.FindPath(path) {
		s.Router.rules[path] = make(map[string]http.HandlerFunc)
	}
	s.Router.rules[path][method] = handler
}

//Add Middleware to http handler
func (s *server) AddMiddleware(handlerFunc http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handlerFunc = middleware(handlerFunc)
	}
	return handlerFunc
}

//Generic Handler type
type Handler func(w http.ResponseWriter, r *http.Request)

//Generic Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

//Write  http response to client

func ServerResponse(w http.ResponseWriter, report string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := fmt.Fprintf(w, report)
	if err != nil {
		log.Panic(err, FailedToResponseServerMsg)
	}

	if len(HealthLogList) > 0 {
		var isSet bool
		lh := &HealthLogList[len(HealthLogList)-1]
		for index, info := range lh.InfoRequests {
			if info.StatusCode == code {
				lh.InfoRequests[index] = InfoRequest{
					StatusCode: code,
					Count:      info.Count + 1,
				}
				isSet = true
			}
		}
		if !isSet {
			lh.InfoRequests = append(lh.InfoRequests, InfoRequest{
				StatusCode: code,
				Count:      1,
			})
		}
	}
}

//Check Dial dns service
func DialServers(dnsList []string) string {
	var s string
	for _, dns := range dnsList {
		timeout := GetTimeout()
		_, errDial := net.DialTimeout(TCPNetwork, dns, time.Duration(timeout))
		if errDial != nil {
			s += StatusFailed
			s += dns
			s += UnableMsg
		}
	}
	return s
}

func GetTimeout() int {
	//number, err := strconv.Atoi(util.GetEnv(DNSTimeout, DNSTimeoutDefaultValue))
	number, err := strconv.Atoi("5000000")
	if err != nil {
		log.Printf("error getting env %s %v", DNSTimeout, err)
	}
	return number
}
