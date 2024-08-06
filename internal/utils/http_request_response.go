package utils

import (
	"encoding/json"
	"net/http"
)

func ReadJson[T any](request *http.Request) (*T, error) {
	var result T
	err := json.NewDecoder(request.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func WriteJson[T any](response http.ResponseWriter, data T) error {
	response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(response).Encode(data)
}
