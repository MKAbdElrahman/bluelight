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
}

type PermissionsService struct {
	permissionsRepository PermissionsRepository
}

func NewPermissionsService(r PermissionsRepository) *PermissionsService {
	return &PermissionsService{
		permissionsRepository: r,
	}
}

func (svc *PermissionsService) GetAllForUser(userID int64) (Permissions, error) {
	return svc.permissionsRepository.GetAllForUser(userID)
}
