package types

type Db interface {
	Save(user *User) error
	GetById(id string) (*User, error)
	Get(username, password string) (*User, error)
}

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}
