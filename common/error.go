package common

// ServiceError 业务异常
type ServiceError struct {
	Message string			// 异常消息
	Cause error				// 捕获的异常信息
	HttpStatus int			// http状态码
	Response *Response		// 响应客户端的数据
}
func (s ServiceError) Error () string {
	return s.Message
}

// SetCause 设置Catch的异常
func (s ServiceError) SetCause(err error) *ServiceError {
	s.Cause = err
	return &s
}

// NewServiceError 创建新的业务异常，指定原始异常，状态码，消息
func NewServiceError (code *Code) *ServiceError {
	return &ServiceError{
		Message: code.Message,
		Cause: nil,
		HttpStatus: code.HttpStatus,
		Response: &Response{
			Success: false,
			Code:    code,
			Message: code.Message,
		},
	}
}

