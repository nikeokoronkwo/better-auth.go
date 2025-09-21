package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AuthClient interface {
	// retrieves the session of the user given the headers
	GetSession(headers http.Header) (*Session, error)

	// signs in a user via email (requires session cookies)
	SignInEmail(body SignInEmailOptions, headers http.Header) (AuthResult[SignInResponse], error)

	// signs in a user via username (requires session cookies)
	SignInUsername(body SignInEmailOptions, headers http.Header) (AuthResult[SignInResponse], error)

	// signs up a user via email
	SignUpEmail(body SignUpEmailOptions) (AuthResult[SignUpResponse], error)

	// signs user out (requires session cookies)
	SignOut(headers http.Header, fetchOptions struct {
		OnSuccess func()
	}) (AuthResult[bool], error)

	ChangeEmail(body struct {
		NewEmail string `json:"new_email"`
		CallbackURL url.URL `json:"callback_url"` // optional
	}) (AuthResult[bool], error)

	// request a password reset for a given user
	RequestPasswordReset(body RequestPasswordResetOptions) (AuthResult[bool], error)

	// actually reset the password
	ResetPassword(body struct {
		NewPassword string `json:"new_password"`
		// gotten from query params
		Token string `json:"token"`
	}) (AuthResult[bool], error)

	// change a user's password
	ChangePassword(body struct {
		NewPassword         string `json:"new_password"`
		CurrentPassword     string `json:"current_password"`
		RevokeOtherSessions bool   `json:"revoke_other_sessions"` // optional
	}) (AuthResult[bool], error)

	// verify a user that signed in via email
	VerifyEmail(query struct {
		Token string `json:"token"`
	})
}

type SignUpResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type SignInResponse struct {
	SignUpResponse `json:"sign_up_response"`
	Redirect       bool    `json:"redirect"` // optional
	Url            url.URL `json:"url"`      // optional
}

type RequestPasswordResetOptions struct {
	Email      string  `json:"email"`
	RedirectTo url.URL `json:"redirect_to"` // optional
}

type SignUpEmailOptions struct {
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`        // optional
	Username    string  `json:"username"`     // optional
	RememberMe  bool    `json:"remember_me"`  // optional = true
	CallbackURL url.URL `json:"callback_url"` // optional
}

type SignInEmailOptions struct {
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	RememberMe  bool    `json:"remember_me"`  // optional = true
	CallbackURL url.URL `json:"callback_url"` // optional
}

type AuthResult[T any] struct {
	Body       T
	Headers    http.Header
	StatusCode int
}

type AuthError error

func (a *AuthResult[T]) AsResponse() (http.Response, error) {
	result, err := json.Marshal(a.Body)
	if err != nil {
		return http.Response{}, err
	}

	var body io.ReadCloser
	if len(result) == 0 {
		body = http.NoBody
	} else {
		body = io.NopCloser(bytes.NewBuffer(result))
	}

	return http.Response{
		Status: fmt.Sprintf("%d %s", 200, "OK"),
		Header: a.Headers,
		Body:   body,
	}, nil
}

func (a *AuthResult[T]) WriteResponse(w http.ResponseWriter) {
	// write headers
	a.Headers.Write(w)
	json.NewEncoder(w).Encode(a.Body)
}
