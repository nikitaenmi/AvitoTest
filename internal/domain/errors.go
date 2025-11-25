package domain

type ErrorType string

const (
	ErrorTypeTeamExists  ErrorType = "TEAM_EXISTS"
	ErrorTypePRExists    ErrorType = "PR_EXISTS"
	ErrorTypePRMerged    ErrorType = "PR_MERGED"
	ErrorTypeNotAssigned ErrorType = "NOT_ASSIGNED"
	ErrorTypeNoCandidate ErrorType = "NO_CANDIDATE"
	ErrorTypeNotFound    ErrorType = "NOT_FOUND"
	ErrorTypeValidation  ErrorType = "VALIDATION_ERROR"
)

type DomainError struct {
	Type    ErrorType
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(errorType ErrorType, message string) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
	}
}

func NewTeamExistsError() *DomainError {
	return &DomainError{
		Type:    ErrorTypeTeamExists,
		Message: "team_name already exists",
	}
}

func NewPRExistsError() *DomainError {
	return &DomainError{
		Type:    ErrorTypePRExists,
		Message: "PR id already exists",
	}
}

func NewPRMergedError() *DomainError {
	return &DomainError{
		Type:    ErrorTypePRMerged,
		Message: "cannot reassign on merged PR",
	}
}

func NewNotAssignedError() *DomainError {
	return &DomainError{
		Type:    ErrorTypeNotAssigned,
		Message: "reviewer is not assigned to this PR",
	}
}

func NewNoCandidateError() *DomainError {
	return &DomainError{
		Type:    ErrorTypeNoCandidate,
		Message: "no active replacement candidate in team",
	}
}

func NewNotFoundError(resource string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeNotFound,
		Message: resource + " not found",
	}
}

func NewValidationError(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}
