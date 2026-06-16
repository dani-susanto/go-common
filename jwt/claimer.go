package jwt

import (
	"context"
	"strconv"
)

func GetClaims(ctx context.Context) (*ClaimsInput, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*ClaimsInput)
	return claims, ok
}

func GetUserID(ctx context.Context) (int, bool) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return 0, false
	}
	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		return 0, false
	}
	return userID, true
}

func GetName(ctx context.Context) (string, bool) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return "", false
	}
	return claims.Name, true
}

func GetEmail(ctx context.Context) (string, bool) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return "", false
	}
	return claims.Email, true
}

func GetRole(ctx context.Context) (string, bool) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return "", false
	}
	return claims.Role, true
}

func GetActorID(ctx context.Context) (int, bool) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return 0, false
	}
	if claims.UserID == "" {
		return 0, false
	}
	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		return 0, false
	}
	return userID, true
}
