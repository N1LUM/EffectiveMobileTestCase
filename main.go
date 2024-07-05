package main

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"test/internal/db"
	"test/internal/handlers/crud/tasks"
	"test/internal/handlers/crud/users"
	"test/internal/logging"
)

func main() {
	logging.InitLogger()
	db.ConnectDB()

	router := mux.NewRouter()

	logging.Log.Info("Создан роутинг")

	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/create", users.CreateUser).Methods("POST")
	usersRouter.HandleFunc("/delete/{id}", users.DeleteUserByID).Methods("DELETE")
	usersRouter.HandleFunc("/update/{id}", users.UpdateUserByID).Methods("PUT")
	usersRouter.HandleFunc("/get/{id}", users.GetUserByID).Methods("GET")
	usersRouter.HandleFunc("/list", users.GetUsers).Methods("GET")
	usersRouter.HandleFunc("/laborCost/{user_id}", users.LaborCost).Methods("GET")

	tasksRouter := router.PathPrefix("/tasks").Subrouter()
	tasksRouter.HandleFunc("/create/{user_id}", tasks.CreateTask).Methods("POST")
	tasksRouter.HandleFunc("/start/{id}", tasks.StartTaskTimer).Methods("POST")
	tasksRouter.HandleFunc("/stop/{id}", tasks.StopTaskTimer).Methods("POST")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	http.ListenAndServe("localhost:8080", router)

	logging.Log.Info("Сервис готов. Открыто соединение для прослушивания запросов")
}
