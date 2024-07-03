package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"test/internal/db"
	"test/internal/handlers/crud"
	"test/internal/logging"
)

func main() {
	logging.InitLogger()
	db.ConnectDB()

	router := mux.NewRouter()

	logging.Log.Info("Создан роутинг")

	router.HandleFunc("/createUser", crud.CreateUser).Methods("POST")
	router.HandleFunc("/deleteUser/{id}", crud.DeleteUserByID).Methods("DELETE")
	http.ListenAndServe("localhost:8080", router)

	logging.Log.Info("Сервис готов. Открыто соединение для прослушивания запросов")
}
