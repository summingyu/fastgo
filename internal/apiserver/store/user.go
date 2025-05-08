package store

import (
	"context"
	"errors"
	"log/slog"

	"github.com/onexstack/fastgo/internal/apiserver/model"
	"github.com/onexstack/fastgo/internal/pkg/errorsx"
	"github.com/onexstack/onexstack/pkg/store/where"
	"gorm.io/gorm"
)

type UserStore interface {
	Create(ctx context.Context, obj *model.User) error
	Update(ctx context.Context, obj *model.User) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.User, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.User, error)

	UserExpansion
}

type UserExpansion interface{}

type userStore struct {
	store *datastore
}

var _ UserStore = (*userStore)(nil)

// newUserStore 创建一个新的 userStore 实例，并返回一个指向该实例的指针
//
// 参数：
//   - store: 一个指向datastore的指针，该datastore用于存储用户数据
//
// 返回值：
//   - *userStore: 返回一个指向新创建的 userStore 实例的指针
func newUserStore(store *datastore) *userStore {
	return &userStore{
		store: store,
	}
}

// Delete implements UserStore.
// Delete 从数据库中删除用户
//
// 参数:
// ctx - 上下文对象
// opts - 删除条件选项
//
// 返回值:
// error - 如果删除失败，返回错误信息；如果删除成功，返回 nil
func (s *userStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx, opts).Delete(new(model.User)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("Failed to delete user from databases", "err", err, "conditions", opts)
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

// Get implements UserStore.
// Get 方法从数据库中获取一个符合给定条件的用户
//
// 参数：
//
//	ctx - 上下文对象，包含请求相关的信息
//	opts - 查询条件
//
// 返回值：
//
//	*model.User - 找到的用户对象指针
//	error - 错误信息，如果发生错误则返回，否则返回 nil
func (s *userStore) Get(ctx context.Context, opts *where.Options) (*model.User, error) {
	var obj model.User
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		slog.Error("Failed to get user from databases", "err", err, "conditions", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotFound
		}
		return nil, errorsx.ErrDBRead.WithMessage(err.Error())
	}
	return &obj, nil
}

// List implements UserStore.
// List 从数据库中获取用户列表
// 参数：
// ctx：上下文对象
// opts：查询条件对象
// 返回值：
// count：用户总数
// ret：用户列表
// err：错误信息
func (s *userStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.User, err error) {
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		slog.Error("Failed to list users from databases", "err", err, "conditions", opts)
		err = errorsx.ErrDBRead.WithMessage(err.Error())
	}
	return
}

// Update implements UserStore.
// Update 在给定的上下文中更新用户信息
//
// 参数：
//
//	ctx: 上下文对象，用于传递请求范围内的数据，如截止日期、取消信号等
//	obj: 指向model.User类型的指针，表示要更新的用户信息
//
// 返回值：
//
//	如果更新成功，返回nil；否则返回一个错误对象
func (s *userStore) Update(ctx context.Context, obj *model.User) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		slog.Error("Failed to update user in databases", "err", err, "user", obj)
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

// Create 函数用于将用户对象插入到数据库中
//
// 参数：
// ctx - 上下文对象
// obj - 待插入的用户对象
//
// 返回值：
// error - 如果插入失败，则返回错误信息；否则返回nil
func (s *userStore) Create(ctx context.Context, obj *model.User) error {
	if err := s.store.DB(ctx).Create(&obj).Error; err != nil {
		slog.Error("Failed to insert user into databases", "err", err, "user", obj)
		return errorsx.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}
