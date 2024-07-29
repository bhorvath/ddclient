package ipaddress

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Return the current public IP address of the client
func TestReturnCurrentIPAddress(t *testing.T) {
	m := NewMockIpifyAPI()
	m.setupRoutes()
	defer m.svr.Close()
	h := NewIpifyIPAddressHandler(m.svr.URL)

	got, err := h.GetCurrent()
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	want := "10.0.0.1"
	if got.String() != want {
		t.Errorf("Got: %v; want: %v", got, want)
	}

}

type MockIpifyAPI struct {
	svr      *httptest.Server
	response string
	calls    int
}

func NewMockIpifyAPI() *MockIpifyAPI {
	return &MockIpifyAPI{response: "10.0.0.1"}
}

func (m *MockIpifyAPI) setupRoutes() {
	mux := http.NewServeMux()
	svr := httptest.NewServer(mux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m.calls++
		fmt.Fprintf(w, m.response)
	})
	m.svr = svr
}
