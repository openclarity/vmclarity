package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/apiserver/common"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/types"
)

// TODO: done.
func (s *ServerImpl) GetUsersUserID(ctx echo.Context, userID models.UserID) error {
	result, err := s.authStore.GetUser(userID)
	if err != nil {
		return handleStoreErr(ctx, err, "get user")
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done.
func (s *ServerImpl) GetCurrentUser(ctx echo.Context) error {
	user := iam.GetUserFromContext(ctx)
	if user == nil {
		return sendError(ctx, http.StatusUnauthorized, "not authenticated")
	}
	return sendResponse(ctx, http.StatusOK, user)
}

// TODO: done.
func (s *ServerImpl) DeleteUsersUserID(ctx echo.Context, userID models.UserID) error {
	err := s.authStore.DeleteUser(userID)
	if err != nil {
		return handleStoreErr(ctx, err, "delete user")
	}
	return sendResponse(ctx, http.StatusOK, "deleted successfully")
}

// TODO: done.
func (s *ServerImpl) PatchUsersUserID(ctx echo.Context, userID models.UserID) error {
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	result, err := s.authStore.UpdateUser(user)
	if err != nil {
		return handleStoreErr(ctx, err, "update user")
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done.
func (s *ServerImpl) PutUsersUserID(ctx echo.Context, userID models.UserID) error {
	updateErr := s.PatchUsersUserID(ctx, userID)
	if updateErr == nil {
		// Successfully updated
		return nil
	}
	if errors.Is(updateErr, types.ErrNotFound) {
		// Resource is not found, create instead
		return s.PostUser(ctx)
	}
	return updateErr
}

// TODO: done.
func (s *ServerImpl) DeleteUserAuthUserID(ctx echo.Context, userID models.UserID) error {
	var userAuth models.UserAuth
	err := ctx.Bind(&userAuth)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	err = s.authStore.RevokeUserAuth(userID, userAuth)
	if err != nil {
		return handleStoreErr(ctx, err, "revoke user auth")
	}
	return sendResponse(ctx, http.StatusOK, "successfully revoked")
}

// TODO: done.
func (s *ServerImpl) GetUserAuthUserID(ctx echo.Context, userID models.UserID) error {
	result, err := s.authStore.GetUserAuth(userID)
	if err != nil {
		return handleStoreErr(ctx, err, "get user auth")
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done.
func (s *ServerImpl) PostUserAuthUserID(ctx echo.Context, userID models.UserID, params models.PostUserAuthUserIDParams) error {
	result, err := s.authStore.CreateUserAuth(userID, models.CredentialType(params.CredentialType), params.CredentialExpiry)
	if err != nil {
		return handleStoreErr(ctx, err, "create user auth")
	}
	return sendResponse(ctx, http.StatusCreated, result)
}

// TODO: done.
func (s *ServerImpl) GetUsers(ctx echo.Context, params models.GetUsersParams) error {
	result, err := s.authStore.GetUsers(params)
	if err != nil {
		return handleStoreErr(ctx, err, "get users")
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done.
func (s *ServerImpl) PostUser(ctx echo.Context) error {
	// Get data from request
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// Create user
	result, err := s.authStore.CreateUser(user)
	if err != nil {
		return handleStoreErr(ctx, err, "create user")
	}
	return sendResponse(ctx, http.StatusCreated, result)
}

func handleStoreErr(ctx echo.Context, err error, opMsg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, types.ErrNotFound) {
		return sendResponse(ctx, http.StatusNotFound, err.Error())
	}
	if errors.Is(err, types.ErrAlreadyExists) {
		return sendResponse(ctx, http.StatusConflict, err.Error())
	}
	if badRequestErr := (*common.BadRequestError)(nil); errors.As(err, &badRequestErr) {
		return sendResponse(ctx, http.StatusBadRequest, badRequestErr.Reason)
	}
	if conflictErr := (*common.ConflictError)(nil); errors.As(err, &conflictErr) {
		return sendResponse(ctx, http.StatusConflict, conflictErr.Reason)
	}
	return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to %s: %v", opMsg, err))
}
