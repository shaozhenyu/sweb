package odm

type AuthError struct {
}

func (AuthError) Error() string {
	return "no auth"
}

var (
	ErrNoAuth = &AuthError{}
)

type IUser interface {
	Uid() int64
}
