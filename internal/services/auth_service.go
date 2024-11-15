package services

import (
    "context"
    "errors"
    "log"

    "github.com/dgrijalva/jwt-go"
    "github.com/yemyoaung/managing-vehicle-tracking-auth-svc/internal/config"
    "github.com/yemyoaung/managing-vehicle-tracking-auth-svc/internal/repositories"
    "github.com/yemyoaung/managing-vehicle-tracking-common"
    "github.com/yemyoaung/managing-vehicle-tracking-models"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
)

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthService interface {
    Login(ctx context.Context, req *LoginRequest) (*models.User, string, error)
    ValidateToken(ctx context.Context, token string) (*models.User, error)
}

type MongoAuthService struct {
    repo       repositories.AuthRepository
    tokenMaker common.TokenMaker
    cfg        *config.EnvConfig
}

func NewMongoAuthService(
    repo repositories.AuthRepository,
    tokenMaker common.TokenMaker,
    cfg *config.EnvConfig,
) AuthService {
    return &MongoAuthService{
        repo:       repo,
        tokenMaker: tokenMaker,
        cfg:        cfg,
    }
}

func (s *MongoAuthService) Login(ctx context.Context, req *LoginRequest) (
    *models.User,
    string,
    error,
) {
    var user models.User
    err := s.repo.FindAdminByEmail(ctx, req.Email, &user)
    if err != nil {
        return nil, "", ErrInvalidCredentials
    }

    if !common.CheckPasswordHash(req.Password, user.Password) {
        return nil, "", ErrInvalidCredentials
    }

    token, err := s.tokenMaker.CreateToken(
        user.Claim(), s.cfg.JwtSecret,
    )

    if err != nil {
        return nil, "", err
    }

    return &user, token, nil
}

func (s *MongoAuthService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
    verifyToken, err := s.tokenMaker.VerifyToken(token, s.cfg.JwtSecret, &jwt.StandardClaims{})
    if err != nil {
        log.Println("Error verifying token", err)
        return nil, err
    }

    switch verifyToken.(type) {
    case *jwt.StandardClaims:
        claims := verifyToken.(*jwt.StandardClaims)
        var user models.User
        err := s.repo.FindAdminByID(ctx, claims.Id, &user)
        if err != nil {
            return nil, err
        }
        return &user, nil
    default:
    }
    return nil, err
}
