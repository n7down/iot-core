package gateway

import (
	//log "github.com/sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DeviceManager struct {
	client mqtt.Client
}

func NewDeviceManager(c mqtt.Client) *DeviceManager {
	return &DeviceManager{
		client: c,
	}
}

func (d DeviceManager) Receive(message string) {
	// TODO: split up the deviceID and the message
	// TODO: if detach - call detach
	// TODO: if attach - call attach
	// TODO: if subscribe - call subscribe
}
