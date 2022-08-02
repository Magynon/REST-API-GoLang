package main

import "STORE/service"

func main() {
	s := service.NewService()
	s.StartWebService()
}
