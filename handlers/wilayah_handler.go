package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/evermos-rakamin/services"
	"github.com/rdsarjito/evermos-rakamin/utils"
)

// WilayahHandler handles wilayah HTTP requests
type WilayahHandler struct {
	wilayahService *services.WilayahService
}

// NewWilayahHandler creates a new wilayah handler
func NewWilayahHandler() *WilayahHandler {
	return &WilayahHandler{
		wilayahService: services.NewWilayahService(),
	}
}

// GetProvinces handles GET /locations/provinces
// @Summary Get all provinces
// @Description Get list of all provinces in Indonesia
// @Tags locations
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /locations/provinces [get]
func (h *WilayahHandler) GetProvinces(c *fiber.Ctx) error {
	provinces, err := h.wilayahService.GetProvinces()
	if err != nil {
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, provinces, "provinces retrieved successfully", fiber.StatusOK)
}

// GetRegencies handles GET /locations/regencies
// @Summary Get regencies by province ID
// @Description Get list of regencies/kabupaten by province ID
// @Tags locations
// @Accept json
// @Produce json
// @Param province_id query string true "Province ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /locations/regencies [get]
func (h *WilayahHandler) GetRegencies(c *fiber.Ctx) error {
	provinceID := c.Query("province_id")
	if provinceID == "" {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "province_id is required"), fiber.StatusBadRequest)
	}

	regencies, err := h.wilayahService.GetRegencies(provinceID)
	if err != nil {
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, regencies, "regencies retrieved successfully", fiber.StatusOK)
}

// GetDistricts handles GET /locations/districts
// @Summary Get districts by regency ID
// @Description Get list of districts/kecamatan by regency ID
// @Tags locations
// @Accept json
// @Produce json
// @Param regency_id query string true "Regency ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /locations/districts [get]
func (h *WilayahHandler) GetDistricts(c *fiber.Ctx) error {
	regencyID := c.Query("regency_id")
	if regencyID == "" {
		return utils.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "regency_id is required"), fiber.StatusBadRequest)
	}

	districts, err := h.wilayahService.GetDistricts(regencyID)
	if err != nil {
		return utils.ErrorResponse(c, err, fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, districts, "districts retrieved successfully", fiber.StatusOK)
}

