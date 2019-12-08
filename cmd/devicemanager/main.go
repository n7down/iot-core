package main

import (
	"fmt"
	"net"

	"github.com/n7down/iot-core/internal/pb/devicemanager"
	servers "github.com/n7down/iot-core/internal/servers/devicemanager"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = "8082"
	cert = "cert.pem"
	key  = "key.pem"

	projectID   = "iota-3345"
	registryID  = "demo-registry"
	cloudRegion = "us-central1"
)

// TODO: gateways talk to this service to see if devices are created

func main() {
	log.SetReportCaller(true)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	deviceManagerServer, err := servers.NewDeviceManagerServer()
	if err != nil {
		log.Fatal("unable to start device manager server")
	}

	log.Infof("Listening on port: %s\n", port)
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	devicemanager_pb.RegisterDeviceManagerServiceServer(grpcServer, deviceManagerServer)
	grpcServer.Serve(lis)
}
