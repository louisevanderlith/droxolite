package roletype

import (
	"strings"
)

type Enum = int

const (
	Admin Enum = iota
	Owner
	User
	Unknown
)

/*
Role Types are assigned depending on how the Users are created.
Users permissions can vary between applications, this allows a less complex system.
Additional permissions can be added via a 'Request' basis, but this is still under development.

Admin: Admins are only created via seed files. As this is the only way to assign a custom 'roletype'.
Owner: Owners are created by Service.Users OR by Registering FOR Cars.APP
User: ...
Unknown: Non-logged-in sessions are Unknown. Even the unknown is someone.
*/

var roleTypes = [...]string{
	"Admin",
	"Owner",
	"User",
	"Unknown"}

func GetRoleType(name string) Enum {
	var result Enum

	for k, v := range roleTypes {
		if strings.ToUpper(name) == v {
			result = Enum(k)
			break
		}
	}

	return result
}
