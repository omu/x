package usl

import "testing"

type testCase struct {
	in  string
	out map[string]string
}

//nolint:funlen
func TestParse(t *testing.T) {
	t.Parallel()

	tests := map[string][]testCase{
		"schemeless usual": {
			{
				"github.com/user/repo", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"inpath": "",
					"name":   "user/repo",
					"ref":    "",
					"scheme": "https",
				},
			},
			{
				"github.com/user/repo@unstable", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"inpath": "",
					"name":   "user/repo",
					"ref":    "unstable",
					"scheme": "https",
				},
			},
			{
				"github.com/user/repo/a/b@unstable", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"inpath": "a/b",
					"name":   "user/repo",
					"ref":    "unstable",
					"scheme": "https",
				},
			},
			{
				"github.com/user/repo.git/a/b", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"inpath": "a/b",
					"name":   "user/repo",
					"ref":    "",
					"scheme": "https",
				},
			},
		},
		"schemeless SSH": {
			{
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
			{
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
			{
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
		},
		"schemeless File": {
			{
				"./a/b", map[string]string{
					"source": "a/b",

					"class":    "",
					"basepath": "a/b",
					"scheme":   "file",
				},
			},
			{
				"../a/b", map[string]string{
					"source": "../a/b",

					"class":    "",
					"basepath": "../a/b",
					"scheme":   "file",
				},
			},
			{
				"./a/b.git/x", map[string]string{
					"source": "a/b.git",

					"class":    "git",
					"basepath": "a/b",
					"inpath":   "x",
					"scheme":   "file",
				},
			},
		},
		"HTTPS scheme": {
			{
				"https://github.com/user/repo", map[string]string{
					"source": "https://github.com/user/repo.git",

					"class":  "git",
					"inpath": "",
					"name":   "user/repo",
					"ref":    "",
					"scheme": "https",
				},
			},
		},
		"SSH scheme": {
			{
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
			{
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
			{
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
		},
	}

	for name, ts := range tests {
		ts := ts // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, tc := range ts {
				got, err := Parse(tc.in)

				if err != nil {
					t.Errorf("Parse(%q) = unexpected err %q", tc.in, err)
					continue
				}

				m, _ := got.Map()

				for ke, ve := range tc.out {
					if va, ok := m[ke]; ok {
						if ve != va {
							t.Errorf("\t%40s    %-12s\twant: %-12s\tgot:  %-12s", tc.in, ke, ve, va)
						}
					}
				}
			}
		})
	}
}
