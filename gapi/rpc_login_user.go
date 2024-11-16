package gapi

import (
	"context"
	"errors"
	
	"github.com/google/uuid"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/pb"
	"github.com/katatrina/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}
	
	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}
	
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}
	
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}
	
	sessionID, err := uuid.Parse(refreshPayload.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse session ID")
	}
	
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           sessionID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}
	
	resp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt.Time),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt.Time),
	}
	return resp, nil
}
