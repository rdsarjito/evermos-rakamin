package dto

// ProvinceResponse represents province response
type ProvinceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RegencyResponse represents regency response
type RegencyResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProvinceID string `json:"province_id"`
}

// DistrictResponse represents district response
type DistrictResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	RegencyID  string `json:"regency_id"`
}

