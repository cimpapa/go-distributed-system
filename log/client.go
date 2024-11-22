package log

import (
	"bytes"
	"fmt"
	stlog "log"
	"luuk/distributed/registry"
	"net/http"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct{
	url string
}
func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer(data)
	res, err := http.Post(cl.url + "/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("fail to send log message. Service response with code: %d", res.StatusCode)
	}
	return len(data), nil
}