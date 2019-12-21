package wdb

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Datum struct {
	Raw  []byte
	Time time.Time
}

type Server struct {
	C chan func(*Server)

	data map[string]Datum
}

func NewServer() *Server {
	s := &Server{
		C:    make(chan func(*Server), 100),
		data: make(map[string]Datum),
	}

	go s.loop()
	return s
}

func (s *Server) loop() {
	for fn := range s.C {
		fn(s)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path

	switch r.Method {
	case "GET":
		dc := DC()
		s.C <- func(s *Server) {
			dc.Put(s.Get(key))
		}
		datum := <-dc

		if datum.Raw == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(datum.Raw)
	case "POST":
		raw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dc := DC()
		s.C <- func(s *Server) {
			s.Set(key, raw)
			dc.Put(Datum{})
		}
		<-dc
		w.WriteHeader(http.StatusCreated)
	case "DELETE":
		s.C <- func(s *Server) {
			s.Remove(key)
		}

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) Set(key string, raw []byte) {
	s.data[key] = Datum{
		Raw:  raw,
		Time: time.Now(),
	}
}

func (s *Server) Remove(key string) {
	delete(s.data, key)
}

func (s *Server) Get(key string) Datum {
	return s.data[key]
}

type DatumChan chan Datum

func DC() DatumChan {
	return make(chan Datum, 1)
}

func (dc DatumChan) Put(datum Datum) {
	dc <- datum
	close(dc)
}
