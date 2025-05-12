package user

import (
	"context"
	"log/slog"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/onexstack/fastgo/internal/apiserver/model"
	"github.com/onexstack/fastgo/internal/apiserver/pkg/conversion"
	"github.com/onexstack/fastgo/internal/apiserver/store"
	"github.com/onexstack/fastgo/internal/pkg/contextx"
	"github.com/onexstack/fastgo/internal/pkg/known"
	apiv1 "github.com/onexstack/fastgo/pkg/api/apiserver/v1"
	"github.com/onexstack/onexstack/pkg/store/where"
	"golang.org/x/sync/errgroup"
)

type UserBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error)
	Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error)
	List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)

	UserExpansion
}

type UserExpansion interface {
}

type userBiz struct {
	store store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(store store.IStore) UserBiz {
	return &userBiz{
		store: store,
	}
}

// Create implements UserBiz.
func (b *userBiz) Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	var userM model.User
	_ = copier.Copy(&userM, rq)

	if err := b.store.User().Create(ctx, &userM); err != nil {
		return nil, err
	}

	return &apiv1.CreateUserResponse{UserID: userM.UserID}, nil
}

// Delete implements UserBiz.
func (b *userBiz) Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	if err := b.store.User().Delete(ctx, where.F("userID", contextx.UserID(ctx))); err != nil {
		return nil, err
	}
	return &apiv1.DeleteUserResponse{}, nil
}

// Get implements UserBiz.
func (b *userBiz) Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.F("userID", contextx.UserID(ctx)))
	if err != nil {
		return nil, err
	}
	return &apiv1.GetUserResponse{User: conversion.UserodelToUserV1(userM)}, nil
}

// List implements UserBiz.
func (b *userBiz) List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	whr := where.P(int(rq.Offset), int(rq.Limit))
	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	// 设置最大并发数为常量 MaxConcurrency
	eg.SetLimit(known.MaxErrGroupConcurrency)

	// 使用 goroutine 并发处理
	for _, user := range userList {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				count, _, err := b.store.Post().List(ctx, where.T(ctx))
				if err != nil {
					return err
				}
				converted := conversion.UserodelToUserV1(user)
				converted.PostCount = count
				m.Store(user.ID, converted)
				return nil
			}
		})
	}
	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, "Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	users := make([]*apiv1.User, 0, len(userList))
	for _, item := range userList {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*apiv1.User))
	}
	slog.DebugContext(ctx, "Get users from backend storage", "count", len(users))

	return &apiv1.ListUserResponse{TotalCount: count, Users: users}, nil
}

// Update implements UserBiz.
func (b *userBiz) Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}
	if rq.Username != nil {
		userM.Username = *rq.Username
	}
	if rq.Email != nil {
		userM.Email = *rq.Email
	}
	if rq.Nickname != nil {
		userM.Nickname = *rq.Nickname
	}
	if rq.Phone != nil {
		userM.Phone = *rq.Phone
	}
	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}
	return &apiv1.UpdateUserResponse{}, nil
}
