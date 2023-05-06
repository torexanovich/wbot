package main

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		fmt.Println("BOT_TOKEN not set!")
		return
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}

	restClient := resty.New()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		fmt.Println("Error getting updates:", err)
		return
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		city := update.Message.Text

		resp, err := restClient.R().
			SetQueryParams(map[string]string{
				"q":     city,
				"appid": os.Getenv("WEATHER_API"),
				"units": "metric",
			}).
			SetResult(&weatherResponse{}).
			Get("https://api.openweathermap.org/data/2.5/weather")

		if err != nil {
			fmt.Println("Error getting weather data:", err)
			continue
		}

		if resp.StatusCode() != 200 {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Could not get weather data."))
			continue
		}

		weather := resp.Result().(*weatherResponse)

		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("The temperature in %s is %.1f°C.", weather.Name, weather.Main.Temp)))
	}
}

type weatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float32 `json:"temp"`
	} `json:"main"`
}
