package metric

import (
	"fmt"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
)

// Producer send metric to an output target
type Producer struct {
	Config   config.Output
	QuitChan chan bool
	Writer   Writer
}

// NewMetricProducer create new metric writer
func NewMetricProducer(conf config.Output) (*Producer, error) {
	writer, err := getMetricWriter(conf.Target, conf.Format)
	if err != nil {
		return nil, fmt.Errorf("unable to get metric writer: %s", err)
	}

	return &Producer{
		Config:   conf,
		QuitChan: make(chan bool),
		Writer:   writer,
	}, nil
}

// Start starts metric producer
func (p *Producer) Start() {
	logger.Debug.Println("Starting metric producer...")
	go func() {
		for {
			select {
			case metric := <-Queue:
				err := p.Writer.Write(metric)
				if err != nil {
					logger.Error.Println("Unable to write metric:", err)
				}
			case <-p.QuitChan:
				return
			}
		}
	}()
}

// Stop stops metric producer
func (p *Producer) Stop() {
	logger.Debug.Println("Stopping metric producer...")
	p.QuitChan <- true
	p.Writer.Close()
}
