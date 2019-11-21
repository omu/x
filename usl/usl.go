package usl

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/purell"
)

const (
	// FallbackScheme should be commented
	FallbackScheme = "https"
)

var (
	// SupportedSchemes should be commented
	SupportedSchemes = NewSupported(
		"https",
		"http",
		"ssh",
		"git",
		"git+ssh",
		"ftp",
		"ftps",
		"file",
	)

	// SupportedProviders should be commented
	SupportedProviders = NewSupported(
		"bitbucket.com",
		"github.com",
		"gitlab.com",
		"salsa.debian.org",
	)

	// SupportedClasses should be commented
	SupportedClasses = NewSupported(
		"git",
		"tar.bz2",
		"tar.gz",
		"tar.xz",
		"tgz",
		"zip",
	)

	// classPattern
	classPattern *regexp.Regexp

	sshPattern *regexp.Regexp

	refPattern *regexp.Regexp
)

func init() {
	classPattern = regexp.MustCompile(`^(.*?)[.]` + groupPatternFromSlice("class", SupportedClasses.List) + `(/.*)?$`)
	// sshPattern was modified from https://golang.org/src/cmd/go/vcs.go.
	sshPattern = regexp.MustCompile(`^([a-zA-Z0-9_]+@)?([a-zA-Z0-9._-]+):(.*)$`)
	refPattern = regexp.MustCompile(`^(.+)@([^@]*)$`)
}

// Supported should be commented
type Supported struct {
	Support map[string]struct{}
	List    []string
}

// NewSupported should be commented
func NewSupported(items ...string) *Supported {
	s := &Supported{
		Support: map[string]struct{}{},
	}

	for _, item := range items {
		s.Support[item] = struct{}{}
		s.List = append(s.List, item)
	}

	return s
}

// Contains should be commented
func (s *Supported) Contains(support string) bool {
	_, ok := s.Support[support]

	return ok
}

// USL should be commented
type USL struct {
	in string

	Fragment string // URL
	Host     string
	Password string
	Path     string
	Port     string
	Scheme   string
	User     string
	Username string
	Class    string // USL spesific
	Domain   string
	InPath   string
	Name     string
	Ref      string
}

// Parse should be commented
func Parse(rawurl string) (*USL, error) {
	u, err := parse(rawurl)
	if err != nil {
		return nil, err
	}

	us := newFromURL(u)
	if err = us.compute(); err != nil {
		return nil, err
	}

	return us, nil
}

func (us *USL) String() string {
	var ui string

	if us.Scheme == "ssh" || us.Scheme == "git+ssh" {
		if us.Username != "" {
			ui = us.Username + "@"
		}

		base := ui + us.Host

		if us.Class == "" {
			return base + ":" + us.Name
		}

		return base + ":" + us.Name + "." + us.Class
	}

	if us.Username != "" {
		ui = us.Username

		if us.Password != "" {
			ui = ui + ":" + us.Password
		}

		ui = ui + "@"
	}

	base := us.Scheme + "://" + ui + us.Host

	if us.Class == "" {
		return base + "/" + us.Path
	}

	return base + "/" + us.Name + "." + us.Class
}

// Source should be commented
func (us *USL) Source() string {
	return us.String()
}

// ID should be commented
func (us *USL) ID() string {
	return ""
}

// Map should be commented
func (us *USL) Map() (map[string]string, []string) {
	m := map[string]string{
		"source": us.String(),
		"id":     us.ID(),
	}

	e := reflect.ValueOf(us).Elem()

	for i := 0; i < e.NumField(); i++ {
		if e.Field(i).CanInterface() {
			k := strings.ToLower(e.Type().Field(i).Name)
			v := reflect.ValueOf(e.Field(i).Interface()).String()

			m[k] = v
		}
	}

	var ks []string

	for k := range m {
		ks = append(ks, k)
	}

	sort.Strings(ks)

	return m, ks
}

// Private functions

func cut(s string, c string) (string, string) {
	i := strings.Index(s, c)

	if i < 0 {
		return s, ""
	}

	return s[:i], s[i+len(c):]
}

func (us *USL) compute() error {
	if strings.HasSuffix(us.Scheme, "+ssh") {
		us.Scheme = "ssh"
	}

	if path, ref, ok := parseRef(us.Path); ok {
		us.Path = path
		us.Ref = ref
	}

	if name, class, inpath, ok := parseClass(us.Path); ok && SupportedClasses.Contains(class) {
		us.Name = relPath(name)
		us.InPath = relPath(inpath)
		us.Class = class
	}

	if SupportedProviders.Contains(us.Host) {
		if us.Class == "" {
			us.Class = "git"
		}

		if us.Name == "" {
			parts := strings.Split(relPath(us.Path), "/")
			if len(parts) < 2 {
				return fmt.Errorf("incomplete repository path %q for provider %q: %q", us.Path, us.Host, us.in)
			}

			us.Name = strings.Join(parts[:2], "/")
			us.InPath = strings.Join(parts[2:], "/")
		}
	}

	if us.Ref != "" && us.Class == "" {
		return fmt.Errorf("malformed url: ref found for non git source: %q", us.in)
	}

	return nil
}

func groupPatternFromSlice(group string, ss []string) string {
	var escaped []string

	for _, s := range ss {
		escaped = append(escaped, regexp.QuoteMeta(s))
	}

	return `(?P<` + group + `>` + strings.Join(escaped, "|") + `)`
}

func matchFile(in string) (map[string]string, bool) {
	if strings.HasPrefix(in, "./") || strings.HasPrefix(in, "/") {
		return map[string]string{"path": in}, true
	}

	return nil, false
}

func namedMatches(re *regexp.Regexp, in string) (map[string]string, bool) {
	match := re.FindStringSubmatch(in)

	if match == nil {
		return nil, false
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result, true
}

func newFromURL(u *url.URL) *USL {
	var username, password, domain, port string

	if u.User != nil {
		username = u.User.Username()
		if p, ok := u.User.Password(); ok {
			password = p
		}
	}
	if u.Host != "" {
		var err error
		if domain, port, err = net.SplitHostPort(u.Host); err != nil {
			domain = ""
			port = ""
		}
	}

	return &USL{
		in: u.String(),

		Domain:   domain,
		Fragment: u.Fragment,
		Host:     u.Host,
		Password: password,
		Path:     u.Path,
		Port:     port,
		Scheme:   u.Scheme,
		Username: username,
	}
}

func parseClass(path string) (string, string, string, bool) {
	if m := classPattern.FindStringSubmatch(path); len(m) > 0 {
		return m[1], m[2], m[3], true
	}

	return path, "", "", false
}

func parseRef(path string) (string, string, bool) {
	if m := refPattern.FindStringSubmatch(path); len(m) >= 3 {
		return m[1], m[2], true
	}

	return path, "", false
}

func parse(rawurl string) (*url.URL, error) {
	in := rawurl

	if scheme, remaining := cut(in, "://"); remaining == "" {
		if m, ok := matchFile(in); ok {
			return parseFile(in, m)
		}

		if m, ok := matchSpecial(in); ok {
			return parseSpecial(in, m)
		}

		if m, ok := matchSSH(in); ok {
			return parseSSH(in, m)
		}

		in = FallbackScheme + "://" + in
	} else {
		scheme = strings.ToLower(scheme)

		if !SupportedSchemes.Contains(scheme) {
			return nil, fmt.Errorf("unsupported scheme %q", scheme)
		}
	}

	return parseUsual(in, nil)
}

func parseFile(_ string, match map[string]string) (*url.URL, error) {
	return &url.URL{
		Host:   "",
		Path:   filepath.Clean(match["path"]),
		Scheme: "file",
	}, nil
}

var reSSH *regexp.Regexp

func init() {
	reSSH = regexp.MustCompile(`^((?P<user>[a-zA-Z0-9_]+)@)?(?P<host>[a-zA-Z0-9._-]+):(?P<path>.*)$`)
}

func matchSSH(in string) (map[string]string, bool) {
	return namedMatches(reSSH, in)
}

func parseSSH(in string, match map[string]string) (*url.URL, error) {
	return &url.URL{
		Host:   match["host"],
		User:   url.User(match["user"]),
		Path:   filepath.Clean(match["path"]),
		Scheme: "ssh",
	}, nil
}

var reSpecial *regexp.Regexp

func init() {
	reSpecial = regexp.MustCompile(
		`^((?P<user>[a-zA-Z0-9_.-]+)@)?` +
			groupPatternFromSlice("provider", SupportedProviders.List) +
			`(?P<sep>[/:])` + `(?P<path>.*)?$`,
	)
}

func matchSpecial(in string) (map[string]string, bool) {
	return namedMatches(reSpecial, in)
}

func parseSpecial(in string, match map[string]string) (*url.URL, error) {
	user := ""
	scheme := "https"

	if match["sep"] == ":" {
		scheme = "ssh"
		user = match["user"]

		if user == "" {
			user = "git"
		} else if user != "git" {
			return nil, fmt.Errorf("user must be git where found %q", match["user"])
		}
	}

	host := match["provider"]
	path := filepath.Clean(match["path"])

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("incomplete repository path %q for provider %q: %q", path, host, in)
	}

	if repo := parts[1]; !strings.HasSuffix(repo, ".git") {
		parts[1] = repo + ".git"
		path = strings.Join(parts, "/")
	}

	return &url.URL{
		Host:   host,
		User:   url.User(user),
		Path:   path,
		Scheme: scheme,
	}, nil
}

func parseUsual(in string, _ map[string]string) (*url.URL, error) {
	normurl, err := purell.NormalizeURLString(
		in, purell.FlagsUsuallySafeGreedy|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveFragment,
	)
	if err != nil {
		return nil, err
	}

	return url.Parse(normurl)
}

func relPath(path string) string {
	return strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
}
