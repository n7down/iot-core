package gateway

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type DeviceManager struct {
	client mqtt.Client
	hub    *Hub
}

func NewDeviceManager(c mqtt.Client, h *Hub) *DeviceManager {
	return &DeviceManager{
		client: c,
		hub:    h,
	}
}

func (d DeviceManager) Detach(deviceID string) {
	detachTopic := fmt.Sprintf("/devices/%s/detach", deviceID)
	if token := d.client.Subscribe(detachTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
	}
	log.Info(fmt.Sprintf("Connected to topic: %s", detachTopic))
}

func (d DeviceManager) Attach(deviceID string) {
	attachTopic := fmt.Sprintf("/devices/%s/attach", deviceID)
	if token := d.client.Subscribe(attachTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
	}
	log.Info(fmt.Sprintf("Connected to topic: %s", attachTopic))
}

func (d DeviceManager) Subscribe(deviceID string) {
	subscribeTopic := fmt.Sprintf("/devices/%s/commands/#", deviceID)
	if token := d.client.Subscribe(subscribeTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to topic: %s", token.Error()))
	}
	log.Info(fmt.Sprintf("Connected to topic: %s", subscribeTopic))
}
