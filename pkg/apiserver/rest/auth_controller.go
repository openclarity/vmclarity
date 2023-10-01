package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"net/http"
)

func (s *ServerImpl) AuthRedirect(ctx echo.Context) error {
	return sendResponse(ctx, http.StatusOK, ctx.Request().Header)
}

// TODO: done
func (s *ServerImpl) GetUsersUserID(ctx echo.Context, userID models.UserID) error {
	result, err := s.authStore.GetUser(userID)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done
func (s *ServerImpl) DeleteUserAuthUserID(ctx echo.Context, userID models.UserID) error {
	var userAuth models.UserAuth
	err := ctx.Bind(&userAuth)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	err = s.authStore.RevokeUserAuth(userID, userAuth)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, "successfully revoked")
}

// TODO: done
func (s *ServerImpl) GetUserAuthUserID(ctx echo.Context, userID models.UserID) error {
	result, err := s.authStore.GetUserAuth(userID)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done
func (s *ServerImpl) PostUserAuthUserID(ctx echo.Context, userID models.UserID, params models.PostUserAuthUserIDParams) error {
	result, err := s.authStore.CreateUserAuth(userID, models.AuthType(params.AuthType), params.ExpiryTime)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done
func (s *ServerImpl) GetCurrentUser(ctx echo.Context) error {
	currUser := iam.GetUserFromContext(ctx)
	if currUser == nil {
		return sendError(ctx, http.StatusBadRequest, "not authenticated")
	}
	return sendResponse(ctx, http.StatusOK, currUser)
}

// TODO: done
func (s *ServerImpl) DeleteUsersUserID(ctx echo.Context, userID models.UserID) error {
	err := s.authStore.DeleteUser(userID)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, "deleted successfully")
}

// TODO: done
func (s *ServerImpl) PatchUsersUserID(ctx echo.Context, userID models.UserID) error {
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	result, err := s.authStore.UpdateUser(user)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done
func (s *ServerImpl) PutUsersUserID(ctx echo.Context, userID models.UserID) error {
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	var fetchedUser models.User
	_, fetchErr := s.authStore.GetUser(userID)
	if fetchErr != nil {
		fetchedUser, err = s.authStore.CreateUser(user)
	} else {
		fetchedUser, err = s.authStore.UpdateUser(user)
	}
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, fetchedUser)
}

// TODO: done
func (s *ServerImpl) GetUsers(ctx echo.Context, params models.GetUsersParams) error {
	result, err := s.authStore.GetUsers(params)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, result)
}

// TODO: done
func (s *ServerImpl) PostUser(ctx echo.Context) error {
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	result, err := s.authStore.CreateUser(user)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, result)
}
