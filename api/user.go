package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/token"
	"github.com/katatrina/simplebank/util"
	"github.com/katatrina/simplebank/validator"
	"github.com/katatrina/simplebank/worker"
	"github.com/rs/zerolog/log"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type createUserResponse struct {
	User db.User `json:"user"`
}

func validateCreateUserRequest(req *createUserRequest) (violations []*validator.FieldViolation) {
	if err := validator.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	
	if err := validator.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	
	if err := validator.ValidateFullName(req.FullName); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	
	if err := validator.ValidateEmail(req.Email); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	
	return violations
}

func (server *Server) createUser(ctx *gin.Context) {
	req := new(createUserRequest)
	
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	violations := validateCreateUserRequest(req)
	if violations != nil {
		ctx.JSON(http.StatusBadRequest, invalidArgumentError(violations))
		return
	}
	
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.Username,
			HashedPassword: hashedPassword,
			FullName:       req.FullName,
			Email:          req.Email,
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			
			err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
			if err != nil {
				log.Error().Err(err).Str("type", worker.TaskSendVerifyEmail).Str("email", user.Email).Msg("failed to distribute task")
				return fmt.Errorf("failed to send verify email to username: %s", user.Username)
			}
			
			return nil
		},
	}
	
	txResult, err := server.store.CreateUserTx(context.Background(), arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolationCode {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Header("Location", fmt.Sprintf("/v1/users/%s", txResult.User.Username))
	
	ctx.JSON(http.StatusOK, createUserResponse{User: txResult.User})
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  db.User   `json:"user"`
}

func validateLoginUserRequest(req *loginUserRequest) (violations []*validator.FieldViolation) {
	if err := validator.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	
	if err := validator.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	
	return violations
}

func (server *Server) loginUser(ctx *gin.Context) {
	req := new(loginUserRequest)
	
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	violations := validateLoginUserRequest(req)
	if violations != nil {
		ctx.JSON(http.StatusBadRequest, invalidArgumentError(violations))
		return
	}
	
	user, err := server.store.GetUser(context.Background(), req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		err = errors.New("incorrect password")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	sessionID, err := uuid.Parse(refreshPayload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	session, err := server.store.CreateSession(context.Background(), db.CreateSessionParams{
		ID:           sessionID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.GetHeader("User-Agent"),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	resp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  user,
	}
	ctx.JSON(http.StatusOK, resp)
}

type updateUserRequest struct {
	Username string  `json:"username"`
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (req *updateUserRequest) getUserName() string {
	if req != nil {
		return req.Username
	}
	
	return ""
}

func (req *updateUserRequest) getFullName() string {
	if req != nil && req.FullName != nil {
		return *req.FullName
	}
	
	return ""
}

func (req *updateUserRequest) getEmail() string {
	if req != nil && req.Email != nil {
		return *req.Email
	}
	
	return ""
}

func (req *updateUserRequest) getPassword() string {
	if req != nil && req.Password != nil {
		return *req.Password
	}
	
	return ""
}

func validateUpdateUserRequest(req *updateUserRequest) (violations []*validator.FieldViolation) {
	if err := validator.ValidateUsername(req.getUserName()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	
	// If the password field is not nil, validate the password.
	// If the password field is nil, it means the user does not want to update the password.
	if req.Password != nil {
		if err := validator.ValidatePassword(req.getPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	
	if req.FullName != nil {
		if err := validator.ValidateFullName(req.getFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	
	if req.Email != nil {
		if err := validator.ValidateEmail(req.getEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	
	return violations
}

func (server *Server) updateUser(ctx *gin.Context) {
	req := new(updateUserRequest)
	
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		ctx.JSON(http.StatusBadRequest, invalidArgumentError(violations))
		return
	}
	
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	
	if authPayload.Role != util.BankerRole && authPayload.Subject != req.getUserName() {
		err := errors.New("cannot update other user's information")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	arg := db.UpdateUserParams{
		Username: req.getUserName(),
		FullName: pgtype.Text{
			String: req.getFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.getEmail(),
			Valid:  req.Email != nil,
		},
	}
	
	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.getPassword())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}
		
		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}
	
	user, err := server.store.UpdateUser(context.Background(), arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = errors.New("user not found")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		err = fmt.Errorf("failed to update user: %w", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, user)
}

type verifyUserEmailRequest struct {
	EmailID    int64  `form:"email_id"`
	SecretCode string `form:"secret_code"`
}

type verifyUserEmailResponse struct {
	IsVerified bool `json:"is_verified"`
}

func validateVerifyUserEmailRequest(req *verifyUserEmailRequest) (violations []*validator.FieldViolation) {
	if err := validator.ValidateEmailID(req.EmailID); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	
	if err := validator.ValidateSecretCode(req.SecretCode); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}
	
	return violations
}

func (server *Server) verifyUserEmail(ctx *gin.Context) {
	req := new(verifyUserEmailRequest)
	
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	violations := validateVerifyUserEmailRequest(req)
	if violations != nil {
		ctx.JSON(http.StatusBadRequest, invalidArgumentError(violations))
		return
	}
	
	txResult, err := server.store.VerifyUserEmailTx(context.Background(), db.VerifyUserEmailTxParams{
		EmailID:    req.EmailID,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		err = fmt.Errorf("failed to verify email: %w", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	resp := verifyUserEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}
	ctx.JSON(http.StatusOK, resp)
}
