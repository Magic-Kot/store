package errcodes

import (
	"git.appkode.ru/pub/go/failure"
)

type ErrCode failure.ErrorCode

const (
	InternalServerError = failure.ErrorCode("INTERNAL_SERVER_ERROR")
	ValidationError     = failure.ErrorCode("VALIDATION_ERROR")
	RefreshTokenInvalid = failure.ErrorCode("REFRESH_TOKEN_INVALID")
	RefreshTokenExpired = failure.ErrorCode("REFRESH_TOKEN_EXPIRED")
	AccessTokenInvalid  = failure.ErrorCode("ACCESS_TOKEN_INVALID")
	AccessTokenExpired  = failure.ErrorCode("ACCESS_TOKEN_EXPIRED")
	NotFound            = failure.ErrorCode("NOT_FOUND")
	Forbidden           = failure.ErrorCode("FORBIDDEN")
	AlreadyExists       = failure.ErrorCode("FILE_ALREADY_EXISTS")
	ReportNotFound      = failure.ErrorCode("REPORT_NOT_FOUND")
	UserNotFound        = failure.ErrorCode("USER_NOT_FOUND")
	DirectoriesBusy     = failure.ErrorCode("DIRECTORIES_BUSY")
)

const (
	InternalServerErrorMessage     = "internal server error"
	ValidationErrorMessage         = "validation error"
	NotFoundErrorMessage           = "not found"
	ForbiddenErrorMessage          = "forbidden"
	RefreshTokenInvalidMessage     = "refresh token invalid"
	RefreshTokenExpiredMessage     = "refresh token expired"
	AccessTokenInvalidMessage      = failure.ErrorCode("access_token_invalid")
	AccessTokenExpiredMessage      = failure.ErrorCode("access_token_expired")
	ErrorFromRemoteServerMessage   = failure.ErrorCode("error_from_remote_server")
	ReportNotFoundMessage          = failure.ErrorCode("report_not_found")
	UserNotFoundMessage            = failure.ErrorCode("user_not_found")
	PushTokenAlreadyDeletedMessage = "Push token already deleted"
)
