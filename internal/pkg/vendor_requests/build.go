package vendor_requests

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

//go:embed config/vendor-data.json
var vendorDataFile []byte

//go:embed config/cookies
var cookies embed.FS

type VendorData []struct {
	AddToCartURL   string            `json:"addToCartURL"`
	CheckoutURL    string            `json:"checkoutURL"`
	Company        string            `json:"company"`
	Verb           string            `json:"verb"`
	Origin         string            `json:"origin"`
	CookieFilePath string            `json:"cookieFilePath"`
	Delay          string            `json:"delay"`
	Payload        interface{}       `json:"payload,omitempty"`
	Headers        map[string]string `json:"headers"`
}

func BuildRequestFromConfig() (*[]http.Request, error) {
	var data VendorData

	err := readDataFromString(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from file: %w", err)
	}

	var requests []http.Request
	for i := 0; i < len(data); i++ {
		parsedURL, _ := url.ParseRequestURI(data[i].AddToCartURL)

		payloadBytes, err := json.Marshal(data[i].Payload)
		if err != nil {
			fmt.Printf("Marshal failed: %v", err)
		}

		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest(data[i].Verb, parsedURL.String(), body)
		if err != nil {
			return nil, err
		}

		req.Header = make(http.Header)
		req.Host = data[i].Origin

		err = setCookie(req, data[i].CookieFilePath)
		if err != nil {
			return nil, err
		}

		for header, value := range data[i].Headers {
			req.Header.Set(header, value)
		}

		requests = append(requests, *req)
	}

	return &requests, nil
}

func setCookie(req *http.Request, cookieFilePath string) error {
	cookieBytes, err := cookies.ReadFile(cookieFilePath)
	if err != nil {
		return fmt.Errorf("failed to get cookie: %w", err)
	}

	req.Header.Set("Cookie", os.ExpandEnv(string(cookieBytes)))
	return nil
}

func readDataFromString(data *VendorData) error {
	err := json.Unmarshal(vendorDataFile, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	return nil
}
