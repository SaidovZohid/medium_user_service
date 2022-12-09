package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SaidovZohid/medium_user_service/genproto/notification_service"
	pb "github.com/SaidovZohid/medium_user_service/genproto/user_service"
	grpcPkg "github.com/SaidovZohid/medium_user_service/pkg/grpc_client"
	"github.com/SaidovZohid/medium_user_service/pkg/utils"
	"github.com/SaidovZohid/medium_user_service/storage"
	"github.com/SaidovZohid/medium_user_service/storage/repo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	storage    storage.StorageI
	inMemory   storage.InMemoryStorageI
	grpcClient grpcPkg.GrpcClientI
}

func NewAuthService(strg storage.StorageI, inMemory storage.InMemoryStorageI, grpc grpcPkg.GrpcClientI) *AuthService {
	return &AuthService{
		storage:    strg,
		inMemory:   inMemory,
		grpcClient: grpc,
	}
}

const (
	RegisterCodeKey   = "register_code_"
	ForgotPasswordKey = "forgot_password_code_"
)

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
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
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	err = s.inMemory.Set("user_"+req.Email, string(userData), 10*time.Minute)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	go func() {
		err := s.sendVereficationCode(RegisterCodeKey, req.Email)
		if err != nil {
			fmt.Printf("failed to send verification code: %v", err)
		}
	}()

	return &emptypb.Empty{}, nil
}

func (s *AuthService) sendVereficationCode(key, email string) error {
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	err = s.inMemory.Set(key+email, code, time.Minute)
	if err != nil {
		return err
	}
	_, err = s.grpcClient.NotificationService().SendEmail(context.Background(), &notification_service.SendEmailRequest{
		To:      email,
		Subject: "Verification Email",
		Body: map[string]string{
			"code": code,
		},
		Type: "verification_email",
	})
	if err != nil {
		return err
	}

	return nil
}
