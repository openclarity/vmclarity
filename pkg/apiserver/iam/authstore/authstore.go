// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authstore

import (
	"context"
	"fmt"
	"time"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/apiserver/common"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/types"

	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	auth_pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/authn"
	management_pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	object_pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	user_pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ types.AuthStore = &zitadelStore{}

type zitadelStore struct {
	projectID string
	mgmt      *management.Client
}

func New() (types.AuthStore, error) {
	// Load config
	config := LoadConfig()

	// Create management client
	mgmtClient, err := management.NewClient(
		config.Issuer,
		config.API,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI(), zitadel.ScopeProjectID(config.ProjectID)},
		zitadel.WithInsecure(),
		zitadel.WithOrgID(config.OrgID),
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(config.AuthKeyPath)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zitadel management client: %w", err)
	}

	return &zitadelStore{
		projectID: config.ProjectID,
		mgmt:      mgmtClient,
	}, nil
}

func (z *zitadelStore) GetUserFromInfo(info *types.UserInfo) (models.User, error) {
	// For Zitadel OIDC, token subject is the user ID
	if info.FromZitadelOIDC {
		return z.GetUser(info.Subject)
	}

	// For Generic OIDC, user must be synced to DB first
	// TODO: Implement fetch/create for a generic OIDC
	return models.User{}, fmt.Errorf("not implemented current user fetcher for generic OIDC")
}

func (z *zitadelStore) GetUser(userID models.UserID) (models.User, error) {
	// Get user from API
	userResp, err := z.mgmt.GetUserByID(context.Background(), &management_pb.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		if isRPCCode(err, codes.NotFound) { // NotFound error
			return models.User{}, types.ErrNotFound
		}
		return models.User{}, fmt.Errorf("failed to get user from API: %w", err) // Generic error
	}

	// Get user roles
	return z.getUserData(userResp.User)
}

func (z *zitadelStore) GetUsers(params models.GetUsersParams) (models.Users, error) {
	// Get grants from API
	grantResp, err := z.mgmt.ListUserGrants(context.Background(), &management_pb.ListUserGrantRequest{
		Query: &object_pb.ListQuery{
			Offset: 0,
			Limit:  0,
			Asc:    false,
		},
		Queries: []*user_pb.UserGrantQuery{
			{
				Query: &user_pb.UserGrantQuery_ProjectIdQuery{
					ProjectIdQuery: &user_pb.UserGrantProjectIDQuery{
						ProjectId: z.projectID,
					},
				},
			},
		},
	})
	if err != nil {
		return models.Users{}, fmt.Errorf("failed to list grants from API: %w", err) // Generic error
	}

	// Extract user roles from grants
	userRolesStore := make(map[string][]string)
	for _, grant := range grantResp.Result {
		currRoles := userRolesStore[grant.UserId]
		userRolesStore[grant.UserId] = append(currRoles, grant.RoleKeys...)
	}

	// Get users from API
	paramSkip := 0
	if params.Skip != nil {
		paramSkip = *params.Skip
	}
	paramTop := 100
	if params.Top != nil {
		paramTop = *params.Top
	}
	userResp, err := z.mgmt.ListUsers(context.Background(), &management_pb.ListUsersRequest{
		Query: &object_pb.ListQuery{
			Offset: uint64(paramSkip),
			Limit:  uint32(paramTop),
		},
	})
	if err != nil {
		return models.Users{}, fmt.Errorf("failed to list users from API: %w", err) // Generic error
	}

	// Transform user response into DB objects
	users := make([]models.User, len(userResp.Result))
	for i, user := range userResp.Result {
		// Get user data
		userRoles := userRolesStore[user.Id]
		userType := models.REGULARUSER
		if user.GetMachine() != nil {
			userType = models.MACHINEUSER
		}

		users[i] = models.User{
			Banned:    toPointer(user.State == user_pb.UserState_USER_STATE_LOCKED),
			CreatedAt: toPointer(user.Details.CreationDate.AsTime()),
			Email:     toPointer(user.UserName),
			Id:        toPointer(user.Id),
			Name:      toPointer(user.UserName),
			Roles:     toPointer(userRoles),
			UpdatedAt: toPointer(user.Details.ChangeDate.AsTime()),
			UserType:  toPointer(userType),
		}
	}

	return models.Users{
		Count: toPointer(len(users)),
		Items: toPointer(users),
	}, nil
}

// nolint:cyclop,gocognit
func (z *zitadelStore) CreateUser(user models.User) (models.User, error) {
	// Validate
	if user.Name == nil {
		return models.User{}, &common.BadRequestError{Reason: "cannot create for empty Name"}
	}
	if user.Email == nil {
		return models.User{}, &common.BadRequestError{Reason: "cannot create for empty Email"}
	}
	if user.Roles == nil {
		return models.User{}, &common.BadRequestError{Reason: "cannot create for empty Roles"}
	}
	if user.UserType == nil {
		return models.User{}, &common.BadRequestError{Reason: "cannot create for empty UserType"}
	}
	if user.Id != nil {
		return models.User{}, &common.BadRequestError{Reason: "cannot create with Id"}
	}
	if user.Banned != nil && *user.Banned {
		return models.User{}, &common.BadRequestError{Reason: "cannot create with Banned"}
	}

	// Create user
	var userID string
	switch userType := *user.UserType; userType {
	case models.REGULARUSER:
		// Create human user to API
		resp, err := z.mgmt.ImportHumanUser(context.Background(), &management_pb.ImportHumanUserRequest{
			UserName: *user.Email,
			Profile: &management_pb.ImportHumanUserRequest_Profile{
				FirstName:   *user.Name,
				DisplayName: *user.Name,
			},
			Email: &management_pb.ImportHumanUserRequest_Email{
				Email:           *user.Email,
				IsEmailVerified: true,
			},
			PasswordChangeRequired:          false,
			RequestPasswordlessRegistration: true,
			Idps:                            nil,
		})
		if err != nil {
			if isRPCCode(err, codes.AlreadyExists) { // AlreadyExists error
				return models.User{}, types.ErrAlreadyExists
			}
			return models.User{}, fmt.Errorf("failed to create user from API: %w", err) // Generic error
		}
		userID = resp.UserId

	case models.MACHINEUSER:
		// Create machine user to API
		resp, err := z.mgmt.AddMachineUser(context.Background(), &management_pb.AddMachineUserRequest{
			UserName:        *user.Email,
			Name:            *user.Name,
			Description:     "Machine user <%s>",
			AccessTokenType: user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
		})
		if err != nil {
			if isRPCCode(err, codes.AlreadyExists) { // AlreadyExists error
				return models.User{}, types.ErrAlreadyExists
			}
			return models.User{}, fmt.Errorf("failed to create user from API: %w", err) // Generic error
		}
		userID = resp.UserId

	default:
		return models.User{}, &common.BadRequestError{
			Reason: fmt.Sprintf("invalid UserType=%s requested", userType),
		}
	}

	// Set user roles
	err := z.setUserRoles(userID, *user.Roles)
	if err != nil {
		return models.User{}, err
	}

	// Fetch user from id
	return z.GetUser(userID)
}

// nolint:cyclop,gocognit
func (z *zitadelStore) UpdateUser(user models.User) (models.User, error) {
	// Validate
	if user.Id == nil {
		return models.User{}, &common.BadRequestError{Reason: "cannot update for empty Id"}
	}
	if user.UserType == nil {
		return models.User{}, &common.BadRequestError{Reason: "cannot update for empty UserType"}
	}

	// Create user
	switch userType := *user.UserType; userType {
	case models.REGULARUSER:
		// Ban
		if user.Banned != nil && *user.Banned {
			_, err := z.mgmt.LockUser(context.Background(), &management_pb.LockUserRequest{
				Id: *user.Id,
			})
			if err != nil {
				if isRPCCode(err, codes.NotFound) { // NotFound error
					return models.User{}, types.ErrNotFound
				}
				return models.User{}, fmt.Errorf("failed to lock user from API: %w", err) // Generic Error
			}
		}
		// Update roles
		if user.Roles != nil {
			err := z.setUserRoles(*user.Id, *user.Roles)
			if err != nil {
				return models.User{}, err
			}
		}
		// Update name
		if user.Name != nil {
			_, err := z.mgmt.UpdateHumanProfile(context.Background(), &management_pb.UpdateHumanProfileRequest{
				UserId:      *user.Id,
				FirstName:   *user.Name,
				DisplayName: *user.Name,
			})
			if err != nil {
				if isRPCCode(err, codes.NotFound) { // NotFound error
					return models.User{}, types.ErrNotFound
				}
				return models.User{}, fmt.Errorf("failed to udate user profile from API: %w", err) // Generic Error
			}
		}
		// Update email
		if user.Email != nil {
			_, err := z.mgmt.UpdateHumanEmail(context.Background(), &management_pb.UpdateHumanEmailRequest{
				UserId:          *user.Id,
				Email:           *user.Email,
				IsEmailVerified: true,
			})
			if err != nil {
				if isRPCCode(err, codes.NotFound) { // NotFound error
					return models.User{}, types.ErrNotFound
				}
				return models.User{}, fmt.Errorf("failed to update user email from API: %w", err) // Generic Error
			}
		}

	case models.MACHINEUSER:
		return models.User{}, &common.BadRequestError{
			Reason: fmt.Sprintf("invalid UserType=%s requested", userType),
		}

	default:
		return models.User{}, &common.BadRequestError{
			Reason: fmt.Sprintf("invalid UserType=%s requested", userType),
		}
	}

	// Fetch user from id
	return z.GetUser(*user.Id)
}

func (z *zitadelStore) DeleteUser(userID models.UserID) error {
	_, err := z.mgmt.RemoveUser(context.Background(), &management_pb.RemoveUserRequest{
		Id: userID,
	})
	if err != nil {
		if isRPCCode(err, codes.NotFound) { // NotFound error
			return types.ErrNotFound
		}
		return fmt.Errorf("failed to delete user from API: %w", err) // Generic error
	}
	return nil
}

func (z *zitadelStore) GetUserAuth(userID models.UserID) (models.UserAuths, error) {
	var auths []models.UserAuth

	// List access tokens
	tokenResp, err := z.mgmt.ListPersonalAccessTokens(context.Background(), &management_pb.ListPersonalAccessTokensRequest{
		UserId: userID,
		Query:  &object_pb.ListQuery{},
	})
	if err != nil {
		if !isRPCCode(err, codes.NotFound) { // Only pass for NotFound errors
			return models.UserAuths{}, fmt.Errorf("failed to list access tokens from API: %w", err)
		}
	} else {
		for _, token := range tokenResp.Result {
			auths = append(auths, models.UserAuth{
				CredType:  toPointer(models.CredTypeTOKEN),
				CreatedAt: toPointer(token.Details.CreationDate.AsTime()),
				CredID:    toPointer(token.Id),
				UserID:    toPointer(userID),
			})
		}
	}

	// List keys
	keyResp, err := z.mgmt.ListMachineKeys(context.Background(), &management_pb.ListMachineKeysRequest{
		UserId: userID,
		Query: &object_pb.ListQuery{
			Offset: 0,
			Limit:  0,
			Asc:    false,
		},
	})
	if err != nil {
		if !isRPCCode(err, codes.NotFound) { // Only pass for NotFound errors
			return models.UserAuths{}, fmt.Errorf("failed to list machine keys from API: %w", err) // Generic error
		}
	} else {
		for _, key := range keyResp.Result {
			auths = append(auths, models.UserAuth{
				CredType:  toPointer(models.CredTypeFILE),
				CreatedAt: toPointer(key.Details.CreationDate.AsTime()),
				CredID:    toPointer(key.Id),
				UserID:    toPointer(userID),
			})
		}
	}

	return models.UserAuths{
		Count: toPointer(len(auths)),
		Items: toPointer(auths),
	}, nil
}

func (z *zitadelStore) CreateUserAuth(userID models.UserID, credType models.CredentialType, credExpiry *models.CredentialExpiry) (models.UserCred, error) {
	// Validate
	var expiryTimestamp *timestamppb.Timestamp
	if credExpiry != nil {
		if credExpiry.Before(time.Now().Add(1 * time.Minute)) {
			return models.UserCred{}, &common.BadRequestError{Reason: "cannot create with expiry in the past"}
		}
		expiryTimestamp = timestamppb.New(*credExpiry)
	}

	switch credType {
	case models.CredentialTypeTOKEN:
		// If access token requested, create access token
		resp, err := z.mgmt.AddPersonalAccessToken(context.Background(), &management_pb.AddPersonalAccessTokenRequest{
			UserId:         userID,
			ExpirationDate: expiryTimestamp,
		})
		if err != nil {
			if isRPCCode(err, codes.NotFound) { // NotFound error
				return models.UserCred{}, types.ErrNotFound
			}
			return models.UserCred{}, fmt.Errorf("failed to create user access token from API: %w", err) // Generic error
		}
		return models.UserCred{
			Credentials: toPointer(resp.Token),
			UserAuth: &models.UserAuth{
				CredType:  toPointer(models.CredTypeTOKEN),
				CreatedAt: toPointer(resp.Details.CreationDate.AsTime()),
				CredID:    toPointer(resp.TokenId),
				UserID:    toPointer(userID),
			},
		}, nil

	case models.CredentialTypeFILE:
		// If service account requested, create key
		resp, err := z.mgmt.AddMachineKey(context.Background(), &management_pb.AddMachineKeyRequest{
			UserId:         userID,
			Type:           auth_pb.KeyType_KEY_TYPE_JSON, // JSON key
			ExpirationDate: expiryTimestamp,
		})
		if err != nil {
			if isRPCCode(err, codes.NotFound) { // NotFound error
				return models.UserCred{}, types.ErrNotFound
			}
			return models.UserCred{}, fmt.Errorf("failed to create user service account from API: %w", err) // Generic error
		}
		return models.UserCred{
			Credentials: toPointer(string(resp.KeyDetails)),
			UserAuth: &models.UserAuth{
				CredType:  toPointer(models.CredTypeFILE),
				CreatedAt: toPointer(resp.Details.CreationDate.AsTime()),
				CredID:    toPointer(resp.KeyId),
				UserID:    toPointer(userID),
			},
		}, nil

	default:
		return models.UserCred{}, &common.BadRequestError{
			Reason: fmt.Sprintf("invalid CredType=%s requested", credType),
		}
	}
}

func (z *zitadelStore) RevokeUserAuth(userID models.UserID, userAuth models.UserAuth) error {
	if userAuth.CredType == nil {
		return &common.BadRequestError{Reason: "cannot create for empty CredType"}
	}

	switch credType := *userAuth.CredType; credType {
	case models.CredTypeTOKEN:
		// If access token requested, create access token
		_, err := z.mgmt.RemovePersonalAccessToken(context.Background(), &management_pb.RemovePersonalAccessTokenRequest{
			UserId:  userID,
			TokenId: *userAuth.CredID,
		})
		if err != nil {
			if isRPCCode(err, codes.NotFound) { // NotFound error
				return types.ErrNotFound
			}
			return fmt.Errorf("failed to delete user access token from API: %w", err) // Generic error
		}
		return nil

	case models.CredTypeFILE:
		// If service account requested, create key
		_, err := z.mgmt.RemoveMachineKey(context.Background(), &management_pb.RemoveMachineKeyRequest{
			UserId: userID,
			KeyId:  *userAuth.CredID,
		})
		if err != nil {
			if isRPCCode(err, codes.NotFound) { // NotFound error
				return types.ErrNotFound
			}
			return fmt.Errorf("failed to delete user service account from API: %w", err) // Generic error
		}
		return nil

	default:
		return &common.BadRequestError{
			Reason: fmt.Sprintf("invalid CredType=%s requested", credType),
		}
	}
}

func (z *zitadelStore) getUserData(userResp *user_pb.User) (models.User, error) {
	// Get user roles
	userRoles, err := z.getUserRoles(userResp.Id)
	if err != nil {
		return models.User{}, err
	}

	// Get user type
	userType := models.REGULARUSER
	if userResp.GetMachine() != nil {
		userType = models.MACHINEUSER
	}

	return models.User{
		Banned:    toPointer(userResp.State == user_pb.UserState_USER_STATE_LOCKED),
		CreatedAt: toPointer(userResp.Details.CreationDate.AsTime()),
		Email:     toPointer(userResp.UserName),
		Id:        toPointer(userResp.Id),
		Name:      toPointer(userResp.UserName),
		Roles:     toPointer(userRoles),
		UpdatedAt: toPointer(userResp.Details.ChangeDate.AsTime()),
		UserType:  toPointer(userType),
	}, nil
}

func (z *zitadelStore) getUserRoles(userID models.UserID) ([]string, error) {
	// Get user roles from grants
	grantResp, err := z.mgmt.ListUserGrants(context.Background(), &management_pb.ListUserGrantRequest{
		Query: &object_pb.ListQuery{},
		Queries: []*user_pb.UserGrantQuery{
			{
				Query: &user_pb.UserGrantQuery_ProjectIdQuery{
					ProjectIdQuery: &user_pb.UserGrantProjectIDQuery{
						ProjectId: z.projectID,
					},
				},
			},
			{
				Query: &user_pb.UserGrantQuery_UserIdQuery{
					UserIdQuery: &user_pb.UserGrantUserIDQuery{
						UserId: userID,
					},
				},
			},
		},
	})
	if err != nil {
		if isRPCCode(err, codes.NotFound) { // Return empty for NotFound error
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list user grants from API: %w", err) // Generic error
	}

	// Extract user roles from grants
	userRoles := make([]string, 0)
	for _, roleResp := range grantResp.Result {
		userRoles = append(userRoles, roleResp.RoleKeys...)
	}
	return userRoles, nil
}

func (z *zitadelStore) setUserRoles(userID models.UserID, roles []string) error {
	// Get user roles from grants
	grantResp, err := z.mgmt.ListUserGrants(context.Background(), &management_pb.ListUserGrantRequest{
		Query: &object_pb.ListQuery{},
		Queries: []*user_pb.UserGrantQuery{
			{
				Query: &user_pb.UserGrantQuery_ProjectIdQuery{
					ProjectIdQuery: &user_pb.UserGrantProjectIDQuery{
						ProjectId: z.projectID,
					},
				},
			},
			{
				Query: &user_pb.UserGrantQuery_UserIdQuery{
					UserIdQuery: &user_pb.UserGrantUserIDQuery{
						UserId: userID,
					},
				},
			},
		},
	})
	// Extract grant ids from grants
	grantID := make([]string, 0)
	if err == nil { // no error, extract grants
		for _, userRole := range grantResp.Result {
			grantID = append(grantID, userRole.ProjectGrantId)
		}
	} else if !isRPCCode(err, codes.NotFound) { // only pass for NotFound error
		return fmt.Errorf("failed to list user grants from API: %w", err)
	}

	// Delete all role grants if possible
	if len(grantID) > 0 {
		_, err = z.mgmt.BulkRemoveUserGrant(context.Background(), &management_pb.BulkRemoveUserGrantRequest{
			GrantId: grantID,
		})
		if err != nil {
			if !isRPCCode(err, codes.NotFound) { // only pass for NotFound error
				return fmt.Errorf("failed to remove user grants from API: %w", err)
			}
		}
	}

	// Create only one user grant with requested roles
	_, err = z.mgmt.AddUserGrant(context.Background(), &management_pb.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: z.projectID,
		RoleKeys:  roles,
	})
	if err != nil {
		if isRPCCode(err, codes.NotFound) { // NotFound error for user
			return types.ErrNotFound
		}
		return fmt.Errorf("failed to add user grants: %w", err) // Generic error
	}

	return nil
}

func toPointer[T any](obj T) *T {
	return &obj
}

func isRPCCode(err error, code codes.Code) bool {
	rpcErrCode := status.Code(err)
	return code == rpcErrCode
}
