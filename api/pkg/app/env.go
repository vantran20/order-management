package app

const (
	// EnvProd means prod env
	EnvProd = Env("prod")
	// EnvQA means qa env
	EnvQA = Env("qa")
	// EnvDev means dev env
	EnvDev = Env("dev")
	// EnvTest means test env
	EnvTest = Env("test")
	// EnvLocal means local env
	EnvLocal = Env("local")
)

// Env denotes the environment
type Env string

// Valid checks if the env is valid or not
func (e Env) Valid() bool {
	switch e {
	case EnvProd, EnvQA, EnvDev, EnvTest, EnvLocal:
		return true
	default:
		return false
	}
}

// String returns the string representation of env
func (e Env) String() string {
	return string(e)
}
