# FM Audio Stream on NATS JetStream
## Dependencies
To enable real-time local playback of FM audio, the following dependencies are required.

For Debian / Ubuntu Linux:
```bash
apt-get install -y pkg-config portaudio19-dev
```
For OS X:
```bash
brew install pkg-config portaudio
```
## Usage
Start the `rtl_rpcd` daemon on the host machine.
```bash
RTLSDR_RPC_SERV_ADDR=127.0.0.1 RTLSDR_RPC_SERV_PORT=40000 rtl_rpcd &
```
### Local Setup
Start the FM audio publisher and NATS JetStream. The publisher extracts FM audio at frequency 99.7M and streams the data to NATS JetStream.
```bash
docker-compose up -d
```
Start the FM audio subscriber locally to listen to the audio stream in real time. The subscriber will connect to NATS JetStream and play the FM audio that is being published.

To start the subscriber:
```bash
NATS_URL=nats://mytoken@127.0.0.1:4222 go run main.go
```
### Remote NATS
Instead of streaming audio data locally, you can also send it to a remote NATS JetStream cluster deployed in the cloud. For example, if you have access to a NATS server at nats://mytoken@1.2.3.4:4222, you can start the publisher container like this:
```bash
docker run --rm -e RTLSDR_RPC_SERV_ADDR=host.docker.internal -e RTLSDR_RPC_SERV_PORT=40000 -e NATS_URL=nats://mytoken@1.2.3.4:4222 minghsu0107/rtlsdr-example-nats-client 99700000
```
This will send the audio streams to the remote NATS server rather than a local NATS container.

Start the subscriber to listen to the audio stream:
```bash
NATS_URL=nats://mytoken@1.2.3.4:4222 go run main.go
```