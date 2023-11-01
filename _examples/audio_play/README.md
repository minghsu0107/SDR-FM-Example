# Real-time FM Audio Stream
This example demonstrates how to process and play real-time FM audio stream.
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
## Getting Started
Start the `rtl_rpcd` daemon on the host machine:
```bash
RTLSDR_RPC_SERV_ADDR=127.0.0.1 RTLSDR_RPC_SERV_PORT=40000 rtl_rpcd &
```
Play real-time FM audio stream from remote SDR hardware at frequency 94.1M:
```bash
RTLSDR_RPC_SERV_ADDR="127.0.0.1" RTLSDR_RPC_SERV_PORT="40000" go run main.go 94100000
```