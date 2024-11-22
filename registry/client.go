package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	stlog "log"
	"luuk/distributed/log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	heartBeatURL, err := url.Parse(r.HeartBeatURL)
	if err != nil {
		return err
	}
	http.HandleFunc(heartBeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{
		Name: r.ServiceName,
	})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(r); err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. Registry service "+
			"respond with code %v", res.StatusCode)
	}
	return nil
}

func UnRegiserService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesURL,
		bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to unregister Service at URL: %s,"+
			" Service State Code is %d", url, resp.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct {
	Name ServiceName
}

func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	if err := dec.Decode(&p); err != nil {
		stlog.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Update Msg received %v\n", p)
	prov.Update(p)
	prov.SetLogger(p, suh.Name)
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

func (p *providers) SetLogger(pat patch, name ServiceName) {
	for _, pE := range pat.Added {
		if pE.Name == LogService {
			if logProvider, err := GetProvider(LogService); err == nil {
				fmt.Printf("Logging service found at: %s\n", logProvider)
				log.SetClientLogger(logProvider, string(name))
				// 调用服务进行测试，是可以写入到日志服务的文件中的
				stlog.Printf("Logging service found at: %s\n", logProvider)
			}
		}
	}
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, pA := range pat.Added {
		if _, ok := p.services[pA.Name]; !ok {
			p.services[pA.Name] = make([]string, 0)
		}
		p.services[pA.Name] = append(p.services[pA.Name], pA.URL)
	}
	for _, pA := range pat.Removed {
		if providerURLs, ok := p.services[pA.Name]; ok {
			for i := range providerURLs {
				if providerURLs[i] == pA.URL {
					p.services[pA.Name] = append(providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}
}

func (p providers) get(name ServiceName) (string, error) {
	ps, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no providers available for service %v", name)
	}
	idx := int(rand.Float32() * float32(len(ps)))
	return ps[idx], nil
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}
