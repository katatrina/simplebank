package gapi

import (
	"context"
	"errors"
	
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/pb"
	"github.com/katatrina/simplebank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username or email already exists: %s", err)
			}
		}
		
		return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
	}
	
	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return resp, nil
}
