# FM Audio Stream on NATS JetStream
## Usage
Start the `rtl_rpcd` daemon on the host machine.
```bash
RTLSDR_RPC_SERV_ADDR=127.0.0.1 RTLSDR_RPC_SERV_PORT=40000 rtl_rpcd &
```
Start FM audio publisher and NATS JetStream.
```bash
docker-compose up -d
```
Start FM audio subscriber, which subscribes the FM audio stream and plays it in real-time.
```bash
NATS_URL=nats://mytoken@127.0.0.1:4222 go run main.go
```