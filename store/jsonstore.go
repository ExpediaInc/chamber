package store

import (
	"encoding/json"
	"io/ioutil"
)

// ensure JSONStore conforms to Store interface
var _ Store = &JSONStore{}

type JSONStore struct {
	services map[string]service
}

type service struct {
	secrets map[string]Secret
}

func NewJSONStore(filePath string) *JSONStore {

	raw, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	var jsonData map[string]interface{}
	json.Unmarshal(raw, &jsonData)

	services := make(map[string]service)
	for k, v := range jsonData {
		services[k] = extractService(v)
	}

	return &JSONStore{
		services: services,
	}
}

func (s *JSONStore) Read(id SecretId, version int) (Secret, error) {
	if service, ok := s.services[id.Service]; ok {
		if value, ok := service.secrets[id.Key]; ok {
			return value, nil
		}
	}

	return Secret{}, ErrSecretNotFound
}

func (s *JSONStore) List(service string, includeValues bool) ([]Secret, error) {
	if service, ok := s.services[service]; ok {
		secrets := make([]Secret, len(service.secrets))
		for _, s := range service.secrets {
			secrets = append(secrets, s)
		}

		return secrets, nil
	}

	return []Secret{}, nil
}

func (s *JSONStore) History(id SecretId) ([]ChangeEvent, error) {
	// History is not supported by JSON Store
	return []ChangeEvent{}, nil
}

func (s *JSONStore) Delete(id SecretId) error {
	// Delete is not supported by JSON Store
	return nil
}

func (s *JSONStore) Write(id SecretId, value string) error {
	// Write is not supported by JSON Store
	return nil
}

func extractService(json interface{}) service {
	jsonData := json.(map[string]interface{})
	secrets := make(map[string]Secret)

	for k, v := range jsonData {
		secrets[k] = extractSecret(k, v)
	}

	return service{
		secrets: secrets,
	}
}

func extractSecret(key string, json interface{}) Secret {
	return Secret{
		Value: json.(*string),
		Meta: SecretMetadata{
			Version: 1,
			Key:     key,
		},
	}
}
