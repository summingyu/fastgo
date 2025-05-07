package rid

import "github.com/onexstack/onexstack/pkg/id"

const defaultABC = "abcdefghijklmnopqrstuvwxyz1234567890"

type ResourceID string

const (
	UserID ResourceID = "user"
	PostID ResourceID = "post"
)

// String 方法将 ResourceID 类型的值转换为字符串并返回。
func (rid ResourceID) String() string {
	return string(rid)
}

// New 方法为 ResourceID 类型生成一个新的字符串表示形式。
// 它接收一个 counter 参数，用于生成唯一的字符串。
// counter 是一个 uint64 类型的值。
//
// 该方法使用 id.NewCode 函数生成唯一的字符串，其中：
// - counter：用于生成唯一字符串的计数器。
// - id.WithCodeChars([]rune(defaultABC))：设置生成唯一字符串时使用的字符集，这里使用默认的 ABC 字符集。
// - id.WithCodeL(6)：设置生成唯一字符串的长度为 6。
// - id.WithCodeSalt(Salt())：设置生成唯一字符串时使用的盐值，通过调用 Salt() 函数获取。
//
// 返回值为一个由 ResourceID 类型和生成的唯一字符串组成的字符串，两者通过 "-" 连接。
func (rid ResourceID) New(counter uint64) string {
	uniqueStr := id.NewCode(
		counter,
		id.WithCodeChars([]rune(defaultABC)),
		id.WithCodeL(6),
		id.WithCodeSalt(Salt()),
	)
	return rid.String() + "-" + uniqueStr
}
