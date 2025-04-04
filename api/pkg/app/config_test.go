package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_IsValid(t *testing.T) {
	type arg struct {
		givenCfg Config
		expErr   error
	}
	tcs := map[string]arg{
		"all good": {
			givenCfg: Config{
				ProjectName:      "project",
				AppName:          "app",
				SubComponentName: "sub",
				Env:              "dev",
				Version:          "1.0.0",
				Server:           "server",
				ProjectTeam:      "team",
			},
		},
		"no team": {
			givenCfg: Config{
				ProjectName:      "project",
				AppName:          "app",
				SubComponentName: "sub",
				Env:              "dev",
				Server:           "server",
			},
		},
		"no version": {
			givenCfg: Config{
				ProjectName:      "project",
				AppName:          "app",
				SubComponentName: "sub",
				Env:              "dev",
				Server:           "server",
			},
		},
		"no server": {
			givenCfg: Config{
				ProjectName:      "project",
				AppName:          "app",
				SubComponentName: "sub",
				Env:              "dev",
				Version:          "1.0.0",
			},
		},
		"no env": {
			givenCfg: Config{
				ProjectName:      "project",
				AppName:          "app",
				SubComponentName: "sub",
			},
			expErr: ErrInvalidAppConfig,
		},
		"invalid env": {
			givenCfg: Config{
				ProjectName:      "project",
				AppName:          "app",
				SubComponentName: "sub",
				Env:              "env",
			},
			expErr: ErrInvalidAppConfig,
		},
		"no sub component": {
			givenCfg: Config{
				ProjectName: "project",
				AppName:     "app",
				Env:         "dev",
			},
			expErr: ErrInvalidAppConfig,
		},
		"no app": {
			givenCfg: Config{
				ProjectName:      "project",
				SubComponentName: "sub",
				Env:              "dev",
			},
			expErr: ErrInvalidAppConfig,
		},

		"no project": {
			givenCfg: Config{
				AppName:          "app",
				SubComponentName: "sub",
				Env:              "dev",
			},
			expErr: ErrInvalidAppConfig,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			require.True(t, errors.Is(tc.givenCfg.IsValid(), tc.expErr))
		})
	}
}
