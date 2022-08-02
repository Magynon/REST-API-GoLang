package main

import "workshop/service"

func main() {
	s := service.NewService()
	s.StartWebService()
}
