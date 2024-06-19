package errors

// Error IDs
const (
	// Authen error
	UnAuthorized = ErrorID(1000)

	// General Error
	RequiredFieldMissing = ErrorID(2000)
	InvalidField         = ErrorID(2001)
	InvalidFormat        = ErrorID(2002)

	// Repository error
	RepoItemNotFound     = ErrorID(3000)
	RepoItemAlreadyExist = ErrorID(3001)

	// Unknow error
	InternalError = ErrorID(5000)
)
