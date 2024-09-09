package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/internal/services/referral"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ApiReferralController struct {
	ReferralService referral.ReferralService
	//middleware  *middleware.ApiController
	logger    *zerolog.Logger
	validator *validator.Validate
}

func NewApiReferralController(userService *referral.ReferralService, logger *zerolog.Logger, validator *validator.Validate) *ApiReferralController {
	return &ApiReferralController{
		ReferralService: *userService,
		logger:          logger,
		validator:       validator,
	}
}

func (arc *ApiReferralController) CreateReferral(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = arc.logger.WithContext(ctx)

	arc.logger.Debug().Msg("starting the handler 'CreateReferral'")

	id := c.Get("id")

	userId, ok := id.(string)
	if ok != true {
		arc.logger.Debug().Msgf("invalid id: %v", id)
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		arc.logger.Debug().Msgf("invalid id: %v", id)
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	body := &models.Request{}
	body.UserId = userIdInt

	// request parsing
	if err := c.Bind(body); err != nil {
		arc.logger.Warn().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	result, err := arc.ReferralService.CreateReferral(ctx, body)
	if err != nil {
		arc.logger.Debug().Msgf("error receiving user data: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (arc *ApiReferralController) GetReferral(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = arc.logger.WithContext(ctx)

	arc.logger.Debug().Msg("starting the handler 'GetReferral'")

	url := c.Param("url")

	result, err := arc.ReferralService.GetReferral(ctx, url)
	if err != nil {
		arc.logger.Debug().Msgf("error receiving a short referral link: %v", err)
		return c.JSON(http.StatusNotFound, err.Error())
	}

	//redirection to the registration page
	return c.Redirect(http.StatusMovedPermanently, result)
}
