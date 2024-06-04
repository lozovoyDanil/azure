package endpoints

import (
	"context"
	"net/http"
	"strings"

	"filmlib/auth"
	"filmlib/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const authHeader string = "Authorization"

var validate = validator.New()

type AuthHandler struct {
	services auth.Service
}

func NewAuthHandler(services auth.Service) *AuthHandler {
	return &AuthHandler{services: services}
}

func (h *AuthHandler) InitRoutes() *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	api := router.Group("/api")

	api.GET("/healthz", func(c echo.Context) error {
		err := h.services.Healthz(context.Background())
		if err != nil {
			return c.String(http.StatusInternalServerError, "Oops, something's off :(")
		}

		return c.String(http.StatusOK, "Doing just fine!")
	})

	api.POST("/sign-up", func(c echo.Context) error {
		var user model.User

		err := c.Bind(&user)
		if err != nil {
			return newErrorResponse(http.StatusBadRequest, err.Error())
		}
		err = validate.Struct(user)
		if err != nil {
			return newErrorResponse(http.StatusBadRequest, err.Error())
		}

		id, err := h.services.CreateUser(user)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"id": id,
		})

	})

	type signInInput struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	api.POST("/sign-in", func(c echo.Context) error {
		var input signInInput

		err := c.Bind(&input)
		if err != nil {
			return newErrorResponse(http.StatusBadRequest, err.Error())
		}
		err = validate.Struct(input)
		if err != nil {
			return newErrorResponse(http.StatusBadRequest, err.Error())
		}

		token, err := h.services.GenerateToken(input.Username, input.Password)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error())

		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": token,
		})
	})

	api.GET("/identity", func(c echo.Context) error {
		identity, err := h.userIdentity(c)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, identity)
	})

	api.GET("/favorites", func(c echo.Context) error {
		identity, err := h.userIdentity(c)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error())
		}

		var movIds struct {
			Ids []primitive.ObjectID `json:"Ids"`
		}

		ids, err := h.services.UserFavorites(identity.Id)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error())
		}
		movIds.Ids = ids

		return c.JSON(http.StatusOK, movIds)
	})

	api.POST("/favorites/:movie_id", func(c echo.Context) error {
		movId := c.Param("movie_id")
		if movId == "" {
			return newErrorResponse(http.StatusBadRequest, "no movie Id provided in req")
		}

		identity, err := h.userIdentity(c)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error())
		}

		err = h.services.AddToFavorites(identity.Id, movId)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "OK",
		})
	})

	api.DELETE("/favorites/:movie_id", func(c echo.Context) error {
		movId := c.Param("movie_id")
		if movId == "" {
			return newErrorResponse(http.StatusBadRequest, "no movie Id provided in req")
		}

		identity, err := h.userIdentity(c)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error())
		}

		err = h.services.RemoveFavorite(identity.Id, movId)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "OK",
		})
	})

	return router
}

func (h *AuthHandler) userIdentity(c echo.Context) (*model.Identity, error) {
	header := c.Request().Header.Get(authHeader)
	if header == "" {
		return nil, newErrorResponse(http.StatusUnauthorized, "empty header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return nil, newErrorResponse(http.StatusUnauthorized, "wrong header type, required BEARER")
	}

	user, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		return nil, newErrorResponse(http.StatusUnauthorized, err.Error())
	}

	return &user, nil
}
