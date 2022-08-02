package service

import (
	"log"
	"net/http"

	"STORE/gateways"

	"github.com/emicklei/go-restful/v3"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) StartWebService() {
	ws := new(restful.WebService)
	restful.Add(ws)

	api := gateways.NewAPI()
	api.RegisterRoutes(ws)

	log.Printf("Started serving on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
