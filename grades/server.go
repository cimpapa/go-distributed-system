package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers() {
	handler := new(studentHandler)
	http.Handle("/student", handler)
	http.Handle("/student/", handler)
}

type studentHandler struct{}

func (sh studentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2:
		sh.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.addGrade(w, r, id)
	}
}

func (sh studentHandler) getAll(w http.ResponseWriter, _ *http.Request) {
	ssMutex.Lock()
	defer ssMutex.Unlock()

	data, err := sh.toJSON(students)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) getOne(w http.ResponseWriter, _ *http.Request, id int) {
	ssMutex.Lock()
	defer ssMutex.Unlock()

	stu, err := students.GetByID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := sh.toJSON(stu)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	ssMutex.Lock()
	defer ssMutex.Unlock()

	stu, err := students.GetByID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var g Grade
	dec := json.NewDecoder(r.Body)
	if err = dec.Decode(&g); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stu.Grades = append(stu.Grades, g)
	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(g)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) toJSON(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(obj); err != nil {
		return nil, fmt.Errorf("fail to serialize obj: %v", err)
	}
	return buf.Bytes(), nil
}
