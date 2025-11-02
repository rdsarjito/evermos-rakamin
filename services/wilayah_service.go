package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rdsarjito/evermos-rakamin/config"
	"github.com/rdsarjito/evermos-rakamin/domain/dto"
)

// WilayahService handles API Wilayah Indonesia integration
type WilayahService struct {
	baseURL    string
	httpClient *http.Client
}

// NewWilayahService creates a new wilayah service
func NewWilayahService() *WilayahService {
	cfg := config.Get()
	baseURL := cfg.External.APILocation
	if baseURL == "" {
		baseURL = "https://www.emsifa.com/api-wilayah-indonesia/api" // default
	}

	return &WilayahService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetProvinces gets all provinces from API Wilayah Indonesia
func (s *WilayahService) GetProvinces() ([]dto.ProvinceResponse, error) {
	url := fmt.Sprintf("%s/provinces.json", s.baseURL)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provinces: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var provinces []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(body, &provinces); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := make([]dto.ProvinceResponse, len(provinces))
	for i, p := range provinces {
		result[i] = dto.ProvinceResponse{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	return result, nil
}

// GetRegencies gets all regencies by province ID from API Wilayah Indonesia
func (s *WilayahService) GetRegencies(provinceID string) ([]dto.RegencyResponse, error) {
	url := fmt.Sprintf("%s/regencies/%s.json", s.baseURL, provinceID)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch regencies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var regencies []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		ProvinceID string `json:"province_id"`
	}

	if err := json.Unmarshal(body, &regencies); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := make([]dto.RegencyResponse, len(regencies))
	for i, r := range regencies {
		result[i] = dto.RegencyResponse{
			ID:         r.ID,
			Name:       r.Name,
			ProvinceID: r.ProvinceID,
		}
	}

	return result, nil
}

// GetDistricts gets all districts by regency ID from API Wilayah Indonesia
func (s *WilayahService) GetDistricts(regencyID string) ([]dto.DistrictResponse, error) {
	url := fmt.Sprintf("%s/districts/%s.json", s.baseURL, regencyID)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch districts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var districts []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		RegencyID string `json:"regency_id"`
	}

	if err := json.Unmarshal(body, &districts); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := make([]dto.DistrictResponse, len(districts))
	for i, d := range districts {
		result[i] = dto.DistrictResponse{
			ID:        d.ID,
			Name:      d.Name,
			RegencyID: d.RegencyID,
		}
	}

	return result, nil
}

