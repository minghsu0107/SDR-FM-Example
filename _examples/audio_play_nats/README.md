# FM Audio Stream on NATS JetStream
## Usage
Start the `rtl_rpcd` daemon on the host machine.
```bash
RTLSDR_RPC_SERV_ADDR=127.0.0.1 RTLSDR_RPC_SERV_PORT=40000 rtl_rpcd &
```
Start the FM audio publisher and NATS JetStream. The publisher extracts FM audio at frequency 99.7M and streams the data to NATS JetStream.
```bash
docker-compose up -d
```
Start the FM audio subscriber, which subscribes the FM audio stream from NATS JetStream and plays it in real-time.
```bash
NATS_URL=nats://mytoken@127.0.0.1:4222 go run main.go
```