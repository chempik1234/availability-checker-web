package tokens

type TokensRepository interface {
	Check(token string) (bool, error)
	Create() (string, error)
	Delete(token string) error
}
