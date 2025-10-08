package common

func HasPermission(userRole string, accessibleRole []string) bool {
	for _, role := range accessibleRole {
		if userRole == role {
			return true
		}
	}
	return false
}
