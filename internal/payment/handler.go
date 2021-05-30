package payment

import (
	"net/http"

	"github.com/labstack/echo"
)

type Handler interface {
	Payment(c echo.Context) error
	Inquiry(c echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return handler{service}
}

func (h handler) Payment(c echo.Context) error {
	req := new(ReqPayment)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.service.Payment(*req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

func (h handler) Inquiry(c echo.Context) error {
	req := new(ReqInquiry)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := h.service.Inquiry(req.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
