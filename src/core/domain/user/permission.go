package user

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionsRepository interface {
	GetAllForUser(userId int64) (Permissions, error)
	AddForUser(userId int64, codes ...string) error
}
