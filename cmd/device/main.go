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
	ID string
	//conn *net.UDPConn
}

func NewDevice(id string) (*Device, error) {
	//conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 10001, Zone: ""})
	//if err != nil {
	//return &Device{}, err
	//}
	d := &Device{
		ID: id,
		//conn: conn,
	}

	d.Send("detach")
	d.Send("attach")
	d.Send("subscribe")
	return d, nil
}

func (d Device) Send(action string) {
	//d.conn.Write([]byte(fmt.Sprintf("%s %s", d.ID, action)))
}

func (d Device) Close() {
	//d.conn.Close()
}

func main() {
	flag.Parse()

	d, err := NewDevice("device0")
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating the device: %v", err))
	}
	d.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Info("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Error(err)
				return
			}

			// TODO: process receive message for a device
			log.Info(fmt.Sprintf("Message: %s", message))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// TODO: process sending messages for a device
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Error(err)
				return
			}
		case <-interrupt:
			log.Info("Interrupt sent")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
