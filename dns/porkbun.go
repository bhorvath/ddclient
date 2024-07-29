package dns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
)

const (
	retrieveEndpoint = "/api/json/v3/dns/retrieveByNameType"
	editEndpoint     = "/api/json/v3/dns/editByNameType"
	createEndpoint   = "/api/json/v3/dns/create"
)

type PorkbunDNSHandler struct {
	baseURL    string
	domain     string
	recordType string
	recordName string
	apiKey     string
	secretKey  string
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

// NewPorkbunDNSHandler allows a DNS record in Porkbun to be read, updated or created. The following
// environment variables must be provided:
// - DOMAIN: The domain of the record; eg. 'test.com'
// - RECORD_TYPE: The type of the record to interact with; eg. 'A'
// - RECORD_NAME: The name of the record to interact with; eg. 'subdomain'
// - PORKBUN_API_KEY: The API key obtained from Porkbun; eg. 'pk1_xxx'
// - PORKBUN_SECRET_KEY: The secret key for the corresponding API key; eg. 'sk1_xxx'
func NewPorkbunDNSHandler(baseURL string) (*PorkbunDNSHandler, error) {
	domain, err := getEnvVar("DOMAIN")
	if err != nil {
		return nil, err
	}

	recordType, err := getEnvVar("RECORD_TYPE")
	if err != nil {
		return nil, err
	}

	recordName, err := getEnvVar("RECORD_NAME")
	if err != nil {
		return nil, err
	}

	apiKey, err := getEnvVar("PORKBUN_API_KEY")
	if err != nil {
		return nil, err
	}

	secretKey, err := getEnvVar("PORKBUN_SECRET_KEY")
	if err != nil {
		return nil, err
	}

	return &PorkbunDNSHandler{
		baseURL:    baseURL,
		domain:     domain,
		recordType: recordType,
		recordName: recordName,
		apiKey:     apiKey,
		secretKey:  secretKey,
	}, nil
}

func getEnvVar(envVar string) (string, error) {
	value := os.Getenv(envVar)
	if value == "" {
		return "", errors.New("environment variable " + envVar + " is missing")
	}

	return value, nil
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
		APIKey:       h.apiKey,
		SecretAPIKey: h.secretKey,
	})
	if err != nil {
		return []record{}, err
	}
	bodyReader := bytes.NewReader(body)

	requestURL := h.baseURL + retrieveEndpoint + "/" + h.domain + "/" + h.recordType + "/" + h.recordName
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
		APIKey:       h.apiKey,
		SecretAPIKey: h.secretKey,
		Content:      ip.String(),
	})
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(body)

	requestURL := h.baseURL + editEndpoint + "/" + h.domain + "/" + h.recordType + "/" + h.recordName
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
		APIKey:       h.apiKey,
		SecretAPIKey: h.secretKey,
		Name:         h.recordName,
		Type:         h.recordType,
		Content:      ip.String(),
	})
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(body)

	requestURL := h.baseURL + createEndpoint + "/" + h.domain
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
