package store

import (
	"context"
	"sync"

	"github.com/onexstack/onexstack/pkg/store/where"
	"gorm.io/gorm"
)

var (
	once sync.Once
	S    *datastore
)

type IStore interface {
	DB(ctx context.Context, wheres ...where.Where) *gorm.DB
	TX(ctx context.Context, fn func(ctx context.Context) error) error

	User() UserStore
	Post() PostStore
}

type transactionKey struct{}

type datastore struct {
	core *gorm.DB

	// fake *gorm.DB
}

var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB) *datastore {
	once.Do(func() {
		S = &datastore{db}
	})

	return S
}

// DB 方法从上下文获取数据库连接，并应用一系列过滤条件
//
// 参数:
// ctx - 上下文，用于传递事务等上下文信息
// wheres - 可变参数，表示一系列的过滤条件
//
// 返回值:
// *gorm.DB - 返回应用了所有过滤条件后的数据库连接对象
func (store *datastore) DB(ctx context.Context, wheres ...where.Where) *gorm.DB {
	db := store.core
	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		db = tx
	}

	for _, whr := range wheres {
		db = whr.Where(db)
	}
	return db
}

// TX 方法启动一个事务，并在事务上下文中执行提供的函数。
//
// 参数：
//
//	ctx: 上下文，用于控制整个操作的生命周期。
//	fn: 要在事务中执行的函数，该函数应接受一个上下文并返回一个错误。
//
// 返回值：
//
//	如果在事务执行期间发生错误，则返回该错误；否则返回 nil。
func (store *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return store.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

// User 方法返回与 datastore 关联的 UserStore 实例。
func (store *datastore) User() UserStore {
	return newUserStore(store)
}

// Post 方法返回一个新的 PostStore 实例
// 参数：
//   - store: 数据存储实例指针
//
// 返回值：
//   - PostStore: 返回一个 PostStore 实例
func (store *datastore) Post() PostStore {
	return newPostStore(store)
}
