package data

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type VignetteProvider interface {
	GetAuthToken() (model.VignetteAuth, error)
	GetVignetteStatus(plate string, auth *model.VignetteAuth) (model.VignetteResult, error)
}

const (
	eDalniceIndex      = "https://edalnice.cz/"
	eDalniceAuth       = "https://auth.edalnice.cz/auth/connect/token"
	eDalniceValidation = "https://eshop.edalnice.cz/api/v3/charge_registrations/3906ba89-153c-4038-8e36-0ca1deb76076/"
	eDalniceTimeLayout = "2006-01-02T15:04:05.000-07:00"
)

type edalniceProvider struct{}

func NewEDalniceVignetteProvider() VignetteProvider {
	return &edalniceProvider{}
}

func getClientCredentials() (string, string, error) {
	res, err := http.Get(eDalniceIndex)
	if err != nil {
		return "", "", fmt.Errorf("error creating client credentails request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("index request failed with status: %s", res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading index response body: %w", err)
	}
	textContent := string(bodyBytes)

	clientIDMatch := regexp.MustCompile(`"REACT_APP_CLIENT_ID":"([^"]*)"`).FindStringSubmatch(textContent)
	clientSecretMatch := regexp.MustCompile(`"REACT_APP_CLIENT_SECRET":"([^"]*)"`).FindStringSubmatch(textContent)

	if len(clientIDMatch) > 1 && len(clientSecretMatch) > 1 {
		clientID := clientIDMatch[1]
		clientSecret := clientSecretMatch[1]
		return clientID, clientSecret, nil
	}

	return "", "", fmt.Errorf("could not find CLIENT_ID or CLIENT_SECRET in the index page content")
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func getAuthToken(clientId, clientSecret string) (authResponse, error) {
	token := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("Scope", "eshop.api eshoppayment.api")

	req, err := http.NewRequest(http.MethodPost, eDalniceAuth, strings.NewReader(data.Encode()))
	if err != nil {
		return authResponse{}, fmt.Errorf("error creating auth request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return authResponse{}, fmt.Errorf("error during authentication: %w", err)
	}
	defer resp.Body.Close()

	var response authResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return authResponse{}, fmt.Errorf("could not decode JSON response from authentication endpoint: %w", err)
	}

	if response.AccessToken == "" {
		return authResponse{}, fmt.Errorf("token not found in authentication response")
	}

	if response.ExpiresIn == 0 {
		return authResponse{}, fmt.Errorf("expiration not found in authentication response")
	}

	return response, nil
}

type statusResponse struct {
	Vehicle struct {
		LicensePlate string `json:"licensePlate"`
	} `json:"vehicle"`
	Charges []struct {
		ValidSince string `json:"validSince"`
		ValidUntil string `json:"validUntil"`
	} `json:"charges"`
}

func getVignetteStatus(token, plate string) (statusResponse, error) {
	req, err := http.NewRequest(http.MethodGet, eDalniceValidation+plate, nil)
	if err != nil {
		return statusResponse{}, fmt.Errorf("error creating status request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statusResponse{}, fmt.Errorf("error during target request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return statusResponse{}, fmt.Errorf("target request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var data statusResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return statusResponse{}, fmt.Errorf("could not decode JSON response from target endpoint: %w", err)
	}

	if data.Vehicle.LicensePlate == "" {
		return statusResponse{}, fmt.Errorf("license plate not found in status response")
	}

	return data, nil
}

func (_ *edalniceProvider) GetAuthToken() (model.VignetteAuth, error) {
	clientId, clientSecret, err := getClientCredentials()
	if err != nil {
		return model.VignetteAuth{}, err
	}

	response, err := getAuthToken(clientId, clientSecret)
	if err != nil {
		return model.VignetteAuth{}, err
	}

	return model.VignetteAuth{
		Token:      response.AccessToken,
		Expiration: time.Now().Add(time.Second * time.Duration(response.ExpiresIn)),
	}, nil
}

func (_ *edalniceProvider) GetVignetteStatus(plate string, auth *model.VignetteAuth) (model.VignetteResult, error) {
	status, err := getVignetteStatus(auth.Token, plate)
	if err != nil {
		return model.VignetteResult{}, err
	}

	result := model.VignetteResult{
		Plate:   plate,
		Charges: make([]model.VignetteCharge, len(status.Charges)),
	}

	for i, charge := range status.Charges {
		since, _ := time.Parse(eDalniceTimeLayout, charge.ValidSince)
		until, _ := time.Parse(eDalniceTimeLayout, charge.ValidUntil)
		result.Charges[i] = model.VignetteCharge{ValidSince: since, ValidUntil: until}
	}

	return result, nil
}
