package sake_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tenntenn/sake"
)

type param struct {
	N int
	S string
}

func (p *param) Set(r *http.Request) error {
	n, err := strconv.Atoi(r.FormValue("n"))
	if err != nil {
		return err
	}

	p.N = n
	p.S = r.FormValue("s")

	return nil
}

func TestStandard(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		param   string
		want    *param
		wanterr bool
	}{
		"noerror": {"n=10&s=hoge", &param{10, "hoge"}, false},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			handler := sake.HandlerFunc[param, *param](func(ctx context.Context, w http.ResponseWriter, r *sake.Request[param, *param]) error {
				if diff := cmp.Diff(r.Param, tt.want); diff != "" {
					t.Error(diff)
				}
				return nil
			})
			h := sake.Standard[param, *param](handler, func(w http.ResponseWriter, err error) {
				switch {
				case tt.wanterr && err == nil:
					t.Error("expected error did not occur")
				case !tt.wanterr && err != nil:
					t.Error("unexpected error:", err)
				}
			})
			r := httptest.NewRequest("GET", "/?"+tt.param, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
		})
	}
}
