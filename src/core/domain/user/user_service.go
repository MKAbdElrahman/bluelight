package user

type UserService struct {
	userRepository UserRepositoty
}

func NewUserService(r UserRepositoty) *UserService {
	return &UserService{
		userRepository: r,
	}
}

type UserRegisterationParams struct {
	Name     string
	Email    string
	Password string
}

func (svc *UserService) RegisterUser(params UserRegisterationParams) (*User, error) {
	u, err := NewUser(params.Name, params.Email, params.Password)
	if err != nil {
		return nil, err
	}

	err = svc.userRepository.Create(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (svc *UserService) UpdateUser(u *User) error {
	return svc.userRepository.Update(u)
}

func (svc *UserService) GetByEmail(email string) (*User, error) {
	return svc.userRepository.GetByEmail(email)
}
