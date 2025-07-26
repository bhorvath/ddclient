package dns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/netip"

	"github.com/bhorvath/ddclient/config"
)

const (
	retrieveEndpoint = "/api/json/v3/dns/retrieveByNameType"
	editEndpoint     = "/api/json/v3/dns/editByNameType"
	createEndpoint   = "/api/json/v3/dns/create"
)

type PorkbunDNSHandler struct {
	baseURL string
	config  *config.App
}

type retrieveRequest struct {
	APIKey       string `json:"apikey"`
	SecretAPIKey string `json:"secretapikey"`
}

type retrieveResponse struct {
	Status  string   `json:"status"`
	Records []record `json:"records"`
}

type record struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl"`
	Prio    string `json:"prio"`
	Notes   string `json:"notes"`
}

type editRequest struct {
	APIKey       string `json:"apikey"`
	SecretAPIKey string `json:"secretapikey"`
	Content      string `json:"content"`
}

type createRequest struct {
	APIKey       string `json:"apikey"`
	SecretAPIKey string `json:"secretapikey"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Content      string `json:"content"`
	TTL          string `json:"ttl"`
}

// NewPorkbunDNSHandler allows a DNS record in Porkbun to be read, updated or created.
func NewPorkbunDNSHandler(baseURL string, config *config.App) (*PorkbunDNSHandler, error) {
	return &PorkbunDNSHandler{
		baseURL: baseURL,
		config:  config,
	}, nil
}

// Update either creates or updates a record based on the current IP address. If the current address
// is the same as the record then no change is made. Update does not currently support making changes
// to multiple records, so an error is thrown if multiple records exist.
func (h *PorkbunDNSHandler) Update(IP netip.Addr) error {
	fmt.Print("Checking whether record exists... ")
	r, err := h.retrieveRecords()
	if err != nil {
		return err
	}

	c := len(r)
	fmt.Printf("Found %v existing record(s).\n", c)
	if c > 1 {
		return errors.New("more than one record to update found")
	} else if c == 1 {
		// Porkbun doesn't gracefully handle update requests if there is no change to the record and
		// let's also avoid an unnecessary network request. Therefore only update if there is a genuine
		// change in IP.
		curIP, err := netip.ParseAddr(r[0].Content)
		if err != nil {
			return err
		}
		if !compareIPs(curIP, IP) {
			fmt.Print("IP has changed. Updating... ")
			err = h.editRecord(IP)
			if err != nil {
				return err
			}
			fmt.Print("Done!\n")
		} else {
			fmt.Println("IP has not changed. Nothing to do.")
		}
	} else {
		// Create new record
		fmt.Print("Creating new record... ")
		err = h.createRecord(IP)
		if err != nil {
			return err
		}
		fmt.Print("Done!\n")
	}

	return nil
}

func (h *PorkbunDNSHandler) retrieveRecords() ([]record, error) {
	body, err := json.Marshal(retrieveRequest{
		APIKey:       h.config.APIKey,
		SecretAPIKey: h.config.SecretKey,
	})
	if err != nil {
		return []record{}, err
	}
	bodyReader := bytes.NewReader(body)

	requestURL := h.baseURL + retrieveEndpoint + "/" + h.config.Domain + "/" + h.config.Type + "/" + h.config.Name
	res, err := http.Post(requestURL, "application/json", bodyReader)
	if err != nil {
		return []record{}, err
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		resBody, _ := io.ReadAll(res.Body)
		return []record{}, errors.New("failed to retrieve records; " + string(resBody))
	}

	var rr retrieveResponse
	err = json.NewDecoder(res.Body).Decode(&rr)
	if err != nil {
		return []record{}, err
	}

	return rr.Records, nil
}

func (h *PorkbunDNSHandler) editRecord(ip netip.Addr) error {
	body, err := json.Marshal(editRequest{
		APIKey:       h.config.APIKey,
		SecretAPIKey: h.config.SecretKey,
		Content:      ip.String(),
	})
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(body)

	requestURL := h.baseURL + editEndpoint + "/" + h.config.Domain + "/" + h.config.Type + "/" + h.config.Name
	res, err := http.Post(requestURL, "application/json", bodyReader)
	if err != nil {
		return err
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		resBody, _ := io.ReadAll(res.Body)
		return errors.New("failed to edit record; " + string(resBody))
	}

	return nil
}

func (h *PorkbunDNSHandler) createRecord(ip netip.Addr) error {
	body, err := json.Marshal(createRequest{
		APIKey:       h.config.APIKey,
		SecretAPIKey: h.config.SecretKey,
		Name:         h.config.Name,
		Type:         h.config.Type,
		Content:      ip.String(),
	})
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(body)

	requestURL := h.baseURL + createEndpoint + "/" + h.config.Domain
	res, err := http.Post(requestURL, "application/json", bodyReader)
	if err != nil {
		return err
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		resBody, _ := io.ReadAll(res.Body)
		return errors.New("failed to create record; " + string(resBody))
	}

	return nil
}

func compareIPs(curIP netip.Addr, newIP netip.Addr) bool {
	return curIP == newIP
}
