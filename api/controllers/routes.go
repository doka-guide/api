// Package controllers - пакет для обработки данных запросов
package controllers

import "github.com/doka-guide/api/api/middlewares"

func (s *Server) initializeRoutes() {

	// Точки входа для сущности Home
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.OptionsHome)).Methods("OPTIONS")
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Точки входа для сущности Login
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.OptionsLogin)).Methods("OPTIONS")
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// Точки входа для сущности User
	s.Router.HandleFunc("/user", middlewares.SetMiddlewareJSON(s.OptionsUsers)).Methods("OPTIONS")
	s.Router.HandleFunc("/user", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/user", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	// Точки входа для сущности Form
	s.Router.HandleFunc("/form", middlewares.SetMiddlewareJSON(s.OptionsForms)).Methods("OPTIONS")
	s.Router.HandleFunc("/form", middlewares.SetMiddlewareJSON(s.CreateForm)).Methods("POST")
	s.Router.HandleFunc("/form", middlewares.SetMiddlewareJSON(s.GetForms)).Methods("GET")
	s.Router.HandleFunc("/form/{id}", middlewares.SetMiddlewareJSON(s.GetForm)).Methods("GET")
	s.Router.HandleFunc("/form/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateForm))).Methods("PUT")
	s.Router.HandleFunc("/form/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteForm)).Methods("DELETE")

	// Точки входа для сущности Subscription
	s.Router.HandleFunc("/subscription", middlewares.SetMiddlewareJSON(s.OptionsSubscriptions)).Methods("OPTIONS")
	s.Router.HandleFunc("/subscription", middlewares.SetMiddlewareJSON(s.CreateSubscription)).Methods("POST")
	s.Router.HandleFunc("/subscription", middlewares.SetMiddlewareJSON(s.GetSubscriptions)).Methods("GET")
	s.Router.HandleFunc("/subscription/{id}", middlewares.SetMiddlewareJSON(s.GetSubscription)).Methods("GET")
	s.Router.HandleFunc("/subscription/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateSubscription))).Methods("PUT")
	s.Router.HandleFunc("/subscription/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteSubscription)).Methods("DELETE")

	// Точки входа для сущности File
	s.Router.HandleFunc("/file", middlewares.SetMiddlewareJSON(s.UploadFile)).Methods("POST")
}
