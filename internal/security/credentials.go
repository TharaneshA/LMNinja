package security

import "github.com/zalando/go-keyring"

const keyRingService = "LMNinja"

// CredentialManager handles secure storage of API keys in the OS keychain.
type CredentialManager struct{}

func NewCredentialManager() *CredentialManager {
	return &CredentialManager{}
}

func (cm *CredentialManager) StoreAPIKey(connectionID string, apiKey string) error {
	return keyring.Set(keyRingService, connectionID, apiKey)
}

func (cm *CredentialManager) GetAPIKey(connectionID string) (string, error) {
	return keyring.Get(keyRingService, connectionID)
}

func (cm *CredentialManager) DeleteAPIKey(connectionID string) error {
	return keyring.Delete(keyRingService, connectionID)
}