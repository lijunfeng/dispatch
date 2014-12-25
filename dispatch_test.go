package dispatch

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"bytes"
	"testing"

	"github.com/lunny/tango"
	"github.com/lunny/log"
)

func TestDispatch(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	l := log.Std
	l.SetOutputLevel(log.Ldebug)

	t1 := tango.NewWithLog(l)
	t1.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("tango 1"))
	}))

	t2 := tango.NewWithLog(l)
	t2.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("tango 2"))
	}))

	dispatch := New(map[string]*tango.Tango{
		"/": t1,
		"/api/": t2,
	})

	t3 := tango.NewWithLog(l)
	t3.Use(dispatch)

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	t3.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "tango 1")

	req, err = http.NewRequest("GET", "http://localhost:8000/api/", nil)
	if err != nil {
		t.Error(err)
	}

	buff.Reset()

	t3.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "tango 2")
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}