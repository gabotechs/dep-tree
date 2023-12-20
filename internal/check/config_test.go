package check

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_Check(t *testing.T) {
	a := require.New(t)

	tests := []struct {
		Name   string
		Config Config
		From   string
		To     string
		Passes bool
	}{
		{
			Name: "white list passes",
			Config: Config{
				WhiteList: map[string][]string{
					"white": {"**pass**"},
				},
			},
			From:   "white",
			To:     "this is going to pass",
			Passes: true,
		},
		{
			Name: "white list fails",
			Config: Config{
				WhiteList: map[string][]string{
					"white": {"**pass**"},
				},
			},
			From:   "white",
			To:     "this doesn't",
			Passes: false,
		},
		{
			Name: "black list passes",
			Config: Config{
				BlackList: map[string][]string{
					"black": {"**fail**"},
				},
			},
			From:   "black",
			To:     "this is going to pass",
			Passes: true,
		},
		{
			Name: "black list fails",
			Config: Config{
				BlackList: map[string][]string{
					"black": {"**fail**"},
				},
			},
			From:   "black",
			To:     "this is going to fail",
			Passes: false,
		},
		{
			Name: "this should never pass",
			Config: Config{
				BlackList: map[string][]string{
					"**": {"**"},
				},
			},
			From:   "black",
			To:     "this is going to fail",
			Passes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			pass, err := tt.Config.Check(tt.From, tt.To)
			a.NoError(err)
			a.Equal(tt.Passes, pass)
		})
	}
}
