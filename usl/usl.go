package usl

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/purell"
)

// FallbackScheme should be commented
const FallbackScheme = "https"

var (
	// classPattern
	classPattern = regexp.MustCompile(`^(.+)[.]([a-zA-Z0-9_.-]+)(/.*)$`)
	// schemePattern
	schemePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+-.]*://`)
	// scpPattern was modified from https://golang.org/src/cmd/go/vcs.go.
	scpPattern = regexp.MustCompile(`^([a-zA-Z0-9_]+@)?([a-zA-Z0-9._-]+):(.*)$`)
	// refPattern
	refPattern = regexp.MustCompile(`^(.+)@([^@]*)$`)
)

var (
	// SupportedProtocols should be commented
	SupportedProtocols = NewSupported(
		"ssh",
		"git",
		"git+ssh",
		"http",
		"https",
		"ftp",
		"ftps",
		"rsync",
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
		"zip",
		"tgz",
	)
)

// Supported should be commented
type Supported struct {
	Support map[string]struct{}
}

// NewSupported should be commented
func NewSupported(items ...string) *Supported {
	s := &Supported{
		Support: map[string]struct{}{},
	}

	for _, i := range items {
		s.Support[i] = struct{}{}
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
	Fragment 	string // URL
	Host     	string
	Password 	string
	Path     	string
	Port     	string
	Scheme   	string
	Username 	string
	Class    	string // USL spesific
	Domain   	string
	ID       	string
	Name     	string
	Ref      	string
	Source   	string
	Target   	string
}

// NewUSLFromURL should be commented
func NewUSLFromURL(u *url.URL) *USL {
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

// Parse should be commented
func Parse(rawurl string) (*USL, error) {
	us, err := parse(rawurl)
	if err == nil {
		us.expand()
	}

	return us, err
}

// Dump should be commented
func (us *USL) Dump(attributes ...string) {
	m, ks := us.ToMap()

	var wanted []string

	if len(attributes) > 0 {
		for i, attribute := range attributes {
			attribute = strings.ToLower(attribute)
			attributes[i] = attribute
			wanted = append(wanted, strings.Title(attribute))
		}
	} else {
		wanted = ks
	}

	for _, attribute := range wanted {
		if value, ok := m[attribute]; ok && value != "" {
			fmt.Printf("%-16s %s\n", attribute, value)
		}
	}
}

// Print should be commented
func (us *USL) Print(attributes ...string) {
	m, ks := us.ToMap()

	var wanted, values []string

	if len(attributes) > 0 {
		for i, attribute := range attributes {
			attribute = strings.ToLower(attribute)
			attributes[i] = attribute
			wanted = append(wanted, strings.Title(attribute))
		}
	} else {
		wanted = ks
	}

	for _, attribute := range wanted {
		if value, ok := m[attribute]; ok {
			values = append(values, value)
		}
	}

	fmt.Println(strings.Join(values[:], " "))
}

// ToMap should be commented
func (us *USL) ToMap() (map[string]string, []string) {
	m := make(map[string]string)
	var ks []string

	e := reflect.ValueOf(us).Elem()

	for i := 0; i < e.NumField(); i++ {
		if e.Field(i).CanInterface() {
			k := e.Type().Field(i).Name
			v := reflect.ValueOf(e.Field(i).Interface()).String()

			m[k] = v
			ks = append(ks, k)
		}
	}

	sort.Strings(ks)

	return m, ks
}

func (us *USL) expand() {
	if path, ref, ok := parseRef(us.Path); ok {
		us.Path = path
		us.Ref = ref
	}

	if name, class, target, ok := parseClass(us.Path); ok && SupportedClasses.Contains(class) {
		us.Name = cleanPath(name)
		us.Target = cleanPath(target)
		us.Class = class
	}

	if us.Class == "" && SupportedProviders.Contains(us.Host) {
		us.Class = "git"
	}

	if us.Class == "" {
		us.Source = us.Scheme + "://" + us.Host + "/" + us.Path
	} else {
		us.Source = us.Scheme + "://" + us.Host + "/" + us.Name
	}
}

func newFileUSL(rawurl string) *USL {
	return NewUSLFromURL(&url.URL{
		Scheme: "file",
		Host:   "",
		Path:   rawurl,
	})
}

func parse(rawurl string) (*USL, error) {
	normurl, err := purell.NormalizeURLString(
		rawurl, purell.FlagsUsuallySafeGreedy|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveFragment,
	)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(rawurl, "./") || strings.HasPrefix(rawurl, "/") {
		return newFileUSL(normurl), nil
	}

	if !HasScheme(normurl) {
		// TODO Supported provider?
		if us, ok := trySCP(normurl); ok {
			return us, nil
		}

		normurl = FallbackScheme + "://" + normurl
	}

	u, err := url.Parse(normurl)
	if err != nil {
		return nil, err
	}

	return NewUSLFromURL(u), nil
}

func parseRef(path string) (string, string, bool) {
	if m := refPattern.FindStringSubmatch(path); len(m) >= 3 {
		return m[1], m[2], true
	}

	return path, "", false
}

func parseClass(path string) (string, string, string, bool) {
	if m := classPattern.FindStringSubmatch(path); len(m) >= 4 {
		return m[1], m[2], m[3], true
	}

	return path, "", "", false
}

func cleanPath(path string) string {
	return strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
}

// HasScheme should be commented
func HasScheme(rawurl string) bool {
	return schemePattern.MatchString(rawurl)
}

func trySCP(rawurl string) (*USL, bool) {
	match := scpPattern.FindAllStringSubmatch(rawurl, -1)
	if len(match) == 0 {
		return nil, false
	}
	m := match[0]

	user := strings.TrimRight(m[1], "@")
	var userinfo *url.Userinfo
	if user != "" {
		userinfo = url.User(user)
	}

	return NewUSLFromURL(&url.URL{
		Scheme: "ssh",
		User:   userinfo,
		Host:   m[2],
		Path:   m[3],
	}), true
}
