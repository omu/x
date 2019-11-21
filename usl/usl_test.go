package usl

import "testing"

var tests = []struct {
	in       string
	expected map[string]string
	out      string
}{
	{
		"https://github.com/user/repo", map[string]string{
			"class":  "git",
			"inpath": "",
			"name":   "user/repo",
			"ref":    "",
			"scheme": "https",
		}, "https://github.com/user/repo.git",
	},
	{
		"github.com/user/repo", map[string]string{
			"class":  "git",
			"inpath": "",
			"name":   "user/repo",
			"ref":    "",
			"scheme": "https",
		}, "https://github.com/user/repo.git",
	},
	{
		"github.com:user/repo", map[string]string{
			"class":    "git",
			"inpath":   "",
			"name":     "user/repo",
			"ref":      "",
			"scheme":   "ssh",
			"username": "git",
		}, "git@github.com:user/repo.git",
	},
	{
		"git@github.com:user/repo", map[string]string{
			"class":    "git",
			"inpath":   "",
			"name":     "user/repo",
			"ref":      "",
			"scheme":   "ssh",
			"username": "git",
		}, "git@github.com:user/repo.git",
	},
}

func TestParse(t *testing.T) {
	for _, sample := range tests {
		got, err := Parse(sample.in)
		if err != nil {
			t.Errorf("Parse(%q) = unexpected err %q, want %q", sample.in, err, sample.expected)
			continue
		}

		m, _ := got.Map()

		for ke, ve := range sample.expected {
			if va, ok := m[ke]; ok {
				if ve != va {
					t.Errorf("Got (%s), expected (%s) for attribute (%s)", va, ve, ke)
					// t.Errorf("Parse(%q) = %q, want %q", sample.in, got, sample.expected)
				}
			} else {
				t.Errorf("Err\n")
			}
		}

		if actual := got.String(); sample.out != actual {
			t.Errorf("Got (%s), expected (%s)", actual, sample.out)
		}
	}
}

// TODO
// "foo://example.com:8042/over/there?name=ferret#nose",
// "urn:example:animal:ferret:nose",
// "jdbc:mysql://test_user:ouupppssss@localhost:3306/sakila?profileSQL=true",
// "ftp://ftp.is.co.za/rfc/rfc1808.txt",
// "http://www.ietf.org/rfc/rfc2396.txt#header1",
// "ldap://[2001:db8::7]/c=GB?objectClass=one&objectClass=two",
// "mailto:John.Doe@example.com",
// "news:comp.infosystems.www.servers.unix",
// "tel:+1-816-555-1212",
// "telnet://192.0.2.16:80/",
// "urn:oasis:names:specification:docbook:dtd:xml:4.1.2",
//
// "ssh://alice@example.com",
// "https://bob:pass@example.com/place",
// "http://example.com/?a=1&b=2+2&c=3&c=4&d=%65%6e%63%6F%64%65%64",
// "https://github.com/roktas///t.git/a/b@next#issue/12",
// "file://foo/bar",
// "./x/y",
// "example.com:a/b",
// "./fexample.com:a/b",
// "any://server.com:a/b",
