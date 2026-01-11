package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SunatResponse struct {
	RazonSocial      string `json:"razon_social"`
	NumeroDocumento  string `json:"numero_documento"`
	Estado           string `json:"estado"`
	Condicion        string `json:"condicion"`
	Direccion        string `json:"direccion"`
	Ubigeo           string `json:"ubigeo"`
	Distrito         string `json:"distrito"`
	Provincia        string `json:"provincia"`
	Departamento     string `json:"departamento"`
	EsAgenteRetencion bool   `json:"es_agente_retencion"`
}

func FetchFromSunat(ruc string) (*SunatResponse, error) {
	url := os.Getenv("SUNAT_URL") + ruc
	token := os.Getenv("RENIEC_TOKEN")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("RUC no encontrado en SUNAT")
	}
	defer resp.Body.Close()

	var result SunatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}