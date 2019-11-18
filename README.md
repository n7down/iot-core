# IOT-Core

## Instructions
1. Generate an RSA256 private key and corresponding X509 public certificate that will be used for authenticating the gateway device with `openssl req -x509 -newkey rsa:2048 -nodes -keyout rsa_private.pem -x509 -days 365 -out rsa_cert.pem -subj "/CN=unused"`
2. Create a registry with `gcloud iot registries create demo-registry0 --region=us-central1 --event-notification-config=topic=gateway-topic --project=<PROJECT_ID>`
3. Create a gateway with - `gcloud iot devices create demo-gateway0 --region=us-central1 --registry=demo-registry --auth-method=device-auth-token-only --device-type=gateway --public-key path=rsa_cert.pem,type=rsa-x509-pem`
4. Create a device with `gcloud iot devices create device0 --region=us-central1 --registry=demo-registry --device-type=non-gateway --project=<PROJECT_ID> --public-key path=rsa_cert.pem,type=rsa-x509-pem`

## TODO
1. [x] Create gateway  and connect to IOT core with `device credentials only` (for proximity of real gateways)
2. [x] Create a device - attach to gateway using JWT
 - [MQTT Bridge](https://cloud.google.com/iot/docs/how-tos/gateways/mqtt-bridge)
3. [x] Send and receive command from IOT core to device
 - [Commands](https://cloud.google.com/iot/docs/how-tos/commands)

## Notes
- [Cloud IOT Core Gateways](https://codelabs.developers.google.com/codelabs/cloud-iot-core-gateways/index.html#0)
- [UDP in GO](https://jameshfisher.com/2016/11/17/udp-in-go/)
