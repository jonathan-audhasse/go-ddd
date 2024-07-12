package errors

// Error IDs
const (
	// Authen error
	UnAuthorized           = ErrorID(1000)
	InvalidAuthProtocol    = ErrorID(1001)
	MissingAuthCredentials = ErrorID(1002)

	// Client Error
	RequiredFieldMissing = ErrorID(2000)
	InvalidField         = ErrorID(2000)
	InvalidFormat        = ErrorID(2001)

	// Repository error
	RepoItemNotFound     = ErrorID(3000)
	RepoItemAlreadyExist = ErrorID(3001)

	// Unknow error
	InternalError = ErrorID(5000)
)
