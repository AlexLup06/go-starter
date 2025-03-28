package service

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"alexlupatsiy.com/personal-website/backend/domain"
	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/helpers/passwords"
	"alexlupatsiy.com/personal-website/backend/repository"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/idtoken"
)

type AppleAuthResponse struct {
	Code    string `form:"code" binding:"required"`
	IDToken string `form:"id_token" binding:"required"`
	State   string `form:"state"`
	User    string `form:"user"`
}

type LoginWithEmailRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type RequestPasswordResetRequest struct {
	Email string `form:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Password string `form:"password" binding:"required"`
	Token    string `form:"token" binding:"required"`
}

type AuthService struct {
	authStorage    repository.AuthStorage
	userService    *UserService
	tokenService   *TokenService
	googleClientId string
}

func NewAuthService(
	authStorage repository.AuthStorage,
	userService *UserService,
	tokenService *TokenService,
	googleClientId string,
) *AuthService {
	return &AuthService{
		authStorage:    authStorage,
		userService:    userService,
		tokenService:   tokenService,
		googleClientId: googleClientId,
	}
}

func (a *AuthService) LoginWithEmail(ctx context.Context, request LoginWithEmailRequest) (UserInfo, error) {
	user, err := a.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return UserInfo{}, err
	}

	emailAuth, err := a.authStorage.GetAuthProviderByUserId(ctx, user.ID, repository.METHOD_EMAIL)
	if err != nil {
		// the auth provider "email" does not exist
		return UserInfo{}, err
	}

	if !passwords.IsSamePassword(request.Password, *emailAuth.PasswordHash) {
		return UserInfo{}, customErrors.NewUnauthorizedError("invalid password")
	}

	userInfo := UserInfo{
		UserId:   user.ID,
		Username: user.Username,
		Email:    *user.Email,
	}

	return userInfo, nil
}

func (a *AuthService) UpdateUserPassword(ctx context.Context, userId string, resetPasswordRequest ResetPasswordRequest) error {
	hashedPassword, err := passwords.HashPassword(resetPasswordRequest.Password)
	if err != nil {
		return err
	}
	err = a.authStorage.UpdateUserPassword(ctx, userId, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

const appleKeysURL = "https://appleid.apple.com/auth/keys"

type ApplePublicKey struct {
	Keys []map[string]interface{} `json:"keys"`
}

// Fetch Apple’s Public Key
func GetApplePublicKey() (*rsa.PublicKey, error) {
	resp, err := http.Get(appleKeysURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var appleKeys ApplePublicKey
	json.Unmarshal(body, &appleKeys)

	// Extract RSA public key from JSON
	keyData := appleKeys.Keys[0] // Apple may have multiple keys; you should match "kid"
	nBytes, _ := base64.RawURLEncoding.DecodeString(keyData["n"].(string))
	eBytes, _ := base64.RawURLEncoding.DecodeString(keyData["e"].(string))

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(new(big.Int).SetBytes(eBytes).Int64()),
	}

	return pubKey, nil
}

// Verify the ID Token
func (a *AuthService) VerifyAppleIDToken(idToken string) (jwt.MapClaims, error) {
	publicKey, err := GetApplePublicKey()
	if err != nil {
		return nil, errors.New("failed to fetch Apple public key")
	}

	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, errors.New("invalid ID token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

const appleTokenURL = "https://appleid.apple.com/auth/token"

type AppleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

func (a *AuthService) ExchangeAppleCodeForTokens(code string) (*AppleTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", "your_client_id")
	data.Set("client_secret", "your_generated_jwt")
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "your_redirect_uri")

	req, _ := http.NewRequest("POST", appleTokenURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse AppleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	if tokenResponse.AccessToken == "" {
		return nil, errors.New("failed to get Apple tokens")
	}

	return &tokenResponse, nil
}

func (a *AuthService) ValidateGoogleIdToken(ctx context.Context, idToken string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, a.googleClientId)
	if err != nil {
		return &idtoken.Payload{}, err
	}
	return payload, err
}

func (a *AuthService) GoogleLogin(ctx context.Context, payload *idtoken.Payload) (UserInfo, error) {

	// check if user already exists with email or sub
	email := payload.Claims["email"].(string)
	givenName := payload.Claims["given_name"].(string)
	sub := payload.Subject
	emailExists := false

	user, err := a.userService.GetUserByEmail(ctx, email)
	if err != nil && err != customErrors.ErrUserDoesNotExist {
		return UserInfo{}, err
	}
	if err != customErrors.ErrUserDoesNotExist {
		emailExists = true
	}

	if emailExists {
		userInfo := UserInfo{
			Email:    *user.Email,
			UserId:   user.ID,
			Username: givenName,
		}

		_, err := a.authStorage.GetAuthProviderByUserId(ctx, user.ID, repository.METHOD_GOOGLE)
		if err != nil && err != customErrors.ErrAuthProviderDoesNotExist {
			return UserInfo{}, err
		}

		// if email and sub exist: login user
		if err != customErrors.ErrAuthProviderDoesNotExist {
			return userInfo, nil
		}

		// if email exists and sub does not: link google auth provider to that account
		googleAuthProvider := domain.AuthProvider{
			UserID:         user.ID,
			Method:         repository.METHOD_GOOGLE.Method,
			ProviderUserID: &sub,
		}
		err = a.authStorage.CreateAuthProvider(ctx, googleAuthProvider)
		if err != nil {
			return UserInfo{}, err
		}
		return userInfo, nil
	}

	// check if sub exists
	googleAuthProvider, err := a.authStorage.GetAuthProviderByProviderId(ctx, sub, repository.METHOD_GOOGLE)

	// if email does not exist and sub exists: users google account changed emails => just change email and log in
	if err == nil {
		err := a.userService.UpdateUserEmail(ctx, googleAuthProvider.UserID, email)
		if err != nil {
			return UserInfo{}, err
		}

		user, err := a.userService.GetUserByEmail(ctx, email)
		if err != nil {
			return UserInfo{}, err
		}

		userInfo := UserInfo{
			Email:    *user.Email,
			UserId:   user.ID,
			Username: givenName,
		}

		return userInfo, nil
	}

	// if email and sub do not exist: create new account
	userInfo, err := a.userService.CreateUser(ctx, email, givenName)
	if err != nil {
		return UserInfo{}, err
	}
	googleAuthProvider = domain.AuthProvider{
		UserID:         userInfo.UserId,
		Method:         repository.METHOD_GOOGLE.Method,
		ProviderUserID: &sub,
	}
	err = a.userService.authStorage.CreateAuthProvider(ctx, googleAuthProvider)

	return userInfo, nil
}

func (a *AuthService) EmailSignUp(ctx context.Context, request SignUpWithEmailRequest) (UserInfo, error) {
	hashedPassword, err := passwords.HashPassword(request.Password)
	if err != nil {
		return UserInfo{}, fmt.Errorf("error hashing the user's password: %w", err)
	}

	emailExists := false

	user, err := a.userService.GetUserByEmail(ctx, request.Email)
	if err != nil && err != customErrors.ErrUserDoesNotExist {
		return UserInfo{}, err
	}
	if err != customErrors.ErrUserDoesNotExist {
		emailExists = true
	}

	if emailExists {
		userInfo := UserInfo{
			Email:    *user.Email,
			UserId:   user.ID,
			Username: user.Username,
		}

		_, err := a.authStorage.GetAuthProviderByUserId(ctx, user.ID, repository.METHOD_EMAIL)
		if err != nil && err != customErrors.ErrAuthProviderDoesNotExist {
			return UserInfo{}, err
		}

		// if email and email authProvider exist: do nothing -> could check password an silently login
		if err != customErrors.ErrAuthProviderDoesNotExist {
			return UserInfo{}, err
		}

		// if email exists and email authProvider does not: link email authProvider to that account
		emailAuthProvider := domain.AuthProvider{
			UserID:       user.ID,
			Method:       repository.METHOD_EMAIL.Method,
			PasswordHash: &hashedPassword,
		}
		err = a.authStorage.CreateAuthProvider(ctx, emailAuthProvider)
		if err != nil {
			return UserInfo{}, err
		}
		return userInfo, nil
	}

	// email does not exist: just create a ne user
	userInfo, err := a.userService.CreateUser(ctx, request.Email, "Alex")
	if err != nil {
		return UserInfo{}, err
	}

	emailAuthProvider := domain.AuthProvider{
		UserID:       userInfo.UserId,
		Method:       repository.METHOD_EMAIL.Method,
		PasswordHash: &hashedPassword,
	}
	err = a.userService.authStorage.CreateAuthProvider(ctx, emailAuthProvider)
	if err != nil {
		return UserInfo{}, err
	}

	return userInfo, err
}
