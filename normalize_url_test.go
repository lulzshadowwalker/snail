package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	cases := []struct {
		Name     string
		Expected string
		Input    string
	}{
		{
			Name:     "Remove scheme",
			Input:    "https://example.com",
			Expected: "example.com",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := NormalizeURL(c.Input)
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}

			want := c.Expected

			if got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}
}
