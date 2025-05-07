package errorsx

import (
	"errors"
	"fmt"
)

// ErrorX定义了 错误类型, 用于描述错误的详细信息.
type ErrorX struct {
	// Code 表示错误的HTTP状态码.
	Code int `json:"code,omitempty"`
	// Reason 表示错误发生的原因, 通常为业务错误码
	Reason string `json:"reason,omitempty"`
	// Message 表示简短的错误信息，通常可直接暴露给用户查看.
	Message string `json:"message,omitempty"`
}

// New 创建一个新的 ErrorX 实例。
//
// 参数:
//
//	code int: 错误代码
//	reason string: 错误原因
//	format string: 错误信息格式
//	args ...any: 错误信息参数列表
//
// 返回值:
//
//	*ErrorX: 创建的 ErrorX 实例
func New(code int, reason string, format string, args ...any) *ErrorX {
	return &ErrorX{
		Code:    code,
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
	}
}

// Error 方法实现了 error 接口，返回自定义错误信息的字符串表示
// 它将 ErrorX 结构体的 Code, Reason 和 Message 三个字段组合成一个字符串并返回
func (err *ErrorX) Error() string {
	return fmt.Sprintf("error: code = %d, reason = %s, message = %s", err.Code, err.Reason, err.Message)
}

// WithMessage 向ErrorX结构体添加一条自定义的错误信息
//
// 参数：
// - format string: 格式化字符串，用于生成错误信息
// - args ...any: 可变参数，用于格式化错误信息
//
// 返回值：
// - *ErrorX: 返回当前ErrorX实例的指针
func (err *ErrorX) WithMessage(format string, args ...any) *ErrorX {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

// FromError 将给定的错误 err 转换为 ErrorX 类型的指针
//
// 参数:
//
//	err: 需要转换的错误对象
//
// 返回值:
//
//	*ErrorX: 如果 err 是 ErrorX 类型，返回其指针；
//	         如果 err 是 nil，返回 nil；
//	         否则，返回一个新的 ErrorX 对象，其状态码为 ErrInternal.Code，
//	         错误原因为 ErrInternal.Reason，错误信息为 err.Error()
func FromError(err error) *ErrorX {
	// 如果传入的错误为nil，则返回nil
	if err == nil {
		return nil
	}

	// 尝试将传入的错误转换为ErrorX类型
	if errx := new(ErrorX); errors.As(err, &errx) {
		// 如果转换成功，则返回转换后的ErrorX对象
		return errx
	}

	// 如果转换失败，则创建一个新的ErrorX对象并返回
	// 使用ErrInternal的错误码和错误原因，以及传入的错误消息
	return New(ErrInternal.Code, ErrInternal.Reason, err.Error())
}
