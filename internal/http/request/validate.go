package request

// Validator is an interface for DTOs that need validation.
type Validator interface {
	Validate() error
}
