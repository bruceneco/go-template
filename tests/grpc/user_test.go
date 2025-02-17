package grpc_test

import (
	"context"
	"github.com/bruceneco/go-template/internal/adapters/db/postgres"
	proto "github.com/bruceneco/go-template/internal/adapters/grpc/proto/gen"
	"github.com/bruceneco/go-template/internal/domain/user"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func userSuiteSetup(t testing.TB) *suiteParams {
	t.Helper()
	lc := func(lc fx.Lifecycle, db *postgres.Connection) {
		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				assert.NoError(t, db.Exec("DELETE FROM users").Error)
				assert.NoError(t, db.Exec("DELETE FROM users").Error)
				return nil
			},
		})
	}
	return setup(t, lc)
}

func TestUserGRPC(t *testing.T) {
	t.Parallel()

	p := userSuiteSetup(t)
	ctx := context.Background()
	client := proto.NewUserServiceClient(p.GRPCClient)

	t.Run("CreateUser", func(t *testing.T) {
		t.Parallel()

		t.Run("success", func(t *testing.T) {
			t.Parallel()
			u := genUser()
			req := &proto.CreateUserRequest{
				Name:     u.Name,
				Email:    u.Email,
				Password: u.Password,
			}
			res, err := client.CreateUser(ctx, req)
			require.NoError(t, err)
			dbUser, err := client.GetUser(ctx, &proto.GetUserRequest{Id: res.GetId()})
			require.NoError(t, err)
			require.Equal(t, res, dbUser)
		})

		t.Run("duplicated email", func(t *testing.T) {
			t.Parallel()
			u := genUser()
			req := &proto.CreateUserRequest{
				Name:     u.Name,
				Email:    u.Email,
				Password: u.Password,
			}
			_, err := client.CreateUser(ctx, req)
			require.NoError(t, err)
			_, err = client.CreateUser(ctx, req)
			require.ErrorIs(t, err, status.Error(codes.AlreadyExists, user.ErrEmailAlreadyExists.Error()))
		})
		t.Run("invalid params", func(t *testing.T) {
			t.Parallel()
			req := &proto.CreateUserRequest{
				Name:     "aa",
				Email:    "invalid",
				Password: "aaa",
			}
			_, err := client.CreateUser(ctx, req)
			require.Error(t, err)
		})
	})

	t.Run("GetUser", func(t *testing.T) {
		t.Parallel()
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			u := genUser()
			req := &proto.CreateUserRequest{
				Name:     u.Name,
				Email:    u.Email,
				Password: u.Password,
			}
			res, err := client.CreateUser(ctx, req)
			require.NoError(t, err)
			dbUser, err := client.GetUser(ctx, &proto.GetUserRequest{Id: res.GetId()})
			require.NoError(t, err)
			require.Equal(t, res, dbUser)
		})
		t.Run("invalid id", func(t *testing.T) {
			t.Parallel()
			_, err := client.GetUser(ctx, &proto.GetUserRequest{Id: "aaa"})
			require.Error(t, err)
		})
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			_, err := client.GetUser(ctx, &proto.GetUserRequest{Id: uuid.New().String()})
			require.ErrorIs(t, err, status.Error(codes.NotFound, user.ErrUserNotFound.Error()))
		})
	})

	t.Run("UpdateUser", func(t *testing.T) {
		t.Parallel()
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			u := genUser()
			req := &proto.CreateUserRequest{
				Name:     u.Name,
				Email:    u.Email,
				Password: u.Password,
			}
			res, err := client.CreateUser(ctx, req)
			require.NoError(t, err)
			res.Name = gofakeit.Name()
			updated, err := client.UpdateUser(ctx, res)
			require.NoError(t, err)
			require.Equal(t, res.GetEmail(), updated.GetEmail())
			require.Equal(t, res.GetName(), updated.GetName())
		})
		t.Run("invalid id", func(t *testing.T) {
			t.Parallel()
			_, err := client.UpdateUser(ctx, &proto.User{Id: "aaa"})
			require.Error(t, err)
		})
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			u := genUser()
			_, err := client.UpdateUser(ctx, &proto.User{Id: uuid.New().String(),
				Name:  u.Name,
				Email: u.Email,
			})
			require.ErrorIs(t, err, status.Error(codes.NotFound, user.ErrUserNotFound.Error()))
		})
		t.Run("email already exists", func(t *testing.T) {
			t.Parallel()
			originalUser := genUser()
			res, err := client.CreateUser(ctx, &proto.CreateUserRequest{
				Name:     originalUser.Name,
				Email:    originalUser.Email,
				Password: originalUser.Password,
			})
			require.NoError(t, err)
			u := genUser()
			_, err = client.CreateUser(ctx, &proto.CreateUserRequest{
				Name:     u.Name,
				Email:    u.Email,
				Password: u.Password,
			})
			require.NoError(t, err)
			res.Email = u.Email
			_, err = client.UpdateUser(ctx, res)
			require.ErrorIs(t, err, status.Error(codes.AlreadyExists, user.ErrEmailAlreadyExists.Error()))
		})
	})

	t.Run("DeleteUser", func(t *testing.T) {
		t.Parallel()
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			u := genUser()
			createdUser, err := client.CreateUser(ctx, &proto.CreateUserRequest{
				Name:     u.Name,
				Email:    u.Email,
				Password: u.Password,
			})
			require.NoError(t, err)
			_, err = client.DeleteUser(ctx, &proto.DeleteUserRequest{Id: createdUser.GetId()})
			require.NoError(t, err)
			_, err = client.GetUser(ctx, &proto.GetUserRequest{Id: createdUser.GetId()})
			require.ErrorIs(t, err, status.Error(codes.NotFound, user.ErrUserNotFound.Error()))
		})
		t.Run("invalid id", func(t *testing.T) {
			t.Parallel()
			_, err := client.DeleteUser(ctx, &proto.DeleteUserRequest{Id: "aaa"})
			require.Error(t, err)
		})
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			_, err := client.DeleteUser(ctx, &proto.DeleteUserRequest{Id: uuid.New().String()})
			require.ErrorIs(t, err, status.Error(codes.NotFound, user.ErrUserNotFound.Error()))
		})
	})
}

func genUser() user.Entity {
	return user.Entity{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 10),
	}
}
