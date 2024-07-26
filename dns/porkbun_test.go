package dns

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
)

var ip, _ = netip.ParseAddr("10.0.0.1")

func TestNewPorkbunDNSHandler(t *testing.T) {
	envSetup(t)

	h, err := NewPorkbunDNSHandler("")
	if h == nil {
		t.Errorf("Instantiation failed; got: %v", h)
	}
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
}

func TestNewPorkbunDNSHandlerMissingEnvVar(t *testing.T) {
	tests := []struct {
		missingEnvVar string
	}{
		{"DOMAIN"},
		{"RECORD_NAME"},
		{"RECORD_TYPE"},
		{"PORKBUN_API_KEY"},
		{"PORKBUN_SECRET_KEY"},
	}

	for _, tt := range tests {
		want := errors.New("environment variable " + tt.missingEnvVar + " is missing")
		envSetup(t)
		t.Setenv(tt.missingEnvVar, "")

		t.Run(tt.missingEnvVar, func(t *testing.T) {
			h, err := NewPorkbunDNSHandler("")
			if h != nil {
				t.Errorf("Expected instantiation to fail; got: %v", h)
			}
			if err.Error() != want.Error() {
				t.Errorf("Got error: %v; want: %v", err, want)
			}
		})
	}
}

// The handler does not support updating multiple records.
func TestUpdateFailOnMultipleRecords(t *testing.T) {
	envSetup(t)
	m := NewMockPorkbunAPI()
	m.setupRoutes()
	defer m.svr.Close()
	h, err := NewPorkbunDNSHandler(m.svr.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	h.Update(ip)
	if m.editCalls != 0 {
		t.Errorf("Got edit calls: %v; want: 0", m.editCalls)
	}
	if m.createCalls != 0 {
		t.Errorf("Got create calls: %v; want: 0", m.createCalls)
	}

}

// If the current IP address is the same as the DNS record then don't edit or create anything.
func TestNoUpdateIfIPHasNotChanged(t *testing.T) {
	envSetup(t)
	m := NewMockPorkbunAPI()
	m.setupRoutes()
	defer m.svr.Close()
	h, err := NewPorkbunDNSHandler(m.svr.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}
	m.retrieveResponse = retrieveResponse{
		"SUCCESS", []record{
			{
				Id:      "test2",
				Content: "10.0.0.1",
			},
		},
	}

	h.Update(ip)
	if m.editCalls != 0 {
		t.Errorf("Got edit calls: %v; want: 0", m.editCalls)
	}
	if m.createCalls != 0 {
		t.Errorf("Got create calls: %v; want: 0", m.createCalls)
	}
}

// If the current IP address is different compared to the DNS record then update the record.
func TestUpdateIfIPHasChanged(t *testing.T) {
	envSetup(t)
	m := NewMockPorkbunAPI()
	m.setupRoutes()
	defer m.svr.Close()
	h, err := NewPorkbunDNSHandler(m.svr.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}
	m.retrieveResponse = retrieveResponse{
		"SUCCESS", []record{
			{
				Id:      "test2",
				Content: "10.0.0.4",
			},
		},
	}

	h.Update(ip)
	if m.editCalls != 1 {
		t.Errorf("Got edit calls: %v; want: 1", m.editCalls)
	}
	if m.createCalls != 0 {
		t.Errorf("Got create calls: %v; want: 0", m.createCalls)
	}
}

// If there is no existing DNS record then create a new one.
func TestCreateIfNoRecord(t *testing.T) {
	envSetup(t)
	m := NewMockPorkbunAPI()
	m.setupRoutes()
	defer m.svr.Close()
	h, err := NewPorkbunDNSHandler(m.svr.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}
	m.retrieveResponse = retrieveResponse{}

	h.Update(ip)
	if m.editCalls != 0 {
		t.Errorf("Got edit calls: %v; want: 0", m.editCalls)
	}
	if m.createCalls != 1 {
		t.Errorf("Got create calls: %v; want: 1", m.createCalls)
	}
}

func envSetup(t *testing.T) {
	t.Setenv("DOMAIN", "test.com")
	t.Setenv("RECORD_TYPE", "A")
	t.Setenv("RECORD_NAME", "test")
	t.Setenv("PORKBUN_API_KEY", "pk1_xxx")
	t.Setenv("PORKBUN_SECRET_KEY", "sk1_xxx")
}

type MockPorkbunAPI struct {
	svr                                   *httptest.Server
	retrieveResponse                      retrieveResponse
	retrieveCalls, editCalls, createCalls int
}

func NewMockPorkbunAPI() *MockPorkbunAPI {
	var retrieveResponse = retrieveResponse{
		"SUCCESS", []record{
			{
				Id:      "test1",
				Content: "10.0.0.2",
			}, {
				Id:      "test2",
				Content: "10.0.0.3",
			},
		},
	}

	p := &MockPorkbunAPI{retrieveResponse: retrieveResponse}

	return p
}

func (m *MockPorkbunAPI) setupRoutes() {
	mux := http.NewServeMux()
	svr := httptest.NewServer(mux)
	mux.HandleFunc(retrieveEndpoint+"/*", func(w http.ResponseWriter, r *http.Request) {
		m.retrieveCalls++
		j, _ := json.Marshal(m.retrieveResponse)
		fmt.Fprintf(w, string(j))
	})
	mux.HandleFunc(editEndpoint+"/*", func(w http.ResponseWriter, r *http.Request) {
		m.editCalls++
	})
	mux.HandleFunc(createEndpoint+"/*", func(w http.ResponseWriter, r *http.Request) {
		m.createCalls++
	})
	m.svr = svr
}
