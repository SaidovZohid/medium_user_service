package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SaidovZohid/medium_user_service/config"
	"github.com/SaidovZohid/medium_user_service/genproto/notification_service"
	pb "github.com/SaidovZohid/medium_user_service/genproto/user_service"
	grpcPkg "github.com/SaidovZohid/medium_user_service/pkg/grpc_client"
	"github.com/SaidovZohid/medium_user_service/pkg/utils"
	"github.com/SaidovZohid/medium_user_service/storage"
	"github.com/SaidovZohid/medium_user_service/storage/repo"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	storage    storage.StorageI
	inMemory   storage.InMemoryStorageI
	grpcClient grpcPkg.GrpcClientI
	cfg        *config.Config
	logger     *logrus.Logger
}

func NewAuthService(strg storage.StorageI, inMemory storage.InMemoryStorageI, grpc grpcPkg.GrpcClientI, cfg *config.Config, log *logrus.Logger) *AuthService {
	return &AuthService{
		storage:    strg,
		inMemory:   inMemory,
		grpcClient: grpc,
		cfg:        cfg,
		logger:     log,
	}
}

const (
	RegisterCodeKey   = "register_code_"
	ForgotPasswordKey = "forgot_password_code_"
)

const (
	VerificationEmail   = "verification_email"
	ForgotPasswordEmail = "forgot_password_email"
)

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.WithError(err).Error("failed to hash password")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	user := repo.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Type:      repo.UserTypeUser,
		Password:  hashedPassword,
	}

	userData, err := json.Marshal(user)
	if err != nil {
		s.logger.WithError(err).Error("failed to marshal user in register func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	err = s.inMemory.Set("user_"+req.Email, string(userData), 10*time.Minute)
	if err != nil {
		s.logger.WithError(err).Error("failed to set user to redis in register func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	go func() {
		err := s.sendVereficationCode(RegisterCodeKey, req.Email, VerificationEmail)
		if err != nil {
			s.logger.WithError(err).Error("failed to send verification code in register func")
		}
	}()

	return &emptypb.Empty{}, nil
}

func (s *AuthService) sendVereficationCode(key, email, email_type string) error {
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		s.logger.WithError(err).Error("failed to generate random code in sendVerification func")
		return err
	}

	err = s.inMemory.Set(key+email, code, time.Minute)
	if err != nil {
		s.logger.WithError(err).Error("failed to send generated code in sendVerificationCode func")
		return err
	}
	_, err = s.grpcClient.NotificationService().SendEmail(context.Background(), &notification_service.SendEmailRequest{
		To:      email,
		Subject: "Verification Email",
		Body: map[string]string{
			"code": code,
		},
		Type: email_type,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.AuthResponse, error) {
	userData, err := s.inMemory.Get("user_" + req.Email)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user data from redis in verify func")
		return nil, status.Errorf(codes.Internal, "internale server error: %v", err)
	}
	var user repo.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		s.logger.WithError(err).Error("failed to unmarshal user data in verify func")
		return nil, status.Errorf(codes.Internal, "internal server errror: %v", err)
	}

	code, err := s.inMemory.Get(RegisterCodeKey + user.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "verification code is expired: %v", err)
	}
	if code != req.Code {
		return nil, status.Errorf(codes.Unknown, "verification code is incorrect: %v", err)
	}

	result, err := s.storage.User().Create(&user)
	if err != nil {
		s.logger.WithError(err).Error("failed to create user in verify func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}
	token, _, err := utils.CreateToken(s.cfg, &utils.TokenParams{
		UserID:   user.ID,
		Email:    user.Email,
		UserType: user.Type,
		Duration: time.Hour * 24 * 360,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create token in verify func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.AuthResponse{
		Id:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		Type:        result.Type,
		CreatedAt:   result.CreatedAt.Format(time.RFC3339),
		AccessToken: token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, err := s.storage.User().GetByEmail(req.Email)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user by email in login func")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not fount please register: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		s.logger.WithError(err).Error("failed to check password in login func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	token, _, err := utils.CreateToken(s.cfg, &utils.TokenParams{
		UserID:   user.ID,
		Email:    user.Email,
		UserType: user.Type,
		Duration: time.Hour * 24 * 360,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create token in login func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.AuthResponse{
		Id:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		AccessToken: token,
	}, nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*emptypb.Empty, error) {
	_, err := s.storage.User().GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "email does not exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	go func() {
		err := s.sendVereficationCode(ForgotPasswordKey, req.Email, ForgotPasswordEmail)
		if err != nil {
			fmt.Printf("failed to send verification code: %v", err)
		}
	}()

	return &emptypb.Empty{}, nil
}

func (s *AuthService) VerifyForgotPassword(ctx context.Context, req *pb.VerifyRequest) (*pb.AuthResponse, error) {
	code, err := s.inMemory.Get(ForgotPasswordKey + req.Email)
	if err != nil {
		s.logger.WithError(err).Error("failed to get code from redis in VerifyForgoPasword func")
		return nil, status.Errorf(codes.Internal, "verification code has been expired: %v", err)
	}

	if req.Code != code {
		return nil, status.Errorf(codes.Internal, "verification code is not true: %v", err)

	}

	result, err := s.storage.User().GetByEmail(req.Email)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user info by email in verifyforgotpassword func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	token, _, err := utils.CreateToken(s.cfg, &utils.TokenParams{
		UserID:   result.ID,
		Email:    result.Email,
		UserType: result.Type,
		Duration: time.Minute * 30,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to create token in VerifyForgotPassword func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.AuthResponse{
		Id:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		Type:        result.Type,
		CreatedAt:   result.CreatedAt.Format(time.RFC3339),
		AccessToken: token,
	}, nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.WithError(err).Error("failed to hashpassword in UpdatePAssword func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	err = s.storage.User().UpdatePassword(&repo.UpdatePassword{
		UserID:   req.UserId,
		Password: hashedPassword,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to update password in UpdatePassword func")
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.AuthPayload, error) {
	accessToken := req.AccessToken

	payload, err := utils.VerifyToken(s.cfg, accessToken)
	if err != nil {
		s.logger.WithError(err).Error("failed to verify token in VerifyToken func")
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return &pb.AuthPayload{
		Id:        payload.Id.String(),
		UserId:    payload.UserID,
		Email:     payload.Email,
		UserType:  payload.UserType,
		IssuedAt:  payload.IssuedAt.Format(time.RFC3339),
		ExpiredAt: payload.ExpiredAt.Format(time.RFC3339),
	}, nil
}
