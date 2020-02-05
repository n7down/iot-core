package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/n7down/iot-core/internal/gateway"
	log "github.com/sirupsen/logrus"
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	gatewayID = flag.String("id", "gateway0", "gateway id")
)

const (
	//projectID = "iota-3345"
	projectID      = "noke-smart-entry-dev"
	registryID     = "demo-registry"
	cloudRegion    = "us-central1"
	privateKeyFile = "rsa_key.pem"
	//privateKeyFile = "rsa_private.pem"
	//publicKeyFile      = "rsa_cert.pem"
	algorithm          = "RS256"
	mqttBridgeHostname = "tls://mqtt.googleapis.com"
	mqttBridgePort     = "443"
	jwtExpiresMinutes  = 1200
	protocolVersion    = 4
)

const (
	DETACH_ACTION      = "detach"
	ATTACH_ACTION      = "attach"
	SUBSCRIBE_ACTION   = "subscribe"
	EVENT_ACTION       = "event"
	UNSUBSCRIBE_ACTION = "unsubscribe"
	REGISTER_ACTION    = "register"
	UNREGISTER_ACTION  = "unregister"
)

func main() {

	flag.Parse()

	log.SetReportCaller(true)

	command := make(chan string, 1000)

	hub := gateway.NewHub(command)
	go hub.Run()

	// onConnect defines the on connect handler which resets backoff variables.
	var onConnect mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Info(fmt.Sprintf("Client connected %s:%s: %t\n", mqttBridgeHostname, mqttBridgePort, client.IsConnected()))
	}

	// onMessage defines the message handler for the mqtt client.
	var onMessage mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())
		log.Info(fmt.Sprintf("Topic: %s Message: %s\n", topic, payload))

		// FIXME: handle the errors from iot core

		// FIXME: use grouping from regex to get the topicDeviceID
		s := strings.Split(topic, "/")
		topicDeviceID := s[2]

		// to send to a device on a gateway - send 'deviceID message' on the gateway commands
		// to send to the device just send the command on deviceID commands
		switch topicDeviceID {
		case *gatewayID:
			m := strings.Split(payload, " ")
			deviceID := m[0]
			message := strings.Join(m[1:], " ")
			hub.Send(deviceID, message)
		default:
			hub.Send(topicDeviceID, payload)
		}
	}

	// onDisconnect defines the connection lost handler for the mqtt client.
	var onDisconnect mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Info("Client disconnected")
	}

	signBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Declare the token with the algorithm used for signing, and the claims
	t := jwt.New(jwt.GetSigningMethod(algorithm))
	t.Claims = jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(jwtExpiresMinutes * time.Minute).Unix(),
		Audience:  projectID,
	}

	jwt, err := t.SignedString(signKey)
	if err != nil {
		log.Fatal(err)
	}

	clientID := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, cloudRegion, registryID, *gatewayID)
	log.Info(fmt.Sprintf("Connecting to: %s", clientID))

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%s", mqttBridgeHostname, mqttBridgePort))
	opts.SetClientID(clientID)
	opts.SetUsername("unused")
	opts.SetPassword(jwt)
	opts.SetProtocolVersion(protocolVersion)
	opts.SetOnConnectHandler(onConnect)
	opts.SetDefaultPublishHandler(onMessage)
	opts.SetConnectionLostHandler(onDisconnect)

	// Create and connect a client using the above options.
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(fmt.Sprintf("Failed to connect client: %s", token.Error()))
	}

	time.Sleep(3 * time.Second)

	gatewayTopic := fmt.Sprintf("/devices/%s/commands/#", *gatewayID)
	if token := client.Subscribe(gatewayTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
	}
	log.Info(fmt.Sprintf("Connected to topic: %s", gatewayTopic))

	errorsTopic := fmt.Sprintf("/devices/%s/errors", *gatewayID)
	if token := client.Subscribe(errorsTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
	}
	log.Info(fmt.Sprintf("Connected to topic: %s", errorsTopic))

	go func(c mqtt.Client) {
		for {
			select {
			case cmd := <-command:

				log.Info(fmt.Sprintf("Received command: %s", cmd))
				command := strings.Fields(cmd)
				id := command[0]
				action := command[1]

				switch action {

				case DETACH_ACTION:

					// detach
					type token struct {
						Authorization string `json:"authorization"`
					}

					t := &token{
						Authorization: jwt,
					}

					payload, err := json.Marshal(t)
					if err != nil {
						log.Error(err)
					}

					detachTopic := fmt.Sprintf("/devices/%s/detach", id)
					log.Info(fmt.Sprintf("Detach: %s topic: %s", id, detachTopic))
					if token := c.Publish(detachTopic, 1, false, payload); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					} else {
						log.Info(fmt.Sprintf("Published to topic: %s", detachTopic))
					}

				case ATTACH_ACTION:

					type token struct {
						Authorization string `json:"authorization"`
					}

					t := &token{
						Authorization: jwt,
					}

					payload, err := json.Marshal(t)
					if err != nil {
						log.Error(err)
					}

					// attach
					attachTopic := fmt.Sprintf("/devices/%s/attach", id)
					log.Info(fmt.Sprintf("Attach: %s topic: %s", id, attachTopic))
					if token := c.Publish(attachTopic, 1, false, payload); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					} else {
						log.Info(fmt.Sprintf("Published to topic: %s", attachTopic))
					}

				case SUBSCRIBE_ACTION:

					// subscribe
					subscribeTopic := fmt.Sprintf("/devices/%s/commands/#", id)
					log.Info(fmt.Sprintf("Subscribe: %s topic: %s", id, subscribeTopic))
					if token := c.Subscribe(subscribeTopic, 0, nil); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					} else {
						log.Info(fmt.Sprintf("Connected to topic: %s", subscribeTopic))
					}

				case UNSUBSCRIBE_ACTION:

					// unsubscribe
					unsubscribeTopic := fmt.Sprintf("/devices/%s/commands/#", id)
					log.Info(fmt.Sprintf("Unsubscribe: %s topic: %s", id, unsubscribeTopic))
					if token := client.Unsubscribe(gatewayTopic, unsubscribeTopic); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					} else {
						log.Info(fmt.Sprintf("Unsubscribed completed: %s topic: %s", id, unsubscribeTopic))
					}

				case REGISTER_ACTION:

					// attach
					type token struct {
						Authorization string `json:"authorization"`
					}

					t := &token{
						Authorization: jwt,
					}

					payload, err := json.Marshal(t)
					if err != nil {
						log.Error(err)
					}

					attachTopic := fmt.Sprintf("/devices/%s/attach", id)
					log.Info(fmt.Sprintf("Attach: %s topic: %s", id, attachTopic))
					if token := c.Publish(attachTopic, 1, false, payload); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					} else {
						log.Info(fmt.Sprintf("Published to topic: %s", attachTopic))
					}

					time.Sleep(5 * time.Second)

					// subscribe
					subscribeTopic := fmt.Sprintf("/devices/%s/commands/#", id)
					log.Info(fmt.Sprintf("Subscribe: %s topic: %s", id, subscribeTopic))
					if token := c.Subscribe(subscribeTopic, 0, nil); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					} else {
						log.Info(fmt.Sprintf("Connected to topic: %s", subscribeTopic))
					}

				case UNREGISTER_ACTION:

					// detach
					type token struct {
						Authorization string `json:"authorization"`
					}

					t := &token{
						Authorization: jwt,
					}

					payload, err := json.Marshal(t)
					if err != nil {
						log.Error(err)
					}

					detachTopic := fmt.Sprintf("/devices/%s/detach", id)
					log.Info(fmt.Sprintf("Detach: %s topic: %s", id, detachTopic))
					if token := c.Publish(detachTopic, 1, false, payload); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					} else {
						log.Info(fmt.Sprintf("Published to topic: %s", detachTopic))
					}

					// unsubscribe
					unsubscribeTopic := fmt.Sprintf("/devices/%s/commands/#", id)
					log.Info(fmt.Sprintf("Unsubscribe: %s topic: %s", id, unsubscribeTopic))
					if token := client.Unsubscribe(gatewayTopic, unsubscribeTopic); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					} else {
						log.Info(fmt.Sprintf("Unsubscribed completed: %s topic: %s", id, unsubscribeTopic))
					}

				case EVENT_ACTION:

					// event
					data := strings.Join(command[2:], " ")
					//attachTopic := fmt.Sprintf("/devices/%s/events/register", id)
					//if token := c.Publish(attachTopic, 1, false, data); token.Wait() && token.Error() != nil {
					//log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					//}
					//log.Info(fmt.Sprintf("Sent event action from %s: %s", id, data))

					eventTopic := fmt.Sprintf("/devices/%s/events", id)
					if token := c.Publish(eventTopic, 1, false, data); token.Wait() && token.Error() != nil {
						log.Error(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
					}
					log.Info(fmt.Sprintf("Sent event action from %s: %s", *gatewayID, data))

				default:
					log.Info(fmt.Sprintf("Unknown action: %s", command[1]))
				}
			}
		}
	}(client)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gateway.ServeWs(command, hub, w, r)
	})

	log.Info(fmt.Sprintf("Running: %s", *addr))
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
