package main

import (
	"chater/database"
	"chater/handlers"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("ошибка инициализации БД: ", err)
	}
	defer db.Close()

	registerRepo := handlers.NewRegisterRepository(db)
	chatRepo := handlers.NewChatRepository(db)
	handlerWithCors := c.Handler(http.DefaultServeMux)

	http.HandleFunc("/API/chater/signIn", registerRepo.SignIn)
	http.HandleFunc("/API/chater/signUp", registerRepo.SignUp)
	http.HandleFunc("/API/chater/chat", chatRepo.SendMessage)

	http.ListenAndServe(":8080", handlerWithCors)
}
