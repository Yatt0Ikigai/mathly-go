package repository

type Repositories interface {
	User() User
}

type repositories struct {
	user User
}

func NewRepositories(databases Databases) Repositories {
	return &repositories{
		user: newUser(databases.DB()),
	}
}

func (r *repositories) User() User {
	return r.user
}
