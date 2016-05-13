package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/callback", func(c *gin.Context) {
		proxyURL, _ := url.Parse(os.Getenv("FIXIE_URL"))
		client := &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		}

		bot, err := linebot.NewClient(123456789, "SECRET", "MID", linebot.WithHTTPClient(client))
		if err != nil {
			fmt.Print("bot init error")
			fmt.Println(err)
			return
		}

		received, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				fmt.Print("invalid sign")
				fmt.Println(err)
			}
			fmt.Print("unknown error")
			return
		}

		for _, result := range received.Results {
			content := result.Content()

			if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
				text, err := content.TextContent()
				if err != nil {
					fmt.Print("invalid text content")
					fmt.Println(err)
					return
				}
				fmt.Println(text.Text)
				newStr := ""
				for _, c := range text.Text {
					newStr += string(c)
					newStr += "゛"
				}
				res, err := bot.SendText([]string{content.From}, newStr)
				if err != nil {
					fmt.Println(res)
					fmt.Println(err)
				}
			}

			if content != nil && content.IsOperation && content.OpType == linebot.OpTypeAddedAsFriend {
				op, err := content.OperationContent()
				if err != nil {
					fmt.Println(err)
					return
				}
				from := op.Params[0]

				prof, err := bot.GetUserProfile([]string{from})
				if err != nil {
					fmt.Println(err)
					return
				}
				g := prof.Contacts[0].DisplayName + "、俺゛と゛勝゛負゛し゛ろ゛お゛ぉ゛ぉ゛ぉ゛ぉ゛ぉ゛！゛！゛"
				_, err = bot.SendText([]string{from}, g)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})

	router.Run(":" + port)
}
