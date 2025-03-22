package handlers

import (
	"chater/database"
	"chater/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

type ChatSession struct {
	conn       *websocket.Conn
	userID     string
	receiverID string
}

type ChatRepository struct {
	Db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{Db: db}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *ChatRepository) SendMessage(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrading connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("user_ID")
	receiverID := r.URL.Query().Get("receiver_id")

	chatSession := newChatSession(conn, userID, receiverID)
	chatSession.Start()
}

func newChatSession(conn *websocket.Conn, userID, receiverID string) *ChatSession {
	return &ChatSession{
		conn:       conn,
		userID:     userID,
		receiverID: receiverID,
	}
}

func (cs *ChatSession) Start() {
	for {
		_, msg, err := cs.conn.ReadMessage()
		if err != nil {
			http.Error(nil, "error reading message: ", http.StatusInternalServerError)
			return
		}

		message := models.Message{}
		if err := json.Unmarshal(msg, &message); err != nil {
			http.Error(nil, "failed format data", http.StatusBadRequest)
			continue
		}
		err = database.SaveMessage(message)
		if err != nil {
			log.Println("error sending message: ", err)
			continue
		}

		err = cs.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("error sending message: ", err)
			break
		}

	}

}

func (cs *ChatSession) GetHistory() ([]models.Message, error) {
	return database.GetMessageHistory(cs.userID, cs.receiverID)
}
