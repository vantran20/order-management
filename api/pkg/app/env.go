package app

const (
	// EnvProd means prod env
	EnvProd = Env("prod")
	// EnvQA means qa env
	EnvQA = Env("qa")
	// EnvDev means dev env
	EnvDev = Env("dev")
)

// Env denotes the environment
type Env string

// Valid checks if the env is valid or not
func (e Env) Valid() bool {
	switch e {
	case EnvProd, EnvQA, EnvDev:
		return true
	default:
		return false
	}
}

// String returns the string representation of env
func (e Env) String() string {
	return string(e)
}
