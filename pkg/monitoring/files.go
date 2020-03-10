package monitoring

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ncarlier/apimon/pkg/config"
	"github.com/ncarlier/apimon/pkg/logger"
)

func getMonitorsFromFile(filename string) []config.Monitor {
	monitors := []config.Monitor{}
	matches, err := filepath.Glob(filename)
	if err != nil {
		logger.Error.Printf("unable to decode filename pattern (%s): %s\n", filename, err.Error())
	}
	for _, _filename := range matches {
		logger.Debug.Println("loading external monitoring file:", _filename)
		data, err := ioutil.ReadFile(_filename)
		if err != nil {
			logger.Error.Printf("unable to read external configuration file (%s): %s\n", _filename, err.Error())
			continue
		}
		conf, err := config.NewConfig(data)
		if err != nil {
			logger.Error.Printf("unable to load external configuration file (%s): %s\n", _filename, err.Error())
			continue
		}
		logger.Info.Printf("%d monitors loaded from external file (%s)\n", len(conf.Monitors), _filename)
		monitors = append(monitors, conf.Monitors...)
	}
	return monitors
}

func (m *Monitoring) getFilesConfig() []config.Monitor {
	monitors := []config.Monitor{}
	filenames := m.conf.MonitorsFiles
	if len(filenames) == 0 {
		return monitors
	}
	for _, filename := range filenames {
		monitors = append(monitors, getMonitorsFromFile(filename)...)
	}
	return monitors
}
