package api

// UserDetails is an interface describing basic behaviour expected from user
// details implementation.
// Everything KYC related should be implemented in `UserDetails` interface
// to avoid ugly user type switches everywhere.
type UserDetails interface {
	// Validate takes into account only KYC entities values, checking if everything
	// is submitted and passes constraint checks.
	// Returns validation error or nil if all good
	Validate() error
	// RequiredDocuments return list of documents required to pass KYC check
	RequiredDocuments() []RequiredDocument
	// DisplayName return user name representation if possible
	DisplayName() *string
}
