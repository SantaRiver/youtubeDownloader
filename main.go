package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kkdai/youtube/v2"
)

func downloadYouTubeVideo(url string) (string, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		return "", err
	}

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		return "", err
	}

	file, err := os.Create(video.Title + ".mp4")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.ReadFrom(stream)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func main() {
	botToken := "7404857795:AAGIYbdVnBv0PD4qGV0KJu6NSklroYsZq5s"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me a YouTube video link and I will download it for you.")
			bot.Send(msg)
		} else if update.Message.Text != "" {
			videoURL := update.Message.Text

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Downloading video...")
			bot.Send(msg)

			videoPath, err := downloadYouTubeVideo(videoURL)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error downloading video: %v", err))
				bot.Send(msg)
				continue
			}

			videoFile := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath(videoPath))
			bot.Send(videoFile)

			os.Remove(videoPath) // Clean up the downloaded file
		}
	}
}
