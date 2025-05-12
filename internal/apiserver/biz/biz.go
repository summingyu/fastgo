package biz

import (
	postv1 "github.com/onexstack/fastgo/internal/apiserver/biz/v1/post"
	userv1 "github.com/onexstack/fastgo/internal/apiserver/biz/v1/user"
	"github.com/onexstack/fastgo/internal/apiserver/store"
)

type IBiz interface {
	UserV1() userv1.UserBiz
	PostV1() postv1.PostBiz
}

type biz struct {
	store store.IStore
}

var _ IBiz = (*biz)(nil)

func NewBiz(store store.IStore) IBiz {
	return &biz{
		store: store,
	}
}

// PostV1 implements IBiz.
func (b *biz) PostV1() postv1.PostBiz {
	return postv1.New(b.store)
}

// UserV1 implements IBiz.
func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store)
}
