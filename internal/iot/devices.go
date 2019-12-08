package iot

import (
	"fmt"
	"io/ioutil"

	"context"

	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

// createDevice creates a device in a registry with one of the following public key formats:
// RSA_PEM, RSA_X509_PEM, ES256_PEM, ES256_X509_PEM, UNAUTH.
func createDevice(projectID string, region string, registryID string, deviceID string, publicKeyFormat string, keyPath string) (*cloudiot.Device, error) {

	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	client, err := cloudiot.New(httpClient)
	if err != nil {
		return nil, err
	}

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	var device cloudiot.Device

	// If no credentials are passed in, create an unauth device.
	if publicKeyFormat == "UNAUTH" {
		device = cloudiot.Device{
			Id: deviceID,
		}
	} else {
		device = cloudiot.Device{
			Id: deviceID,
			Credentials: []*cloudiot.DeviceCredential{
				{
					PublicKey: &cloudiot.PublicKeyCredential{
						Format: publicKeyFormat,
						Key:    string(keyBytes),
					},
				},
			},
		}
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}

// deleteDevice deletes a device from a registry.
func deleteDevice(projectID string, region string, registryID string, deviceID string) (*cloudiot.Empty, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.Delete(path).Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}

// listDevices gets the identifiers of devices for a specific registry.
func listDevices(projectID string, region string, registryID string) ([]*cloudiot.Device, error) {
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
	response, err := client.Projects.Locations.Registries.Devices.List(parent).Do()
	if err != nil {
		return nil, err
	}

	return response.Devices, nil
}

// getDevice retrieves a specific device and prints its details.
func getDevice(projectID string, region string, registryID string, device string) (*cloudiot.Device, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, device)
	response, err := client.Projects.Locations.Registries.Devices.Get(path).Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}
