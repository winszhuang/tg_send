package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Config struct {
	Token  string `yaml:"token"`
	XCode  string `yaml:"x_code"`
	ApiUrl string `yaml:"api_url"`
}

func (c *Config) Load(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, c)
}

func XCodeAuthMiddleware(cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		xCodeHeader := c.GetHeader("XCode")
		if xCodeHeader != cfg.XCode {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or missing XCode header",
			})
			c.Abort() // 阻止請求繼續處理
			return
		}
		c.Next() // 如果驗證通過，繼續處理後續的 Handler
	}
}

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

	// 群組 ID 存儲
	groupIDs := make(map[int64]bool)

	// 處理機器人加入新群組的事件
	bot.Handle(tele.OnAddedToGroup, func(c tele.Context) error {
		chatId := c.Chat().ID
		groupIDs[chatId] = true
		fmt.Printf("加入了新群組: %d\n", chatId)
		return nil
	})

	bot.Handle("/hello", func(c tele.Context) error {
		fmt.Println(c.Chat().ID)
		return c.Send("Hello!")
	})

	bot.Group().Handle("/check-group", func(c tele.Context) error {
		groupId := c.Chat().ID
		if _, isExist := groupIDs[groupId]; !isExist {
			groupIDs[groupId] = true
		}
		return nil
	})

	go bot.Start()
	defer bot.Close()

	// Setup Gin router
	router := gin.Default()

	router.Use(XCodeAuthMiddleware(&cfg))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "pong"})
	})

	// POST /push endpoint
	router.POST("/push", func(c *gin.Context) {
		// 讀取純文本請求體
		message, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Bad request"})
			return
		}

		if len(groupIDs) == 0 {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "There is no group to send message to"})
			return
		}

		// 轉換 message 為 string
		messageStr := string(message)
		for chatID := range groupIDs {
			_, err := bot.Send(tele.ChatID(chatID), messageStr)
			if err != nil {
				if strings.Contains(err.Error(), "the group chat was deleted") {
					fmt.Println(err.Error())
					delete(groupIDs, chatID)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message sent"})
	})

	go pingServerPeriodically(cfg, router)

	router.Run()
}

func pingServerPeriodically(cfg Config, router *gin.Engine) {
	minutes := 13
	ticker := time.NewTicker(time.Duration(minutes) * time.Minute)

	for {
		select {
		case <-ticker.C:
			url := cfg.ApiUrl + "/ping"
			_, err := http.Get(url)
			if err != nil {
				fmt.Println("Error pinging web server:", err)
			} else {
				fmt.Println("Web server ping successful")
			}
		}
	}
}
