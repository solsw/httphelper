package httphelper

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/solsw/generichelper"
)

type O struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// newServer returns a server replying with the given status and body, plus a client and request targeting it.
func newServer(t *testing.T, status int, body string) (*http.Client, *http.Request) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	rq, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	return srv.Client(), rq
}

func TestReqJSON_Success(t *testing.T) {
	cl, rq := newServer(t, http.StatusOK, `{"name":"Alice","age":30}`)
	got, err := ReqJSON[O, generichelper.NoType](cl, rq, IsNotStatus2xx)
	if err != nil {
		t.Fatalf("ReqJSON() error = %v", err)
	}
	want := &O{Name: "Alice", Age: 30}
	if got == nil || *got != *want {
		t.Errorf("ReqJSON() = %v, want %v", got, want)
	}
}

func TestReqJSON_ErrorWithObject(t *testing.T) {
	cl, rq := newServer(t, http.StatusBadRequest, `{"I":7,"S":"bad"}`)
	_, err := ReqJSON[O, E](cl, rq, IsNotStatus2xx)
	herr, ok := err.(*Error[E])
	if !ok {
		t.Fatalf("ReqJSON() error type = %T, want *Error[E]", err)
	}
	if herr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", herr.StatusCode, http.StatusBadRequest)
	}
	if (herr.Object != E{I: 7, S: "bad"}) {
		t.Errorf("Object = %v, want {7 bad}", herr.Object)
	}
	if herr.Message != "" {
		t.Errorf("Message = %q, want empty", herr.Message)
	}
}

func TestReqJSON_ErrorWithText(t *testing.T) {
	cl, rq := newServer(t, http.StatusInternalServerError, "boom")
	_, err := ReqJSON[O, E](cl, rq, IsNotStatus2xx)
	herr, ok := err.(*Error[E])
	if !ok {
		t.Fatalf("ReqJSON() error type = %T, want *Error[E]", err)
	}
	if herr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", herr.StatusCode, http.StatusInternalServerError)
	}
	if herr.Message != "boom" {
		t.Errorf("Message = %q, want %q", herr.Message, "boom")
	}
}

func TestReqJSON_ErrorEmptyBody(t *testing.T) {
	cl, rq := newServer(t, http.StatusNotFound, "")
	_, err := ReqJSON[O, E](cl, rq, IsNotStatus2xx)
	herr, ok := err.(*Error[E])
	if !ok {
		t.Fatalf("ReqJSON() error type = %T, want *Error[E]", err)
	}
	if herr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", herr.StatusCode, http.StatusNotFound)
	}
	if herr.Message != "" || (herr.Object != E{}) {
		t.Errorf("got Object=%v Message=%q, want empty", herr.Object, herr.Message)
	}
}

func TestReqBody_TooLarge(t *testing.T) {
	orig := MaxResponseBodySize
	MaxResponseBodySize = 8
	t.Cleanup(func() { MaxResponseBodySize = orig })

	cl, rq := newServer(t, http.StatusOK, strings.Repeat("x", 100))
	_, err := ReqBody[generichelper.NoType](cl, rq, IsNotStatus2xx)
	if err != ErrResponseBodyTooLarge {
		t.Errorf("ReqBody() error = %v, want %v", err, ErrResponseBodyTooLarge)
	}
}

func TestJSONReader(t *testing.T) {
	r, err := JSONReader(nil)
	if err != nil {
		t.Fatalf("JSONReader(nil) error = %v", err)
	}
	if r != nil {
		t.Errorf("JSONReader(nil) = %v, want nil", r)
	}

	r, err = JSONReader(O{Name: "Bob", Age: 5})
	if err != nil {
		t.Fatalf("JSONReader() error = %v", err)
	}
	bb, _ := io.ReadAll(r)
	if got, want := string(bb), `{"name":"Bob","age":5}`; got != want {
		t.Errorf("JSONReader() = %s, want %s", got, want)
	}
}

func TestAuthBasic(t *testing.T) {
	// RFC 7617 example: "Aladdin":"open sesame" -> "QWxhZGRpbjpvcGVuIHNlc2FtZQ=="
	if got, want := AuthBasic("Aladdin", "open sesame"), "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ=="; got != want {
		t.Errorf("AuthBasic() = %q, want %q", got, want)
	}
}

func TestStatusCheckers(t *testing.T) {
	tests := []struct {
		code          int
		notOK, not2xx bool
	}{
		{http.StatusOK, false, false},
		{http.StatusCreated, true, false},
		{http.StatusNoContent, true, false},
		{http.StatusMultipleChoices, true, true},
		{http.StatusNotFound, true, true},
	}
	for _, tt := range tests {
		rs := &http.Response{StatusCode: tt.code}
		if got := IsNotStatusOK(rs); got != tt.notOK {
			t.Errorf("IsNotStatusOK(%d) = %v, want %v", tt.code, got, tt.notOK)
		}
		if got := IsNotStatus2xx(rs); got != tt.not2xx {
			t.Errorf("IsNotStatus2xx(%d) = %v, want %v", tt.code, got, tt.not2xx)
		}
	}
}
