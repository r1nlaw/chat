package database

import (
	"chater/models"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

var db *sqlx.DB

func InitDB() (*sqlx.DB, error) {

	if err := godotenv.Load(); err != nil {
		log.Fatal("ошибка загрузки .env файла")
	}

	conf := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL"),
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", conf.Host, conf.Port, conf.Username, conf.DBName, conf.Password, conf.SSLMode)

	var err error
	db, err = sqlx.Connect("postgres", connStr) // Присваиваем глобальную переменную db
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %w", err)
	}

	return db, nil
}

func SaveMessage(message models.Message) error {
	if db == nil {
		return fmt.Errorf("database is not initialized")
	}

	// Сохраняем сообщение в базе данных
	_, err := db.Exec("INSERT INTO messages (sender_id, receiver_id, message) VALUES ($1, $2, $3)", message.SenderID, message.ReceiverID, message.Message)
	return err
}

func GetMessageHistory(userID, receiverID string) ([]models.Message, error) {
	// Получаем историю сообщений
	rows, err := db.Query("SELECT sender_id, receiver_id, message FROM messages WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)", userID, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.SenderID, &message.ReceiverID, &message.Message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
