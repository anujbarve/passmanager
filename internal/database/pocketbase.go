// internal/database/pocketbase.go
package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"passmanager/internal/models"
	"strings"
	"time"
)

type PocketBaseClient struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

type AuthResponse struct {
	Token  string `json:"token"`
	Record struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"record"`
}

type ListResponse struct {
	Items      []models.Credential `json:"items"`
	TotalItems int                 `json:"totalItems"`
}

type ConfigListResponse struct {
	Items []models.VaultConfig `json:"items"`
}

func NewPocketBaseClient(baseURL string) *PocketBaseClient {
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &PocketBaseClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *PocketBaseClient) TestConnection() error {
	resp, err := p.httpClient.Get(p.baseURL + "/api/health")
	if err != nil {
		return fmt.Errorf("cannot reach PocketBase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PocketBase health check failed: status %d", resp.StatusCode)
	}
	return nil
}

func (p *PocketBaseClient) Authenticate(email, password string) error {
	authData := map[string]string{
		"identity": email,
		"password": password,
	}

	jsonData, err := json.Marshal(authData)
	if err != nil {
		return err
	}

	endpoints := []string{
		"/api/collections/_superusers/auth-with-password",
		"/api/collections/users/auth-with-password",
		"/api/admins/auth-with-password",
	}

	var lastErr error
	for _, endpoint := range endpoints {
		resp, err := p.httpClient.Post(
			p.baseURL+endpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			lastErr = err
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var authResp AuthResponse
			if err := json.Unmarshal(body, &authResp); err != nil {
				lastErr = err
				continue
			}
			p.authToken = authResp.Token
			return nil
		}

		lastErr = fmt.Errorf("endpoint %s failed: %s", endpoint, string(body))
	}

	return fmt.Errorf("authentication failed: %v", lastErr)
}

func (p *PocketBaseClient) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", p.baseURL, endpoint), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if p.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+p.authToken)
	}

	return p.httpClient.Do(req)
}

func (p *PocketBaseClient) CreateCredential(cred models.Credential) (*models.Credential, error) {
	resp, err := p.doRequest("POST", "/api/collections/credentials/records", cred)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create credential: %s", string(body))
	}

	var created models.Credential
	if err := json.Unmarshal(body, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

func (p *PocketBaseClient) GetCredential(id string) (*models.Credential, error) {
	resp, err := p.doRequest("GET", fmt.Sprintf("/api/collections/credentials/records/%s", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("credential not found")
	}

	var cred models.Credential
	if err := json.NewDecoder(resp.Body).Decode(&cred); err != nil {
		return nil, err
	}

	return &cred, nil
}

func (p *PocketBaseClient) ListCredentials(search string) ([]models.Credential, error) {
	endpoint := "/api/collections/credentials/records?perPage=500&sort=-created"
	if search != "" {
		endpoint += "&filter=" + url.QueryEscape(fmt.Sprintf(
			"title~'%s' || username~'%s' || url~'%s' || category~'%s'",
			search, search, search, search))
	}

	resp, err := p.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var listResp ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, err
	}

	return listResp.Items, nil
}

func (p *PocketBaseClient) UpdateCredential(id string, cred models.Credential) (*models.Credential, error) {
	resp, err := p.doRequest("PATCH", fmt.Sprintf("/api/collections/credentials/records/%s", id), cred)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update credential: %s", string(body))
	}

	var updated models.Credential
	if err := json.Unmarshal(body, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (p *PocketBaseClient) DeleteCredential(id string) error {
	resp, err := p.doRequest("DELETE", fmt.Sprintf("/api/collections/credentials/records/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete credential")
	}

	return nil
}

func (p *PocketBaseClient) SaveVaultConfig(config models.VaultConfig) error {
	resp, err := p.doRequest("POST", "/api/collections/vault_config/records", config)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to save config: %s", string(body))
	}

	return nil
}

func (p *PocketBaseClient) GetVaultConfig() (*models.VaultConfig, error) {
	resp, err := p.doRequest("GET", "/api/collections/vault_config/records?perPage=1", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var listResp ConfigListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, err
	}

	if len(listResp.Items) == 0 {
		return nil, fmt.Errorf("vault not initialized")
	}

	return &listResp.Items[0], nil
}

func (p *PocketBaseClient) UpdateVaultConfig(id string, config models.VaultConfig) error {
	resp, err := p.doRequest("PATCH", fmt.Sprintf("/api/collections/vault_config/records/%s", id), config)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update config: %s", string(body))
	}

	return nil
}

func (p *PocketBaseClient) GetCredentialCount() (int, error) {
	resp, err := p.doRequest("GET", "/api/collections/credentials/records?perPage=1", nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var listResp ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return 0, err
	}

	return listResp.TotalItems, nil
}