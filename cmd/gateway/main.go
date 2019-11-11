package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	projectID          = "iota-3345"
	registryID         = "demo-registry"
	gatewayID          = "demo-gateway0"
	cloudRegion        = "us-central1"
	privateKeyFile     = "../../rsa_private.pem"
	algorithm          = "RS256"
	caCerts            = "../../roots.pem"
	mqttBridgeHostname = "mqtt.googleapis.com"
	mqttBridgePort     = "8883"
	jwtExpiresMinutes  = 1200
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	log.SetReportCaller(true)

	signBytes, err := ioutil.ReadFile(privateKeyFile)
	log.Fatal(err)

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	log.Fatal(err)

	// Declare the token with the algorithm used for signing, and the claims
	t := jwt.New(jwt.GetSigningMethod(algorithm))
	t.Claims = jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(jwtExpiresMinutes * time.Minute).Unix(),
		Audience:  projectID,
	}

	tokenString, err := t.SignedString(signKey)
	log.Fatal(err)

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("%s:%s", mqttBridgeHostname, mqttBridgePort))
	opts.SetClientID("go-simple")
	opts.SetUsername("unused")
	opts.SetPassword(tokenString)
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("go-mqtt/sample", 0, false, text)
		token.Wait()
	}

	time.Sleep(3 * time.Second)

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
