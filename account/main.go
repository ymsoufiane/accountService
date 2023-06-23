package main

import (
	"account/context"
	"account/routes"
	"net/http"

)

func main() {

	port := context.Config.Server.Port
	context.InfoLogger.Println("server start in port: " + port)

	http.ListenAndServe(":"+port, routes.Router)
	//log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization","Token"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(routes.Router)))

}
