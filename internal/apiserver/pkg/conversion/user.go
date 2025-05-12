package conversion

import (
	"github.com/onexstack/onexstack/pkg/core"

	"github.com/onexstack/fastgo/internal/apiserver/model"

	apiv1 "github.com/onexstack/fastgo/pkg/api/apiserver/v1"
)

func UserodelToUserV1(userModel *model.User) *apiv1.User {
	var protoUser apiv1.User
	_ = core.CopyWithConverters(&protoUser, userModel)
	return &protoUser
}

func UserV1ToUserodel(protoUser *apiv1.User) *model.User {
	var userModel model.User
	_ = core.CopyWithConverters(&userModel, protoUser)
	return &userModel
}
