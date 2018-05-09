package rest

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Get performs a GET request on the provided URL with the VSTS Access Token defined and returns the body
func Get(url string, accessToken string) []byte {
	accessTokenBase64 := base64.StdEncoding.EncodeToString([]byte(":" + accessToken))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create request: ", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", accessTokenBase64))

	client := &http.Client{}
	log.Print("Sending URL Request to ", url)

	timeSent := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("GET Request failed: ", err)
	}
	defer resp.Body.Close()
	elapsedTime := time.Since(timeSent)
	log.Printf("Received response in %.fms Status: %s, Content Type: %s", elapsedTime.Seconds()*1000, resp.Status, resp.Header.Get("Content-Type"))

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Response code was: ", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read body, ", err)
	}
	return body
}
