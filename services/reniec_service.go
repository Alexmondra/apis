package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DecolectaResponse struct {
	FirstName      string `json:"first_name"`
	FirstLastName  string `json:"first_last_name"`
	SecondLastName string `json:"second_last_name"`
	DocumentNumber string `json:"document_number"`
}

func FetchFromReniec(dni string) (*DecolectaResponse, error) {
	url := os.Getenv("RENIEC_URL") + dni
	token := os.Getenv("RENIEC_TOKEN")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("DNI no encontrado en RENIEC")
	}
	defer resp.Body.Close()

	var result DecolectaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}