package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/michael-a-shelton/vendor-scraper/internal/pkg/vendor_requests"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Calling API...")

	client := &http.Client{}

	requests, err := vendor_requests.BuildRequestFromConfig()
	if err != nil {
		fmt.Printf("build requests failed %v",errors.Unwrap(err))
	}

	for _, req := range *requests {

		resp, err := client.Do(&req)
		if err != nil {
			fmt.Printf("Request failed: %+v\n", err)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed Request: %+v\nResponse: %+v\nBody:%+v\n", req, resp, req.Body)
		}

		defer func() {
			err = resp.Body.Close()
			if err != nil {
				fmt.Printf("error closing body: %+v", err)
			}
		}()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}

		var responseObject http.Response
		_ = json.Unmarshal(bodyBytes, &responseObject)
		fmt.Printf("API Response as struct %+v\n", responseObject)
	}

}
