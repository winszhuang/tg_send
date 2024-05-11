package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Config struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chat_id"`
}

func (c *Config) Load(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, c)
}

// chatId 1474301143

func main() {
	// 載入設定檔案
	var cfg Config
	err := cfg.Load("config.yml")
	if err != nil {
		log.Fatalf("無法讀取設定檔案: %v", err)
	}

	// Setup Telegram bot
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/hello", func(c tele.Context) error {
		fmt.Println(c.Chat().ID)
		return c.Send("Hello!")
	})

	go bot.Start()

	defer bot.Close()

	// Setup Gin router
	router := gin.Default()

	// POST /push endpoint
	router.POST("/push", func(c *gin.Context) {
		// 讀取純文本請求體
		message, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Bad request"})
			return
		}

		// 轉換 message 為 string
		messageStr := string(message)

		// Send message to Telegram
		//_, err := bot.Send(, json.Message) // Specify the ChatID correctly
		_, err = bot.Send(tele.ChatID(cfg.ChatID), messageStr) // Specify the ChatID correctly
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message sent"})
	})

	router.Run()
}
