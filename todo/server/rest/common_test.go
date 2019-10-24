package rest

import (
	"bytes"
	"encoding/json"
	"github.com/vkonstantin/wg/todo/controller"
	"github.com/vkonstantin/wg/todo/storage/memory"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newTestService() *restTestService {
	mainService := controller.NewMainService(memory.NewDefault())
	s := New("", mainService).(*rest)
	rs := restTestService{
		service: s,
	}
	return &rs
}

type restTestService struct {
	service *rest
}

func (s *restTestService) request(method, path string, body interface{}, headers ...string) (code int, respBody string, err error) {
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		return 0, "", err
	}
	if len(headers) > 0 && len(headers)%2 == 0 {
		for i := 0; i < len(headers); i += 2 {
			key := headers[i]
			value := headers[i+1]
			req.Header.Add(key, value)
		}
	}

	w := httptest.NewRecorder()
	s.service.engine.ServeHTTP(w, req)

	respBody = w.Body.String()
	return w.Code, respBody, nil
}

func randRequestID() string {
	i := rand.Intn(10000000)
	s := strconv.Itoa(i)
	return s
}
