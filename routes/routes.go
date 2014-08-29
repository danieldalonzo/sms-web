package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/votinginfoproject/sms-web/queue"
	"github.com/votinginfoproject/sms-web/sms"
	"github.com/votinginfoproject/sms-web/status"
)

type Server struct {
	handler http.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	log.Print(fmt.Sprintf("[INFO] [REQUEST] Method: %s - Path: %s - Host: %s, FormData: %s", r.Method, r.URL.RequestURI(), r.Host, r.Form))

	s.handler.ServeHTTP(w, r)
}

func New(q queue.ExternalQueueService) *Server {
	routes := httprouter.New()

	routes.PanicHandler = func(res http.ResponseWriter, req *http.Request, _ interface{}) {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "text/plain")
		log.Print("[ERROR] : ", req)
	}

	routes.GET("/status", status.Get)

	if q != nil {
		sms.WireUp(q)
	}
	routes.POST("/", sms.Receive)

	return &Server{routes}
}
