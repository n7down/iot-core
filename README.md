# IOT-Core

## Instructions
- Generate an RSA256 private key and corresponding X509 public certificate that will be used for authenticating the gateway device - `generate an RSA256 private key and corresponding X509 public certificate that will be used for authenticating your gateway device - `openssl req -x509 -newkey rsa:2048 \
  -nodes -keyout rsa_private.pem -x509 -days 365 -out rsa_cert.pem -subj "/CN=unused"`
- Create a gateway with - `gcloud iot devices create demo-gateway0 --region=us-central1 --registry=demo-registry --auth-method=device-auth-token-only --device-type=gateway --public-key path=rsa_cert.pem,type=rsa-x509-pem`

## TODO
1. [ ] Create gateway  and connect to IOT core with `device credentials only` (for proximity of real gateways)
2. [ ] Create a device - attach to gateway using JWT
 - [MQTT Bridge](https://cloud.google.com/iot/docs/how-tos/gateways/mqtt-bridge)
3. [ ] Send and receive command from IOT core to device
 - [Commands](https://cloud.google.com/iot/docs/how-tos/commands)
4. [ ] Send info logs to stackdriver - seperate by site
5. [ ] Send command to device to send logs of what gateway it is attached to

## Notes
- [Cloud IOT Core Gateways](https://codelabs.developers.google.com/codelabs/cloud-iot-core-gateways/index.html#0)

