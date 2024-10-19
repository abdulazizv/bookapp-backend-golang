package util

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"gitlab.com/bookapp/api/models"
	"google.golang.org/api/option"
)

func SendNotify(req *models.BookNotify) error {
	conf := &firebase.Config{
		ProjectID: "book-app-7f4e1",
	}
	opt := option.WithCredentialsFile("./config/android-notify-key.json")
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		fmt.Println("error initializing app:  ", err.Error())
	}
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize FCM client: %v\n", err)
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    req.Title,
			Body:     req.Body,
			ImageURL: req.ImageUrl,
		},
		Topic: "/topics/users_book_app",
	}
	response, err := client.Send(ctx, message)
	if err != nil {
		fmt.Println("Failed to send message: ", err.Error())
		return err
	}

	log.Printf("Successfully sent message: %v\n", response)
	return nil
}
