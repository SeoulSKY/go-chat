package routes

import (
	"encoding/json"
	"log"
	"main/models"
	"main/mongo"
	"net/http"
)

const path string = "/chat"

func init() {
	RegisterRoute(&Route{
		path:    path,
		handler: get,
		method:  "GET",
	})

	RegisterRoute(&Route{
		path:    path,
		handler: post,
		method:  "POST",
	})
}

func get(w http.ResponseWriter, r *http.Request) {
	chats := mongo.GetInstance().Find()
	if err := json.NewEncoder(w).Encode(chats); err != nil {
		log.Panic(err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	chat := new(models.Chat)
	json.NewDecoder(r.Body).Decode(chat)

	mongo.GetInstance().Insert(chat)
}
