package app

import (
	"fmt"

	fastRouter "github.com/fasthttp/router"
)

// type RequestMethod int

// const (
// 	GET RequestMethod = iota
// 	POST
// )

func createRouter() *fastRouter.Router {

	fmt.Println("Router new")
	router := fastRouter.New()
	return router
}

func initHandlers() {
	app.router.GET("/start/{botid}", startBotHandler)

	app.router.GET("/health", healthHandler)

	app.router.POST("/new", newBotHandler)
}

// func registerHandler(method RequestMethod, handler fasthttp.RequestHandler) {
// 	switch method {
// 	case GET:
// 		app.router.GET("/health", handler)
// 	case POST:
// 		app.router.POST("/health", handler)
// 	}
// }
