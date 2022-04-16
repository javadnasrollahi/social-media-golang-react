package service

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"jupiter/app/common/token"
	"jupiter/app/model"
	"jupiter/config"
	"net/http"
)

type AuthService struct {
}

func (s AuthService) GenerateAccessToken(user *model.User) (*token.JwtCustomToken, error) {
	claims := &token.JwtCustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	t, err := token.Generate(token.AccessToken, claims)

	if err != nil {
		return nil, errors.New("service.authService.GenerateAccessToken")
	}

	return t, nil
}

func (s AuthService) GenerateRefreshToken(user *model.User) (*token.JwtCustomToken, error) {
	claims := &token.JwtCustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	t, err := token.Generate(token.RefreshToken, claims)

	if err != nil {
		return nil, errors.New("service.authService.GenerateAccessToken")
	}

	return t, nil
}

func (s AuthService) Response(c echo.Context, user *model.User) error {
	accessToken, atErr := s.GenerateAccessToken(user)
	if atErr != nil {
		return echo.ErrInternalServerError
	}

	refreshToken, rtErr := s.GenerateRefreshToken(user)
	if rtErr != nil {
		return echo.ErrInternalServerError
	}

	fmt.Println(refreshToken.ExpiresAt())
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.String(),
		Path:     "/",
		Expires:  refreshToken.ExpiresAt(),
		Secure:   config.GetConfig().App.Production,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": accessToken.String(),
	})
}

//func (s AuthService) User(c echo.Context) (*model.User, error) {
//	jwtUser := c.Get("user").(*jwt.Token)
//	claims := jwtUser.Claims.(*token.JwtCustomClaims)
//
//	db := app.GetDB()
//	user := &model.User{
//		ID: claims.UserID,
//	}
//
//	result := db.First(user)
//
//	if result.Error != nil {
//		return nil, errors.New("user not found")
//	}
//
//	return user, nil
//}
