package devicemanager

import (
	"context"

	devicemanager_pb "github.com/n7down/iot-core/internal/pb/devicemanager"
)

type DeviceManagerServer struct {
	db *DeviceManagerDB

	projectID   string
	registryID  string
	cloudRegion string
}

func NewDeviceManagerServer() (*DeviceManagerServer, error) {
	deviceManagerServer := &DeviceManagerServer{}
	db, err := NewDeviceManagerDB()
	if err != nil {
		return deviceManagerServer, err
	}
	return &DeviceManagerServer{db: db}, nil
}

func (s *DeviceManagerServer) DeviceCreate(ctx context.Context, req *devicemanager_pb.DeviceCreateRequest) (*devicemanager_pb.DeviceCreateResponse, error) {

	// TODO: check if the device has been created
	// TODO: if it has not then create it

	res := &devicemanager_pb.DeviceCreateResponse{}
	return res, nil
}

func (s *DeviceManagerServer) GatewayCreate(ctx context.Context, req *devicemanager_pb.GatewayCreateRequest) (*devicemanager_pb.GatewayCreateResponse, error) {
	res := &devicemanager_pb.GatewayCreateResponse{}
	return res, nil
}
