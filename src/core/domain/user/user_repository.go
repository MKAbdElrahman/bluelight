package user

type UserRepositoty interface {
	Create(u *User) error
	GetByEmail(email string) (*User, error)
	Update(u *User) error
}
