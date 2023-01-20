package router

import (
	"github.com/gorilla/mux"
	"crud/middleware"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users/adduser", middleware.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/getusers", middleware.GetAllUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/getuser/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/updateuser/{id}", middleware.UpdateUser).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/users/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}