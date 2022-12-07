package service

import (
	"context"
	"time"

	pb "github.com/SaidovZohid/medium_user_service/genproto/user_service"
	"github.com/SaidovZohid/medium_user_service/storage"
	"github.com/SaidovZohid/medium_user_service/storage/repo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
}

func NewUserService(strg storage.StorageI, inMemory storage.InMemoryStorageI) *UserService {
	return &UserService{
		storage:  strg,
		inMemory: inMemory,
	}
}

func (s *UserService) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().Create(&repo.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     &req.PhoneNumber,
		Email:           req.Email,
		Gender:          &req.Gender,
		Password:        req.Password,
		UserName:        &req.Username,
		ProfileImageUrl: &req.ProfileImageUrl,
		Type:            req.Type,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     *user.PhoneNumber,
		Email:           user.Email,
		Gender:          *user.Gender,
		Password:        user.Password,
		Username:        *user.UserName,
		ProfileImageUrl: *user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}, nil
}
