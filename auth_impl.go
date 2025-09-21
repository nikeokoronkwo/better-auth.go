package main

import (
	"database/sql"
	"net/http"
)

type AuthClientImpl struct {
	db *sql.DB
	schema map[string]Model
	EmailAndPassword EmailAndPasswordBaseOptions
	Username UsernameOptions
	SecondaryStorage SecondaryStorage
}

func (a AuthClientImpl) SignUpEmail(body SignUpEmailOptions) (AuthResult[SignUpResponse], error) {
	return AuthResult[SignUpResponse]{}, nil
}

type AuthClientOptions struct {
	Database *sql.DB
	UserModelName string // optional
	SessionModelName string // optional
	EmailAndPassword struct { // optional
		EmailAndPasswordBaseOptions
		Enable bool
	}
	Username UsernameOptions // optional
	SecondaryStorage SecondaryStorage // optional
}

type UsernameOptions struct { 
	Enable bool
	MinUsernameLength int // optional = 3
	MaxUsernameLength int // optional = 30
	UsernameValidator func (string) bool
	UsernameNormalization func (string) string
}

type EmailAndPasswordBaseOptions struct {
	MinPasswordLength int // optional = 8
	MaxPasswordLength int // optional = 128
	ResetPasswordTokenExpiresIn int // optional = 1 hour (3600)
	Password struct {
		Hash func(string) (string, error)
		Verify func(password, hash string) (bool, error)
	}
	SendVerificationEmail func(user string, url string, token string, request http.Request)
	SendResetPassword func(user string, url string, token string, request http.Request)
	OnPasswordChange func(user string, request http.Request)
}

func InitialiseClient(options AuthClientOptions) AuthClientImpl {
	// set up models

	// initialise database
	if options.Database == nil {
		panic("database cannot be nil")
	}

	var schemas map[string]Model
	schemas["user"] = Model{
		Fields: map[string]Field {
			"id": Field{
				Type: String,
				Primary: true,
				NotNull: true,
			},
			"name": Field{
				Type: String,
				NotNull: true,
			},
			"email": Field{
				Type: String,
				NotNull: true,
				Unique: true,
			},
			"email_verified": Field{
				Type: Boolean,
				NotNull: true,
				Default: false,
			},
			"created_at": Field{
				Type: Timestamptz,
				NotNull: true,
				DefaultExpression: "NOW()",
			},
			"updated_at": Field{
				Type: Timestamptz,
				NotNull: true,
				DefaultExpression: "NOW()",
			},
		},
	}
	if options.Username.Enable {
		schemas["user"].Fields["username"] = Field{
			Type: String,
			NotNull: true,
			Unique: true,
		}
	}
	schemas["session"] = Model{
			Fields: map[string]Field {
				"id": Field{
					Type: String,
					Primary: true,
					NotNull: true,
				},
				"user_id": Field{
					Type: String,
					NotNull: true,
					References: func() (Model, string) {
						return schemas["user"], "id"
					},
				},
				"token": Field{
					Type: String,
					NotNull: true,
				},
				"created_at": Field{
					Type: Timestamptz,
					NotNull: true,
					DefaultExpression: "NOW()",
				},
				"expires_at": Field{
					Type: Timestamptz,
					NotNull: true,
				},
				"updated_at": Field{
					Type: Timestamptz,
					NotNull: true,
					DefaultExpression: "NOW()",
				},
				"ip_address": Field{
					Type: String,
				},
				"user_agent": Field{
					Type: String,
				},
			},
		}

	// set up migration
	var emailAndPasswordOpts EmailAndPasswordBaseOptions
	if options.EmailAndPassword.Enable {
		if options.EmailAndPassword.Password.Hash != nil {
			emailAndPasswordOpts.Password.Hash = options.EmailAndPassword.Password.Hash
		} else {
			emailAndPasswordOpts.Password.Hash = HashPassword
		}
		if options.EmailAndPassword.Password.Verify != nil {
			emailAndPasswordOpts.Password.Verify = options.EmailAndPassword.Password.Verify
		} else {
			emailAndPasswordOpts.Password.Verify = VerifyPassword
		}
	}


	// return client
	return AuthClientImpl{
		db: options.Database,
		EmailAndPassword: options.EmailAndPassword.EmailAndPasswordBaseOptions,
		Username: options.Username,
		SecondaryStorage: options.SecondaryStorage,
		schema: schemas,
	}
}