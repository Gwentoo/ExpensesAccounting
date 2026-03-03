package jwt

import (
	"backend/internal/config"
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ContextKey string

const ClaimsKey ContextKey = "claims"

func CookieAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if info.FullMethod == "/auth.Auth/Login" ||
		info.FullMethod == "/auth.Auth/Register" ||
		info.FullMethod == "/auth.Auth/Verify" ||
		info.FullMethod == "/auth.Auth/NewPass" ||
		info.FullMethod == "/auth.Auth/ValidateRestoreToken" ||
		info.FullMethod == "/auth.Auth/SetNewPassword" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is missing")
	}
	var token string

	if cookies := md.Get("cookie"); len(cookies) > 0 {
		for _, c := range cookies {
			if strings.HasPrefix(c, "jwt=") {
				token = strings.TrimPrefix(c, "jwt=")
				break
			}
		}
	}

	if token == "" {
		if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
			for _, authHeader := range authHeaders {
				if strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
					break
				}
			}
		}
	}

	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "authorization cookie missing")
	}

	claims, err := ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	newCtx := context.WithValue(ctx, ClaimsKey, claims)

	return handler(newCtx, req)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {

	if tokenString == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}

	cfg := config.LoadConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JwtSecretKey), nil
	})

	if err != nil || token == nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, status.Error(codes.Unauthenticated, "invalid claims")
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {

	claims, ok := ctx.Value(ClaimsKey).(jwt.MapClaims)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "claims not found in context")
	}

	userIDFloat, ok := claims["uid"].(float64)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "user id not found in token")
	}

	return int64(userIDFloat), nil
}
