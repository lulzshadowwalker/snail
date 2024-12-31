package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestReadLinks(t *testing.T) {
	cases := []struct {
		Name      string
		InputURL  string
		InputBody string
		Expected  []string
	}{
		{
			Name:     "absolute and relative URLs",
			InputURL: "https://blog.boot.dev",
			InputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			Expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			Name:      "no links",
			InputURL:  "https://blog.boot.dev",
			InputBody: "<html></html>",
			Expected:  []string{},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
      u, err := url.Parse(c.InputURL)
      if err != nil {
        t.Errorf("failed to parse url %s", err)
      }

			got, err := ReadLinks(c.InputBody, u)
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}

			if len(got) != len(c.Expected) {
				t.Errorf("got %d urls, expected %d", len(got), len(c.Expected))
			}

			if !reflect.DeepEqual(got, c.Expected) {
				t.Errorf("got %#v want %#v", got, c.Expected)
			}
		})
	}
}
