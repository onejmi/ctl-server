package main

import (
	"fmt"
	"net/http"

	"./data"
	"./routes"
	"./routes/util"
)

func main() {

	setupDatabase()
	configureRoutes()

	data.SetupCronJobs()

	fmt.Println("Starting server...")
	http.ListenAndServe(":8081", nil)
}

func configureRoutes() {
	http.Handle("/register", util.OnlyMethod("POST", util.JSONResponse(handler(routes.Register))))
	http.Handle("/login", util.OnlyMethod("POST", util.JSONResponse(handler(routes.Login))))
	http.Handle("/claim", util.OnlyMethod("POST",
		util.JSONResponse(util.IsAuthenticated(handler(routes.Claim)))))
	http.Handle("/profile",
		util.OnlyMethod("GET", util.JSONResponse(util.IsAuthenticated(handler(routes.GetProfile)))))
}

func setupDatabase() {
	data.Connect()
}

func handler(endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(endpoint)
}
