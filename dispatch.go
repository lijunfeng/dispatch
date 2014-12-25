package dispatch

import (
	"strings"
	"net/http"

	"github.com/lunny/tango"
)

type Dispatch struct {
	logger tango.Logger
	tangos map[string]*tango.Tango
}

func New(m map[string]*tango.Tango) *Dispatch {
	if m != nil {
		return &Dispatch{tangos: m}
	}
	return &Dispatch{tangos: make(map[string]*tango.Tango)}
}

func (d *Dispatch) Add(name string, t *tango.Tango) {
	d.tangos[name] = t
}

func (d *Dispatch) SetLogger(logger tango.Logger) {
	d.logger = logger
}

func (d *Dispatch) Handle(ctx *tango.Context) {
	fields := strings.Split(ctx.Req().URL.Path, "/")
	if len(fields) == 2 {
		if t, ok := d.tangos["/"]; ok {
			t.ServeHTTP(ctx.ResponseWriter, ctx.Req())
			return
		}
	} else {
		if t, ok := d.tangos[strings.Join(fields[0:2], "/")+"/"]; ok {
			t.ServeHTTP(ctx.ResponseWriter, ctx.Req())
			return
		}
	}

	ctx.WriteHeader(http.StatusNotFound)
	ctx.Write([]byte("Not Found"))
}