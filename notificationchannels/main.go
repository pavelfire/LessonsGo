package main

import "fmt"

type Notifier interface {
	Send() string
	ChannelName() string
}

type EmailNotification struct {
	To      string
	Subject string
}

func (e EmailNotification) ChannelName() string {
	return "email"
}

func (e EmailNotification) Send() string {
	return fmt.Sprintf("письмо отправлено на %s: %s", e.To, e.Subject)
}

type PushNotification struct {
	DeviceID string
	Message  string
}

func (p PushNotification) ChannelName() string {
	return "push"
}

func (p PushNotification) Send() string {
	return fmt.Sprintf("push отправлен на устройство %s: %s", p.DeviceID, p.Message)
}

func main() {
	notifiers := []Notifier{
		EmailNotification{To: "user@example.com", Subject: "Новый сезон доступен"},
		PushNotification{DeviceID: "device-42", Message: "Ваш сериал уже онлайн"},
	}

	for _, n := range notifiers {
		fmt.Printf("Канал: %s\nРезультат: %s\n\n", n.ChannelName(), n.Send())
	}
}
