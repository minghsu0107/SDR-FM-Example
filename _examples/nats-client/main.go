package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
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

	natsURL := getenv("NATS_URL", "nats://127.0.0.1:4222")
	natsSubject := getenv("NATS_SUBJECT", "fm.raw")

	marshaler := &nats.GobMarshaler{}
	logger := watermill.NewStdLogger(false, false)
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}
	publisher, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:         natsURL,
			NatsOptions: options,
			Marshaler:   marshaler,
			JetStream: nats.JetStreamConfig{
				Disabled:       false,
				AutoProvision:  false,
				PublishOptions: nil,
				TrackMsgId:     false,
				AckAsync:       false,
			},
		},
		logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	var i uint64 = 0
	audio := make([]byte, 2*8192)
	for {
		n, err := stdout.Read(audio)
		for n < 16384 {
			bytesRead, err := stdout.Read(audio[n:])
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
			}
			n += bytesRead
			if n < 16384 {
				time.Sleep(10 * time.Millisecond)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		payload := make([]byte, 2*8192)
		copy(payload, audio)
		msg := message.NewMessage(strconv.FormatUint(i, 10), payload)
		if err := publisher.Publish(natsSubject, msg); err != nil {
			fmt.Println(err)
		}
		fmt.Println("pub: ", msg.UUID)
		i++
		select {
		case <-sig:
			if err = cmd.Process.Kill(); err != nil {
				fmt.Println(err)
			}
			if err = publisher.Close(); err != nil {
				fmt.Println(err)
			}
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
