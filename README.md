# SDR FM Audio Example
This example service utilizes software-defined radio (SDR) to extract audio streams from FM radio broadcasts. It leverages the popular open-source `librtlsdr` library and related utilities like `rtl_fm` and `rtl_power` to interface with SDR hardware connected to the edge node. Specifically, it controls the tuning and data capture from the hardware to receive transmissions at specified frequencies. For example, it can tune to an FM station frequency, acquire the audio broadcast data, and then process it.

SDR shifts the radio signal processing like tuning and demodulation from specialized analog hardware to software running on a computer's digital CPU. This allows SDR to access a wider range of the radio spectrum than a traditional analog receiver locked to a single band like FM. If physical SDR hardware is not available, the service can use mock data that emulates a radio transmission, providing flexibility. So even without a real over-the-air signal, the service can ingest audio streams as if they were broadcast over FM radio.

For more information about the `librtlsdr` library and the `rtl_fm` usage, see [librtlsdr](https://github.com/librtlsdr/librtlsdr) and [Rtl_fm Guide](http://kmkeen.com/rtl-demod-guide/). Overall, this service exemplifies using SDR to acquire and process wireless signals at the edge.

## Install Dependencies
In order to extract IQ data from the SDR hardware, the `librtlsdr` binaries have to be installed on the host machine.

First, have `gcc`, `g++`, and `make` installed. 

```bash
sudo apt-get update
sudo apt-get -y install build-essential
```

Then install `cmake` and `libusb`.

```bash
sudo apt-get -y install cmake libusb-1.0-0-dev
```

Next, build and install `librtlsdr` binaries and libraries.

```bash
git clone https://github.com/minghsu0107/librtlsdr
cd librtlsdr
mkdir build && cd build
cmake ../
sudo make && sudo make install
```
After building and installing librtlsdr, the files are located in the following directories:
- Header files are installed to `/usr/local/include`
- Library files are installed to `/usr/local/lib`
- Executable binaries are installed to `/usr/local/bin`
### For Mac Users (Apple Chips)
Install  `cmake` and `libusb` via Homebrew.
```bash
brew install cmake libusb
```
Check the version and library paths of `libusb`.
```bash
brew ls libusb
```
Build and install the librtlsdr binaries and libraries, setting the appropriate configuration and library paths for the system. For example, on arm64 with libusb version `1.0.26`:
```bash
git clone https://github.com/minghsu0107/librtlsdr
cd librtlsdr
mkdir build && cd build
cmake -DCMAKE_HOST_SYSTEM_PROCESSOR:STRING=arm64 -DLIBUSB_INCLUDE_DIR=/opt/homebrew/Cellar/libusb/1.0.26/include/libusb-1.0 -DLIBUSB_LIBRARY=/opt/homebrew/lib/libusb-1.0.dylib ../
sudo make && sudo make install
```
## Getting Started
Start a `rtl_rpcd` daemon on the host machine, which allows remote access of SDR hardware at `127.0.0.1:40000` via `librtlsdr` command-line tools.

```bash
RTLSDR_RPC_SERV_ADDR=127.0.0.1 RTLSDR_RPC_SERV_PORT=40000 rtl_rpcd &
```

Run the API server inside a container, which retrieves IQ data remotely from the `rtl_rpcd` daemon on the host machine and exposes audio data via HTTP APIs:

```bash
docker run -d --rm -p 8080:8080 -e RTLSDR_RPC_SERV_ADDR=host.docker.internal -e RTLSDR_RPC_SERV_PORT=40000 minghsu0107/rtlsdr-example-api
```
The API server is now configured to listen on port 8080 inside the container, which is forwarded from port 8080 on the host machine. This allows the API server to be accessed from the host machine at `http://localhost:8080`.
## APIs
### `/freqs`
Get a list of the frequencies of strong radio stations.

```bash
curl localhost:8080/freqs
```

Example response:

If the SDR hardware is present it will return a list of string FM stations.

`{"origin":"sdr_hardware","freqs":[89700000,91100000,91900000,93300000,94500000,95700000,96100000,97700000,99100000,101500000,102300000,103700000,105100000,107900000]}`

If the SDR hardware is not present or can not be used for some reason it will return a single station of frequency 0.

`{"origin":"fake","freqs":[0]}`

### `/audio/<freq>`
Get a 30 second chunk of raw audio.

```bash
curl localhost:8080/audio/99700000 | play --rate 32k -t raw -e s -b 16 -c 1 -V1 -
```

If the SDR hardware is not present or can not be used for some reason you can curl the fake station at frequency 0.

```bash
curl localhost:8080/audio/0 | play --rate 32k -t raw -e s -b 16 -c 1 -V1 -
```
