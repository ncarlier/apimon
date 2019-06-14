package output

import (
	"fmt"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
	"github.com/ncarlier/apimon/pkg/output/writer"
)

// Provider format and send metric to an output target
type Provider struct {
	Config   config.Output
	QuitChan chan bool
	Writer   writer.Writer
}

// NewOutputProvider create new metric provider
func NewOutputProvider(conf config.Output) (*Provider, error) {
	writer, err := writer.NewOutputWriter(conf.Target, conf.Format)
	if err != nil {
		return nil, fmt.Errorf("unable to get output writer: %s", err)
	}

	return &Provider{
		Config:   conf,
		QuitChan: make(chan bool),
		Writer:   writer,
	}, nil
}

// Start starts output provider
func (p *Provider) Start() {
	logger.Debug.Println("starting output provider...")
	go func() {
		for {
			select {
			case metric := <-Queue:
				err := p.Writer.Write(metric)
				if err != nil {
					logger.Error.Println("unable to write metric:", err)
				}
			case <-p.QuitChan:
				return
			}
		}
	}()
}

// Stop stops metric provider
func (p *Provider) Stop() {
	p.QuitChan <- true
	p.Writer.Close()
	logger.Debug.Println("output provider stopped")
}
