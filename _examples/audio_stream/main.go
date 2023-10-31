package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gordonklaus/portaudio"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func captureAudio(freq int) {
	cmd := exec.Command("rtl_fm", "-M", "fm", "-s", "170k", "-o", "4", "-A", "fast", "-r", "32k", "-l", "0", "-E", "deemp", "-f", strconv.Itoa(freq))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	rtlsdrRpcServAddr := getenv("RTLSDR_RPC_SERV_ADDR", "127.0.0.1")
	rtlsdrRpcServPort := getenv("RTLSDR_RPC_SERV_PORT", "40000")
	cmd.Env = append(cmd.Env, "RTLSDR_RPC_IS_ENABLED=1", "RTLSDR_RPC_SERV_ADDR="+rtlsdrRpcServAddr, "RTLSDR_RPC_SERV_PORT="+rtlsdrRpcServPort)
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// get audio format information
	rate := 32000

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int16, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, float64(rate), len(out), &out)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stream.Close()

	if err = stream.Start(); err != nil {
		fmt.Println(err)
		return
	}
	defer stream.Stop()
	for {
		audio := make([]byte, 2*len(out))
		_, err = stdout.Read(audio)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		if err = binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out); err != nil {
			fmt.Println(err)
			return
		}
		if err = stream.Write(); err != nil {
			fmt.Println(err)
			return
		}
		select {
		case <-sig:
			return
		default:
		}
	}
}

func main() {
	freqStr := os.Args[1]

	freq, err := strconv.Atoi(freqStr)
	if err != nil {
		fmt.Println("Error converting argument to integer:", err)
		return
	}
	captureAudio(freq)
}
