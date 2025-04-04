package app

import "github.com/friendsofgo/errors"

// Config holds basic information about the app
type Config struct {
	ProjectName      string // mandatory
	AppName          string // mandatory
	SubComponentName string // mandatory
	Env              Env    // mandatory
	Version          string
	Server           string
	ProjectTeam      string
}

// IsValid checks if the config is valid or not
func (c Config) IsValid() error {
	switch {
	case
		c.ProjectName == "",
		c.AppName == "",
		c.SubComponentName == "",
		!c.Env.Valid():
		return errors.WithStack(ErrInvalidAppConfig)
	default:
		return nil
	}
}
