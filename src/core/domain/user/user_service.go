package user

type UserService struct {
	userRepository UserRepositoty
}

func NewUserService(r UserRepositoty) *UserService {
	return &UserService{
		userRepository: r,
	}
}

func (svc *UserService) CreateUser(u *User) error {
	return svc.userRepository.Create(u)
}

func (svc *UserService) UpdateUser(u *User) error {
	return svc.userRepository.Update(u)
}

func (svc *UserService) GetByEmail(email string) (*User, error) {
	return svc.userRepository.GetByEmail(email)
}
