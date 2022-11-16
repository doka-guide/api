// Package controllers - пакет для обработки данных запросов
package controllers

import "github.com/doka-guide/api/api/middlewares"

func (server *Server) initializeRoutes() {

	// Точки входа для сущности Home
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.OptionsHome)).Methods("OPTIONS")
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.Home)).Methods("GET")

	// Точки входа для сущности Login
	server.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(server.OptionsLogin)).Methods("OPTIONS")
	server.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")

	// Точки входа для сущности User
	server.Router.HandleFunc("/user", middlewares.SetMiddlewareJSON(server.OptionsUsers)).Methods("OPTIONS")
	server.Router.HandleFunc("/user", middlewares.SetMiddlewareJSON(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/user", middlewares.SetMiddlewareJSON(server.GetUsers)).Methods("GET")
	server.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateUser))).Methods("PUT")
	server.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteUser)).Methods("DELETE")

	// Точки входа для сущности Form
	server.Router.HandleFunc("/form", middlewares.SetMiddlewareJSON(server.OptionsForms)).Methods("OPTIONS")
	server.Router.HandleFunc("/form", middlewares.SetMiddlewareJSON(server.CreateForm)).Methods("POST")
	server.Router.HandleFunc("/form", middlewares.SetMiddlewareJSON(server.GetForms)).Methods("GET")
	server.Router.HandleFunc("/form/{id}", middlewares.SetMiddlewareJSON(server.GetForm)).Methods("GET")
	server.Router.HandleFunc("/form/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateForm))).Methods("PUT")
	server.Router.HandleFunc("/form/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteForm)).Methods("DELETE")

	// Точки входа для сущности Subscription
	server.Router.HandleFunc("/subscription", middlewares.SetMiddlewareJSON(server.OptionsSubscriptions)).Methods("OPTIONS")
	server.Router.HandleFunc("/subscription", middlewares.SetMiddlewareJSON(server.CreateSubscription)).Methods("POST")
	server.Router.HandleFunc("/subscription", middlewares.SetMiddlewareJSON(server.GetSubscriptions)).Methods("GET")
	server.Router.HandleFunc("/subscription/{id}", middlewares.SetMiddlewareJSON(server.GetSubscription)).Methods("GET")
	server.Router.HandleFunc("/subscription/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateSubscription))).Methods("PUT")
	server.Router.HandleFunc("/subscription/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteSubscription)).Methods("DELETE")

	// Точки входа для сущности ProfileLink
	server.Router.HandleFunc("/profile-link", middlewares.SetMiddlewareJSON(server.OptionsProfileLinks)).Methods("OPTIONS")
	server.Router.HandleFunc("/profile-link", middlewares.SetMiddlewareJSON(server.CreateProfileLink)).Methods("POST")
	server.Router.HandleFunc("/profile-link", middlewares.SetMiddlewareJSON(server.GetProfileLinks)).Methods("GET")
	server.Router.HandleFunc("/profile-link/{id}", middlewares.SetMiddlewareJSON(server.GetProfileLink)).Methods("GET")
	server.Router.HandleFunc("/profile-link/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteProfileLink)).Methods("DELETE")

	// Точки входа для сущности SubscriptionReport
	server.Router.HandleFunc("/subscription-report", middlewares.SetMiddlewareJSON(server.OptionsProfileLinks)).Methods("OPTIONS")
	server.Router.HandleFunc("/subscription-report", middlewares.SetMiddlewareJSON(server.CreateProfileLink)).Methods("POST")
	server.Router.HandleFunc("/subscription-report", middlewares.SetMiddlewareJSON(server.GetProfileLinks)).Methods("GET")
	server.Router.HandleFunc("/subscription-report/{id}", middlewares.SetMiddlewareJSON(server.GetProfileLink)).Methods("GET")
	server.Router.HandleFunc("/subscription-report/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteProfileLink)).Methods("DELETE")

	// Точки входа для сущности File
	server.Router.HandleFunc("/file", middlewares.SetMiddlewareJSON(server.UploadFile)).Methods("POST")
}
