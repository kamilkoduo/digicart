package main

import "github.com/kamilkoduo/digicart/src/service"

func main() {
	//rest.Run()
	srv:=service.CartApiServer{}
	_, _ = srv.GetCart("cart1")
}