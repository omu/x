package usl

import (
	"testing"
)

func TestSchemelessUsual(t *testing.T) {
	test(t, []testCase{
		testCase{
			"github.com/user/repo", map[string]string{
				"source": "https://github.com/user/repo.git",

				"class":  "git",
				"inpath": "",
				"name":   "user/repo",
				"ref":    "",
				"scheme": "https",
			},
		},
		testCase{
			"github.com/user/repo@unstable", map[string]string{
				"source": "https://github.com/user/repo.git",

				"class":  "git",
				"inpath": "",
				"name":   "user/repo",
				"ref":    "unstable",
				"scheme": "https",
			},
		},
		testCase{
			"github.com/user/repo/a/b@unstable", map[string]string{
				"source": "https://github.com/user/repo.git",

				"class":  "git",
				"inpath": "a/b",
				"name":   "user/repo",
				"ref":    "unstable",
				"scheme": "https",
			},
		},
		testCase{
			"github.com/user/repo.git/a/b", map[string]string{
				"source": "https://github.com/user/repo.git",

				"class":  "git",
				"inpath": "a/b",
				"name":   "user/repo",
				"ref":    "",
				"scheme": "https",
			},
		},
	})
}

func TestSchemelessSSH(t *testing.T) {
	test(t, []testCase{
		testCase{
			"github.com:user/repo", map[string]string{
				"source": "git@github.com:user/repo.git",

				"class":    "git",
				"inpath":   "",
				"name":     "user/repo",
				"ref":      "",
				"scheme":   "ssh",
				"username": "git",
			},
		},
		testCase{
			"git@github.com:user/repo", map[string]string{
				"source": "git@github.com:user/repo.git",

				"class":    "git",
				"inpath":   "",
				"name":     "user/repo",
				"ref":      "",
				"scheme":   "ssh",
				"username": "git",
			},
		},
		testCase{
			"user@example.com:a/b", map[string]string{
				"source": "user@example.com:a/b",

				"class":    "",
				"inpath":   "",
				"name":     "a/b",
				"path":     "a/b",
				"ref":      "",
				"scheme":   "ssh",
				"username": "user",
			},
		},
	})
}

func TestSchemelessFile(t *testing.T) {
	test(t, []testCase{
		testCase{
			"./a/b", map[string]string{
				"source": "a/b",

				"class":    "",
				"fullpath": "a/b",
				"scheme":   "file",
			},
		},
		testCase{
			"../a/b", map[string]string{
				"source": "../a/b",

				"class":    "",
				"fullpath": "../a/b",
				"scheme":   "file",
			},
		},
		testCase{
			"./a/b.git/x", map[string]string{
				"source": "a/b.git",

				"class":    "git",
				"fullpath": "a/b",
				"inpath":   "x",
				"scheme":   "file",
			},
		},
	})
}

func TestSchemeHTTPS(t *testing.T) {
	test(t, []testCase{
		testCase{
			"https://github.com/user/repo", map[string]string{
				"source": "https://github.com/user/repo.git",

				"class":  "git",
				"inpath": "",
				"name":   "user/repo",
				"ref":    "",
				"scheme": "https",
			},
		},
	})
}

func TestSchemeSSH(t *testing.T) {
	test(t, []testCase{
		testCase{
			"ssh://git@github.com/user/repo", map[string]string{
				"source": "git@github.com:user/repo.git",

				"class":    "git",
				"inpath":   "",
				"name":     "user/repo",
				"ref":      "",
				"scheme":   "ssh",
				"username": "git",
			},
		},
		testCase{
			"ssh://user:pass@example.com/a/b", map[string]string{
				"source": "user@example.com:a/b",

				"class":    "",
				"inpath":   "",
				"name":     "a/b",
				"password": "pass",
				"path":     "/a/b",
				"ref":      "",
				"scheme":   "ssh",
				"username": "user",
			},
		},
		testCase{
			"ssh://user:pass@example.com:22/a/b", map[string]string{
				"source": "ssh://user:pass@example.com:22/a/b",

				"class":    "",
				"domain":   "example.com",
				"host":     "example.com:22",
				"inpath":   "",
				"name":     "a/b",
				"password": "pass",
				"path":     "/a/b",
				"port":     "22",
				"ref":      "",
				"scheme":   "ssh",
				"username": "user",
			},
		},
	})
}

type testCase struct {
	in  string
	out map[string]string
}

func test(t *testing.T, cs []testCase) {
	t.Parallel()

	for _, c := range cs {
		got, err := Parse(c.in)

		if err != nil {
			t.Errorf("Parse(%q) = unexpected err %q", c.in, err)
			continue
		}

		m, _ := got.Map()

		for ke, ve := range c.out {
			if va, ok := m[ke]; ok {
				if ve != va {
					t.Errorf("\t%40s %-12s\tgot:  %-12s\twant: %-12s", c.in, ke, va, ve)
				}
			}
		}
	}
}
