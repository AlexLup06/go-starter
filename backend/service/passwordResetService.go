package service

import (
	"context"
	"fmt"

	"alexlupatsiy.com/personal-website/backend/domain"
	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/helpers/passwords"
	"alexlupatsiy.com/personal-website/backend/helpers/token"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type PasswordResetService struct {
	passwordResetStorage repository.PasswordResetStorage
	userService          *UserService
	tokenService         *TokenService
	mailService          *MailService
}

func NewPasswordResetService(
	passwordResetStorage repository.PasswordResetStorage,
	userService *UserService,
	tokenService *TokenService,
	mailService *MailService,
) *PasswordResetService {
	return &PasswordResetService{
		passwordResetStorage: passwordResetStorage,
		userService:          userService,
		tokenService:         tokenService,
		mailService:          mailService,
	}
}

func (p *PasswordResetService) GenerateResetToken(ctx context.Context, email string) (string, error) {
	user, err := p.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// check if user has generated reset token recently
	// => if has geneated two in the last 15 minutes then custom error with GeneratedTooMany
	amountTokens, err := p.passwordResetStorage.GetAmountTokensYoungerThan15Min(ctx, user.ID)
	if err != nil {
		return "", err
	}

	if amountTokens > 2 {
		return "", customErrors.ErrGeneratedTooManyResetTokens
	}

	// delete all expired reset tokens
	err = p.passwordResetStorage.DeleteAllTokensOlderThan15min(ctx, user.ID)

	// revoke the other tokens
	err = p.passwordResetStorage.RevokeAllTokens(ctx, user.ID)

	// generate reset Token
	userInfo := UserInfo{
		UserId:   user.ID,
		Username: user.Username,
		Email:    *user.Email,
	}
	tokenString, token, _, err := p.tokenService.GenerateUserInfoJWT(userInfo, 15)
	if err != nil {
		return "", err
	}
	hashedToken := passwords.HashToken(tokenString)

	// save in db
	passwordReset := domain.PasswordReset{
		UserID:     user.ID,
		ExpiresAt:  token.ExpiresAt.Time,
		ResetToken: hashedToken,
		Used:       false,
	}
	err = p.passwordResetStorage.CreatePasswordResetToken(ctx, user.ID, passwordReset)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (p *PasswordResetService) SendResetPasswordEmailAsync(ctx context.Context, resetToken string) {
	go func() {
		body := "http://localhost:3000/auth/reset-password?token=" + resetToken
		// err := p.mailService.SendMail("noreply@alexlupatsiy.com", "alexander@lupatsiy.de", "Password Reset", body)
		// if err != nil {
		// 	// TODO: Handle logging
		// 	fmt.Println(err)
		// }
		fmt.Println(body)
	}()
}

func (p *PasswordResetService) CheckResetPasswordToken(ctx context.Context, resetPasswordRequest ResetPasswordRequest) (*token.CustomClaims[UserInfo], error) {
	hashedToken := passwords.HashToken(resetPasswordRequest.Token)
	err := p.passwordResetStorage.CheckToken(ctx, hashedToken)
	if err != nil {
		return nil, err
	}

	token, err := p.tokenService.ParseUserInfoJWT(resetPasswordRequest.Token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (p *PasswordResetService) RevokeAllTokens(ctx context.Context, userId string) error {
	err := p.passwordResetStorage.RevokeAllTokens(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}
