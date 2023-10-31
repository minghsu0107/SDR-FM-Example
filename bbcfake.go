package main

import (
	"fmt"
	"os"

	mp3 "github.com/hajimehoshi/go-mp3"
)

// func downsample(input []byte) (out []byte) {
// 	out = make([]byte, len(input)/6+1)
// 	for i := 0; i < len(input); i += 12 {
// 		out[i/6] = input[i]
// 		out[(i/6)+1] = input[i+1]
// 	}
// 	return
// }

// downloadAndSplit fetches audio from a file and returns a slice of chunks.
func downloadAndSplit(path string) ([][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d, err := mp3.NewDecoder(f)
	if err != nil {
		return nil, err
	}

	chunks := make([][]byte, 0, 10997)
Loop:
	for {
		buf := make([]byte, 512)
		chunk := make([]byte, 0, 512*10997)
		for i := 0; i < 10997; i++ {
			bytesRead, err := d.Read(buf)
			if err != nil {
				chunks = append(chunks, chunk)
				break Loop
			}
			chunk = append(chunk, buf[:bytesRead]...)
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

// NewFakeRadio creates a new fake radio which will give chunks of audio from BBC news.
func NewFakeRadio() FakeRadio {
	return FakeRadio{
		chunksPointer: 0,
		chunks:        [][]byte{},
	}
}

// FakeRadio holds the state for the radio.
type FakeRadio struct {
	chunksPointer int
	chunks        [][]byte
}

func (fr *FakeRadio) refreshChunks() {
	var err error
	fr.chunks, err = downloadAndSplit("mock_audio.mp3")
	if err != nil {
		panic(err)
	}
	fr.chunksPointer = len(fr.chunks)
	fmt.Println("got", len(fr.chunks), "chunks of fake audio")
}

// GetNextChunk may fetch a new audio file
func (fr *FakeRadio) GetNextChunk() (chunk []byte) {
	if fr.chunksPointer == 0 {
		fmt.Println("no chunks, refreshing...")
		fr.refreshChunks()
	}
	fr.chunksPointer--
	return fr.chunks[fr.chunksPointer]
}
