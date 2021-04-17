package auth

type Captcha interface {
	Generate() (string, string, error)
	Verify(id, answer string) bool
}

type Token interface {
	Create(interface{}) (string, error)
	Verify(string) (interface{}, error)
}

