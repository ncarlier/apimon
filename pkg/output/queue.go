package output

import "github.com/ncarlier/apimon/pkg/model"

// Queue metric queue
var Queue = make(chan model.Metric)
