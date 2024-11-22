package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const ServerPort = ":3000"
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()

	if err := r.sendRequiredServices(reg); err != nil {
		return err
	}

	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})

	return nil
}

func (r registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequiredService {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, added := range fullPatch.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					if err := r.sendPatch(p, reg.ServiceUpdateURL); err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

func (r registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqReg := range reg.RequiredService {
			if reqReg == serviceReg.ServiceName {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	if err := r.sendPatch(p, reg.ServiceUpdateURL); err != nil {
		return err
	}
	return nil
}

func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	if _, err = http.Post(url, "application/json", bytes.NewBuffer(d)); err != nil {
		return err
	}
	return nil
}

func (r *registry) remove(url string) error {
	for i, reg := range r.registrations {
		if reg.ServiceURL == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: reg.ServiceName,
						URL:  reg.ServiceURL,
					},
				},
			})
			return nil
		}
	}
	return fmt.Errorf("service at URL: %s not found", url)
}


func (r *registry) heartBeat(t time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registrations {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var success = true
				for i := 0; i < 3; i++ {
					res, err := http.Get(reg.HeartBeatURL)
					if err != nil {
						log.Println(err)
					}
					if res.StatusCode != http.StatusOK {
						r.remove(reg.ServiceURL)
						success = false
					} else {
						if !success {
							r.add(reg)
						}
						break
					}
				}
				if success {
					log.Printf("health check at service: %s, state is ok", reg.ServiceName)
				} else {
					log.Printf("health check at service: %s, state is bad", reg.ServiceName)
				}
				time.Sleep(t)
			}()
		}
		wg.Wait()
	}
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartBeat(3 * time.Second)
	})
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		if err = reg.add(r); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		paload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(paload)
		log.Println("Removing Service at URL: " + url)
		if err = reg.remove(url); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
