package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	addr = flag.String("addr", "localhost:8080", "http service address")
)

type Device struct {
	ID   string
	send chan string
	conn *websocket.Conn
}

func NewDevice(id string, u url.URL) (*Device, error) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &Device{}, err
	}
	log.Info(fmt.Sprintf("Connected: %s", u.String()))

	d := &Device{
		ID:   id,
		send: make(chan string, 10),
		conn: c,
	}

	d.Send("register", "")
	return d, nil
}

func (d Device) Close() {
	d.conn.Close()
}

func (d Device) Send(action string, data string) {
	message := fmt.Sprintf("%s %s %s", d.ID, action, data)
	d.send <- message
}

func (d Device) Run() {
	log.Info(fmt.Sprintf("Running: %s", d.ID))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// ticker
	go func() {
		for {
			select {
			case t := <-ticker.C:
				d.Send("event", t.String())
			}
		}
	}()

	// read
	go func() {
		defer close(done)
		for {
			_, message, err := d.conn.ReadMessage()
			if err != nil {
				log.Error(err)
				return
			}

			log.Info(fmt.Sprintf("Message: %s", message))
		}
	}()

	// write
	for {
		select {
		case <-done:
			return
		case t := <-d.send:
			err := d.conn.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil {
				log.Error(err)
				return
			}
		case <-interrupt:
			log.Info("Interrupt sent")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := d.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func main() {
	log.SetReportCaller(true)
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

	d, err := NewDevice("device0", u)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating the device: %v", err))
	}
	d.Run()
}
