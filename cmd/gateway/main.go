package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	projectID          = "iota-3345"
	registryID         = "demo-registry"
	gatewayID          = "demo-gateway0"
	cloudRegion        = "us-central1"
	privateKeyFile     = "rsa_private.pem"
	publicKeyFile      = "rsa_cert.pem"
	algorithm          = "RS256"
	mqttBridgeHostname = "tls://mqtt.googleapis.com"
	mqttBridgePort     = "8883"
	jwtExpiresMinutes  = 1200
	protocolVersion    = 4
)

func main() {
	log.SetReportCaller(true)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// onConnect defines the on connect handler which resets backoff variables.
	var onConnect mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Info(fmt.Sprintf("Client connected %s:%s: %t\n", mqttBridgeHostname, mqttBridgePort, client.IsConnected()))
	}

	// onMessage defines the message handler for the mqtt client.
	var onMessage mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Info(fmt.Sprintf("Topic: %s Message: %s\n", msg.Topic(), msg.Payload()))
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

	clientID := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, cloudRegion, registryID, gatewayID)

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

	gatewayTopic := fmt.Sprintf("/devices/%s/commands/#", gatewayID)
	if token := client.Subscribe(gatewayTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
	}
	log.Info(fmt.Sprintf("Connected to topic: %s", gatewayTopic))

	<-c

	log.Info(fmt.Sprintf("Disconnecting from: %s:%s", mqttBridgeHostname, mqttBridgePort))
	client.Disconnect(10)
}
