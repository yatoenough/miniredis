package aof

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/yatoenough/miniredis/internal/resp"
)

type AOF struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &AOF{
		file: f,
		rd:   bufio.NewReader(f),
	}

	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *AOF) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *AOF) Read(callback func(value resp.Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	respr := resp.NewRESP(aof.file)

	for {
		value, err := respr.Read()
		if err == nil {
			callback(value)
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
	}

	return nil
}

func (aof *AOF) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}
