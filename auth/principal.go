package auth


import (
	"github.com/mindhash/goFlow/base"
)



// A Principal is an abstract object that can have access to activities
type Principal interface {
	// The Principal's identifier.
	Name() string
	
	validate() error
	
}

// Role is basically the same as Principal, just concrete. Users can inherit permissions from Roles.
type Role interface {
	Principal
}

// A User is a Principal that can log in and have multiple Roles.
type User interface {
	Principal

	// The user's email address.
	Email() string

	// Sets the user's email address.
	SetEmail(string) error

	// If true, the user is unable to authenticate.
	Disabled() bool

	// Sets the disabled property
	SetDisabled(bool)

	// Authenticates the user's password.
	Authenticate(password string) bool

	// Changes the user's password.
	SetPassword(password string)

	// The set of Roles the user belongs to  
	RoleNames() base.Set

	// The roles the user was explicitly granted access to thru the admin API.
	ExplicitRoles() base.Set

	// Sets the explicit roles the user belongs to.
	SetExplicitRoles(base.Set)  
 
}
