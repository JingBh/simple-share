package controllers

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jingbh/simple-share/app/context"
	_oidc "github.com/jingbh/simple-share/internal/oidc"
	"github.com/jingbh/simple-share/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type userInfoResponse struct {
	Subject  string `json:"subject"`
	Username string `json:"username"`
}

func setCallbackCookie(c echo.Context, name, value string) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   true,
		HttpOnly: true,
	})
}

func removeCallbackCookie(c echo.Context, name string) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	})
}

func AuthRedirectLogin(c echo.Context) error {
	nonce, err := _oidc.GenerateNonce()
	if err != nil {
		return err
	}

	state, err := _oidc.GenerateNonce()
	if err != nil {
		return err
	}

	url := _oidc.OAuthConfig.AuthCodeURL(state, oidc.Nonce(nonce))

	setCallbackCookie(c, "_nonce", nonce)
	setCallbackCookie(c, "_state", state)
	return c.Redirect(http.StatusFound, url)
}

func AuthCallback(c echo.Context) error {
	state, err := c.Cookie("_state")
	if err != nil || c.QueryParam("state") != state.Value {
		return echo.NewHTTPError(http.StatusForbidden, "invalid state")
	}

	token, err := _oidc.OAuthConfig.Exchange(c.Request().Context(), c.QueryParam("code"))
	if err != nil {
		return err
	}

	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "invalid token")
	}

	verifier := _oidc.Provider.Verifier(_oidc.Config)
	idToken, err := verifier.Verify(c.Request().Context(), rawIdToken)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusForbidden,
			Message:  "invalid token",
			Internal: err,
		}
	}

	nonce, err := c.Cookie("_nonce")
	if err != nil || idToken.Nonce != nonce.Value {
		return echo.NewHTTPError(http.StatusForbidden, "invalid nonce")
	}

	removeCallbackCookie(c, "_nonce")
	removeCallbackCookie(c, "_state")
	c.SetCookie(&http.Cookie{
		Name:     "_token",
		Value:    rawIdToken,
		Path:     "/",
		Expires:  idToken.Expiry,
		Secure:   true,
		HttpOnly: true,
	})

	return c.Redirect(http.StatusFound, utils.Url("/"))
}

func AuthGetUserinfo(c echo.Context) error {
	cc := c.(context.CustomContext)

	if cc.Token == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	return c.JSON(http.StatusOK, &userInfoResponse{
		Subject:  cc.Token.Subject,
		Username: cc.Username,
	})
}
