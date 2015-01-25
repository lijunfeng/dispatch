package dispatch

import (
	"path"
	"sort"
	"strings"

	"github.com/lunny/tango"
)

type mTango struct {
	route string
	tg    *tango.Tango
}

type Dispatch []*mTango

func (d *Dispatch) Len() int {
	return len(*d)
}

func (d *Dispatch) Less(i, j int) bool {
	return len((*d)[i].route) > len((*d)[j].route)
}

func (d *Dispatch) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

func New(m map[string]*tango.Tango) *Dispatch {
	var dispatch Dispatch = make([]*mTango, 0)
	if m == nil {
		m = make(map[string]*tango.Tango)
	}
	for k, t := range m {
		dispatch = append(dispatch, &mTango{
			route: k,
			tg:    t,
		})
	}
	sort.Sort(&dispatch)
	return &dispatch
}

func (d *Dispatch) Add(name string, t *tango.Tango) {
	*d = append(*d, &mTango{
		route: name,
		tg:    t,
	})
	sort.Sort(d)
}

func (d *Dispatch) Handle(ctx *tango.Context) {
	var tg *mTango
	for _, t := range *d {
		if strings.HasPrefix(ctx.Req().URL.Path, t.route) {
			tg = t
			break
		}
	}

	if tg != nil {
		ctx.Req().URL.Path = path.Join("/", strings.TrimLeft(ctx.Req().URL.Path, tg.route))
		tg.tg.ServeHTTP(ctx.ResponseWriter, ctx.Req())
	} else {
		ctx.NotFound()
	}
}
