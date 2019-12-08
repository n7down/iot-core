package iot

import (
	"fmt"
	"io/ioutil"

	"context"

	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

// createGateway creates a new IoT Core gateway with a given id, public key, and auth method.
// gatewayAuthMethod can be one of: ASSOCIATION_ONLY, DEVICE_AUTH_TOKEN_ONLY, ASSOCIATION_AND_DEVICE_AUTH_TOKEN.
// https://cloud.google.com/iot/docs/reference/cloudiot/rest/v1/projects.locations.registries.devices#gatewayauthmethod
func createGateway(projectID string, region string, registryID string, gatewayID string, gatewayAuthMethod string, publicKeyPath string) (*cloudiot.Device, error) {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	client, err := cloudiot.New(httpClient)
	if err != nil {
		return nil, err
	}

	keyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	gateway := &cloudiot.Device{
		Id: gatewayID,
		Credentials: []*cloudiot.DeviceCredential{
			{
				PublicKey: &cloudiot.PublicKeyCredential{
					Format: "RSA_X509_PEM",
					Key:    string(keyBytes),
				},
			},
		},
		GatewayConfig: &cloudiot.GatewayConfig{
			GatewayType:       "GATEWAY",
			GatewayAuthMethod: gatewayAuthMethod,
		},
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, gateway).Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}

// listGateways lists all the gateways in a specific registry.
func listGateways(projectID string, region string, registryID string) ([]*cloudiot.Device, error) {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	client, err := cloudiot.New(httpClient)
	if err != nil {
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.List(parent).GatewayListOptionsGatewayType("GATEWAY").Do()

	if err != nil {
		return nil, fmt.Errorf("ListGateways: %v", err)
	}

	if len(response.Devices) == 0 {
		return response.Devices, nil
	}

	return response.Devices, nil
}
