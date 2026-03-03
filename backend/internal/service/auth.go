package service

import (
	"backend/internal/jwt"
	"backend/internal/storage/cache"
	"backend/internal/storage/postgres"
	"backend/pkg/generateNumber"
	"backend/pkg/gmail"
	authv1 "backend/proto/gen/auth"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authv1.UnimplementedAuthServer
	db    *postgres.DB
	redis *cache.DB
}

func NewAuthService(db *postgres.DB, redis *cache.DB) *AuthService {
	return &AuthService{
		db:    db,
		redis: redis,
	}
}

func (s *AuthService) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	isExists, err := s.db.IsUserExists(ctx, req.Email)
	if err != nil {
		return &authv1.RegisterResponse{Success: false}, err
	}

	if isExists {
		return &authv1.RegisterResponse{Success: false}, fmt.Errorf("user already exists")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	code, err := generatenumber.GenerateFiveDigitNumber()
	if err != nil {
		return nil, err
	}

	pending := &cache.PendingUser{
		Email:    req.Email,
		Password: string(hash),
		Username: req.Username,
		Code:     strconv.Itoa(code),
		Expires:  time.Now().Add(5 * time.Minute),
	}
	if err := s.redis.SavePendingUser(ctx, pending, 5*time.Minute); err != nil {
		return nil, err
	}

	go func() {
		if err := gmail.SendVerificationCode(req.Email, strconv.Itoa(code)); err != nil {
			log.Printf("failed to send email to %s: %v", req.Email, err)
		}
	}()

	return &authv1.RegisterResponse{Success: true}, nil
}

func (s *AuthService) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &authv1.LoginResponse{Token: ""}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &authv1.LoginResponse{Token: ""}, err
	}

	token, err := jwt.NewToken(user.ID, req.Email, 8760*time.Hour)
	if err != nil {
		return &authv1.LoginResponse{Token: ""}, err
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *AuthService) Verify(ctx context.Context, req *authv1.VerifyRequest) (*authv1.VerifyResponse, error) {
	pending, err := s.redis.GetPendingUser(ctx, req.Email)
	if err != nil {
		return &authv1.VerifyResponse{Success: false, Token: ""}, nil
	}

	if pending.Code != req.Code || pending.Expires.Before(time.Now()) {
		return &authv1.VerifyResponse{Success: false, Token: ""}, nil
	}

	id, err := s.db.CreateUser(ctx, pending.Email, pending.Username, pending.Password)
	if err != nil {
		return &authv1.VerifyResponse{Success: false, Token: ""}, err
	}

	err = s.redis.DeletePendingUser(ctx, req.Email)
	if err != nil {
		return &authv1.VerifyResponse{Success: false, Token: ""}, err
	}

	token, err := jwt.NewToken(id, req.Email, 8760*time.Hour)
	if err != nil {
		return &authv1.VerifyResponse{Success: false, Token: ""}, err
	}

	return &authv1.VerifyResponse{Success: true, Token: token}, nil
}

func (s *AuthService) NewPass(ctx context.Context, req *authv1.NewPassRequest) (*authv1.NewPassResponse, error) {
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &authv1.NewPassResponse{Success: false}, err
	}
	token, err := jwt.NewToken(user.ID, user.Email, 10*time.Minute)
	if err != nil {
		return &authv1.NewPassResponse{Success: false}, err
	}
	link := "http://localhost:3000/restore-pass/" + token

	go func() {
		if err := gmail.SendPasswordResetEmail(req.Email, link); err != nil {
			log.Printf("failed to send email to %s: %v", req.Email, err)
		}
	}()
	pass := &cache.NewPass{
		Email: req.Email,
		Token: token,
	}
	err = s.redis.SavePendingPass(ctx, pass)
	if err != nil {
		return &authv1.NewPassResponse{Success: false}, err
	}
	return &authv1.NewPassResponse{
		Success: true,
		Message: "Письмо отправлено",
	}, nil
}

func (s *AuthService) ValidateRestoreToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	claims, err := jwt.ValidateToken(req.Token)

	if err != nil {
		return &authv1.ValidateTokenResponse{
			Success: false,
			Email:   "",
		}, nil
	}

	userID, ok := claims["uid"].(float64)
	if !ok {
		return &authv1.ValidateTokenResponse{
			Success: false,
			Email:   "",
		}, nil
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return &authv1.ValidateTokenResponse{
			Success: false,
			Email:   "",
		}, nil
	}

	if exp, ok := claims["exp"].(int64); ok {
		if time.Now().After(time.Unix(exp, 0)) {
			return &authv1.ValidateTokenResponse{
				Success: false,
				Email:   "",
			}, nil
		}
	}

	user, err := s.db.GetUserByID(ctx, int64(userID))
	if err != nil {
		return &authv1.ValidateTokenResponse{
			Success: false,
			Email:   "",
		}, nil
	}

	if user.Email != email {
		return &authv1.ValidateTokenResponse{
			Success: false,
			Email:   "",
		}, nil
	}

	return &authv1.ValidateTokenResponse{Success: true, Email: email}, nil
}

func (s *AuthService) SetNewPassword(ctx context.Context, req *authv1.SetNewPassRequest) (*authv1.SetNewPassResponse, error) {
	newPass := req.NewPass
	claims, err := jwt.ValidateToken(req.Token)
	if err != nil {
		return &authv1.SetNewPassResponse{Success: false}, err
	}
	userID, ok := claims["uid"].(float64)
	if !ok {
		return &authv1.SetNewPassResponse{
			Success: false,
		}, nil
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return &authv1.SetNewPassResponse{
			Success: false,
		}, nil
	}

	err = s.redis.DeletePendingPass(ctx, email)
	if err != nil {
		return &authv1.SetNewPassResponse{Success: false}, err
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)

	err = s.db.UpdatePassword(ctx, int64(userID), string(hash))
	if err != nil {
		return &authv1.SetNewPassResponse{
			Success: false,
		}, err
	}

	return &authv1.SetNewPassResponse{Success: true}, nil
}
