package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/gordonklaus/portaudio"
	nc "github.com/nats-io/nats.go"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
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
	subscribeOptions := []nc.SubOpt{
		nc.DeliverNew(),
		nc.AckExplicit(),
	}
	jsConfig := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    false,
		SubscribeOptions: subscribeOptions,
		TrackMsgId:       false,
		AckAsync:         false,
		DurablePrefix:    "",
	}
	subscriber, err := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:              natsURL,
			QueueGroupPrefix: "",
			SubscribersCount: 1,
			CloseTimeout:     time.Minute,
			NatsOptions:      options,
			Unmarshaler:      marshaler,
			JetStream:        jsConfig,
		},
		logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := subscriber.Subscribe(context.Background(), natsSubject)
	if err != nil {
		fmt.Println(err)
		return
	}

	out := make([]int16, 8192)
	// get audio format information
	rate := 16000

	portaudio.Initialize()
	defer portaudio.Terminate()
	// stereo audio
	stream, err := portaudio.OpenDefaultStream(0, 2, float64(rate), len(out), &out)
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

	go func() {
		<-sig
		if err = subscriber.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for msg := range msgs {
		fmt.Println("sub: ", msg.UUID)
		audio := make([]byte, 2*8192)
		copy(audio, msg.Payload)
		if err = binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out); err != nil {
			fmt.Println(err)
		}
		if err = stream.Write(); err != nil {
			fmt.Println(err)
		}
		msg.Ack()
	}
}
