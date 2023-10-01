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
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/types"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	auth_pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/authn"
	management_pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	object_pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	user_pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
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
		return nil, err
	}

	return &zitadelStore{
		projectID: config.ProjectID,
		mgmt:      mgmtClient,
	}, nil
}

func (z *zitadelStore) GetUserFromInfo(info *types.UserInfo) (models.User, error) {
	// Get user from email
	var err error
	var user *user_pb.User
	if info.FromZitadelOIDC {
		userResp, err := z.mgmt.GetUserByID(context.Background(), &management_pb.GetUserByIDRequest{
			Id: info.GetSubject(),
		})
		if err != nil {
			return models.User{}, err
		}
		user = userResp.User
	} else if info.FromGenericOIDC {
		// TODO: Implement create from a generic store
		return models.User{}, fmt.Errorf("not implemented current user fetcher for generic OIDC")
	}
	if err != nil {
		return models.User{}, nil
	}

	// Get user roles
	return z.getUserData(user)
}

func (z *zitadelStore) GetUser(userID models.UserID) (models.User, error) {
	// Get user
	userResp, err := z.mgmt.GetUserByID(context.Background(), &management_pb.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		return models.User{}, err
	}

	// Get user roles
	return z.getUserData(userResp.User)
}

func (z *zitadelStore) GetUsers(params models.GetUsersParams) (models.Users, error) {
	paramSkip := 0
	if params.Skip != nil {
		paramSkip = *params.Skip
	}
	paramTop := 100
	if params.Top != nil {
		paramTop = *params.Top
	}

	// Get users
	userResp, err := z.mgmt.ListUsers(context.Background(), &management_pb.ListUsersRequest{
		Query: &object_pb.ListQuery{
			Offset: uint64(paramSkip),
			Limit:  uint32(paramTop),
		},
	})
	if err != nil {
		return models.Users{}, err
	}

	// Get roles
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
		return models.Users{}, err
	}
	userRoleMap := make(map[string][]string)
	for _, grant := range grantResp.Result {
		currRoles, _ := userRoleMap[grant.UserId]
		userRoleMap[grant.UserId] = append(currRoles, grant.RoleKeys...)
	}

	// Generate user data
	var users []models.User
	for _, user := range userResp.Result {
		// Get user data
		userRoles, _ := userRoleMap[user.Id]
		userType := models.REGULARUSER
		if user.GetMachine() != nil {
			userType = models.MACHINEUSER
		}

		users = append(users, models.User{
			Banned:    toPointer(user.State == user_pb.UserState_USER_STATE_LOCKED),
			CreatedAt: toPointer(user.Details.CreationDate.AsTime()),
			Email:     toPointer(user.UserName),
			Id:        toPointer(user.Id),
			Name:      toPointer(user.UserName),
			Roles:     toPointer(userRoles),
			UpdatedAt: toPointer(user.Details.ChangeDate.AsTime()),
			UserType:  toPointer(userType),
		})
	}

	return models.Users{
		Count: toPointer(len(users)),
		Items: toPointer(users),
	}, nil
}

func (z *zitadelStore) CreateUser(user models.User) (models.User, error) {
	// Validate
	if user.Name == nil {
		return models.User{}, fmt.Errorf("cannot create for empty Name")
	}
	if user.Email == nil {
		return models.User{}, fmt.Errorf("cannot create for empty Email")
	}
	if user.Roles == nil {
		return models.User{}, fmt.Errorf("cannot create for empty Roles")
	}
	if user.UserType == nil {
		return models.User{}, fmt.Errorf("cannot create for empty UserType")
	}
	if user.Banned != nil && !*user.Banned {
		return models.User{}, fmt.Errorf("cannot create for Banned")
	}

	// Create user
	var userId string
	userType := *user.UserType
	switch userType {
	case models.REGULARUSER:
		// Create human user
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
			return models.User{}, err
		}
		userId = resp.UserId

	case models.MACHINEUSER:
		// Create machine user
		resp, err := z.mgmt.AddMachineUser(context.Background(), &management_pb.AddMachineUserRequest{
			UserName:        *user.Email,
			Name:            *user.Name,
			Description:     "Machine user <%s>",
			AccessTokenType: user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
		})
		if err != nil {
			return models.User{}, err
		}
		userId = resp.UserId
	}

	// Set user roles
	err := z.setUserRoles(userId, *user.Roles)
	if err != nil {
		return models.User{}, err
	}

	// Fetch user from id
	return z.GetUser(userId)
}

func (z *zitadelStore) UpdateUser(user models.User) (models.User, error) {
	// Validate
	if user.Id == nil {
		return models.User{}, fmt.Errorf("cannot create for empty Id")
	}
	if user.UserType == nil {
		return models.User{}, fmt.Errorf("cannot create for empty UserType")
	}

	// Create user
	userType := *user.UserType
	switch userType {
	case models.REGULARUSER:
		// Update name
		if user.Name != nil {
			_, err := z.mgmt.UpdateHumanProfile(context.Background(), &management_pb.UpdateHumanProfileRequest{
				UserId:      *user.Id,
				FirstName:   *user.Name,
				DisplayName: *user.Name,
			})
			if err != nil {
				return models.User{}, err
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
				return models.User{}, err
			}
		}
		// Update roles
		if user.Roles != nil {
			err := z.setUserRoles(*user.Id, *user.Roles)
			if err != nil {
				return models.User{}, err
			}
		}
		// Ban
		if user.Banned != nil {
			_, err := z.mgmt.LockUser(context.Background(), &management_pb.LockUserRequest{
				Id: *user.Id,
			})
			if err != nil {
				return models.User{}, err
			}
		}

	case models.MACHINEUSER:
		return models.User{}, fmt.Errorf("machine user cannot be updated")
	}

	// Fetch user from id
	return z.GetUser(*user.Id)
}

func (z *zitadelStore) DeleteUser(userID models.UserID) error {
	_, err := z.mgmt.RemoveUser(context.Background(), &management_pb.RemoveUserRequest{
		Id: userID,
	})
	return err
}

func (z *zitadelStore) GetUserAuth(userID models.UserID) (models.UserAuths, error) {
	var auths []models.UserAuth

	// List access tokens
	tokenResp, err := z.mgmt.ListPersonalAccessTokens(context.Background(), &management_pb.ListPersonalAccessTokensRequest{
		UserId: userID,
		Query: &object_pb.ListQuery{
			Offset: 0,
			Limit:  0,
			Asc:    false,
		},
	})
	if err != nil {
		return models.UserAuths{}, err
	}
	for _, token := range tokenResp.Result {
		auths = append(auths, models.UserAuth{
			AuthType:   toPointer(models.AuthTypeACCESSTOKEN),
			CreatedAt:  toPointer(token.Details.CreationDate.AsTime()),
			UserAuthID: toPointer(token.Id),
			UserID:     toPointer(userID),
		})
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
		return models.UserAuths{}, err
	}
	for _, key := range keyResp.Result {
		auths = append(auths, models.UserAuth{
			AuthType:   toPointer(models.AuthTypeSERVICEACCOUNT),
			CreatedAt:  toPointer(key.Details.CreationDate.AsTime()),
			UserAuthID: toPointer(key.Id),
			UserID:     toPointer(userID),
		})
	}

	return models.UserAuths{
		Count: toPointer(len(auths)),
		Items: toPointer(auths),
	}, nil
}

func (z *zitadelStore) CreateUserAuth(userID models.UserID, authType models.AuthType, expiryDate *time.Time) (models.UserCred, error) {
	var expiryTimestamp *timestamppb.Timestamp
	if expiryDate != nil {
		expiryTimestamp = timestamppb.New(*expiryDate)
	}

	switch authType {
	case models.AuthTypeACCESSTOKEN:
		// If access token requested, create access token
		resp, err := z.mgmt.AddPersonalAccessToken(context.Background(), &management_pb.AddPersonalAccessTokenRequest{
			UserId:         userID,
			ExpirationDate: expiryTimestamp,
		})
		if err != nil {
			return models.UserCred{}, err
		}
		return models.UserCred{
			Credentials: toPointer(map[string]interface{}{
				"token": resp.Token,
			}),
			UserAuth: &models.UserAuth{
				AuthType:   toPointer(authType),
				CreatedAt:  toPointer(resp.Details.CreationDate.AsTime()),
				UserAuthID: toPointer(resp.TokenId),
				UserID:     toPointer(userID),
			},
		}, nil

	case models.AuthTypeSERVICEACCOUNT:
		// If service account requested, create key
		resp, err := z.mgmt.AddMachineKey(context.Background(), &management_pb.AddMachineKeyRequest{
			UserId:         userID,
			Type:           auth_pb.KeyType_KEY_TYPE_JSON, // JSON key
			ExpirationDate: expiryTimestamp,
		})
		if err != nil {
			return models.UserCred{}, err
		}
		return models.UserCred{
			Credentials: toPointer(map[string]interface{}{
				"key": string(resp.KeyDetails),
			}),
			UserAuth: &models.UserAuth{
				AuthType:   toPointer(authType),
				CreatedAt:  toPointer(resp.Details.CreationDate.AsTime()),
				UserAuthID: toPointer(resp.KeyId),
				UserID:     toPointer(userID),
			},
		}, nil

	default:
		// Otherwise, unknown type requested
		return models.UserCred{}, fmt.Errorf("invalid AuthType=%s requested", authType)
	}
}

func (z *zitadelStore) RevokeUserAuth(userID models.UserID, userAuth models.UserAuth) error {
	if userAuth.AuthType == nil {
		return fmt.Errorf("cannot create for empty AuthType")
	}

	authType := *userAuth.AuthType
	switch authType {
	case models.AuthTypeACCESSTOKEN:
		// If access token requested, create access token
		_, err := z.mgmt.RemovePersonalAccessToken(context.Background(), &management_pb.RemovePersonalAccessTokenRequest{
			UserId:  userID,
			TokenId: *userAuth.UserAuthID,
		})
		return err

	case models.AuthTypeSERVICEACCOUNT:
		// If service account requested, create key
		_, err := z.mgmt.RemoveMachineKey(context.Background(), &management_pb.RemoveMachineKeyRequest{
			UserId: userID,
			KeyId:  *userAuth.UserAuthID,
		})
		return err

	default:
		// Otherwise, unknown type requested
		return fmt.Errorf("invalid AuthType=%s requested", authType)
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

func (z *zitadelStore) getUserRoles(userId models.UserID) ([]string, error) {
	// Get user roles
	userRoleResp, err := z.mgmt.ListUserGrants(context.Background(), &management_pb.ListUserGrantRequest{
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
			{
				Query: &user_pb.UserGrantQuery_UserIdQuery{
					UserIdQuery: &user_pb.UserGrantUserIDQuery{
						UserId: userId,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	userRoles := make([]string, 0)
	for _, roleResp := range userRoleResp.Result {
		userRoles = append(userRoles, roleResp.RoleKeys...)
	}
	return userRoles, nil
}

func (z *zitadelStore) setUserRoles(userId models.UserID, roles []string) error {
	// Fetch all role grants
	userRolesResp, err := z.mgmt.ListUserGrants(context.Background(), &management_pb.ListUserGrantRequest{
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
			{
				Query: &user_pb.UserGrantQuery_UserIdQuery{
					UserIdQuery: &user_pb.UserGrantUserIDQuery{
						UserId: userId,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	grantId := make([]string, 0)
	for _, userRole := range userRolesResp.Result {
		grantId = append(grantId, userRole.ProjectGrantId)
	}

	// Delete all role grants
	_, err = z.mgmt.BulkRemoveUserGrant(context.Background(), &management_pb.BulkRemoveUserGrantRequest{
		GrantId: grantId,
	})
	if err != nil {
		return err
	}

	// Create only one user grant with requested roles
	_, err = z.mgmt.AddUserGrant(context.Background(), &management_pb.AddUserGrantRequest{
		UserId:    userId,
		ProjectId: z.projectID,
		RoleKeys:  roles,
	})
	return err
}

func toPointer[T any](obj T) *T {
	return &obj
}
