package dotloop_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/RealGeeks/dotloop-go"
)

func TestDotloop_InvalidToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		io.WriteString(w, `{"error":"invalid_token","error_description":"Invalid access token: fb0e9121"}`)
	}))
	defer ts.Close()

	cli := &dotloop.Dotloop{
		URL:   ts.URL + "/",
		Token: "fb0e9121",
	}
	err := cli.LoopIt(dotloop.Loop{})

	equals(t, &dotloop.ErrInvalidToken{Msg: "dotloop: Invalid access token: fb0e9121"}, err)
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
