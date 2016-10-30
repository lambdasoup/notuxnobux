package steam

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const (
	steamLogin = "https://steamcommunity.com/openid/login"

	openIDNs         = "http://specs.openid.net/auth/2.0"
	openIDIdentifier = "http://specs.openid.net/auth/2.0/identifier_select"
)

var (
	validationRegexp    = regexp.MustCompile("^(http|https)://steamcommunity.com/openid/id/[0-9]{15,25}$")
	extractDigitsRegexp = regexp.MustCompile("\\D+")
)

// AuthURLFor returns the openID auth url for the given request context
func AuthURLFor(r *http.Request) *url.URL {
	c := appengine.NewContext(r)

	// ignore err bc is const
	u, _ := url.Parse(steamLogin)

	q := u.Query()
	q.Set("openid.claimed_id", openIDIdentifier)
	q.Set("openid.identity", openIDIdentifier)
	q.Set("openid.mode", "checkid_setup")
	q.Set("openid.ns", openIDNs)
	q.Set("openid.realm", realm(r).String())
	q.Set("openid.return_to", requestURI(r).String())
	u.RawQuery = q.Encode()

	log.Debugf(c, "redirect url: %v", u.String())

	return u
}

// FetchID gets the steam id request's user
func FetchID(r *http.Request) (string, error) {
	c := appengine.NewContext(r)

	// validate mode
	if r.URL.Query().Get("openid.mode") != "id_res" {
		return "", fmt.Errorf("unexpected mode. was %v", r.Form.Get("openid.mode"))
	}

	// validate return_to
	uri1 := requestURI(r)
	uri2, err := url.ParseRequestURI(r.URL.Query().Get("openid.return_to"))
	if err != nil {
		return "", errors.New("could not parse return_to URI")
	}
	if uri1.Scheme != uri2.Scheme || uri1.Host != uri2.Host || uri1.Path != uri2.Path {
		return "", fmt.Errorf("The \"return_to url\" must match the url of current request. (%v) != (%v)", uri1, uri2)
	}

	params := make(url.Values)
	params.Set("openid.assoc_handle", r.URL.Query().Get("openid.assoc_handle"))
	params.Set("openid.signed", r.URL.Query().Get("openid.signed"))
	params.Set("openid.sig", r.URL.Query().Get("openid.sig"))
	params.Set("openid.ns", r.URL.Query().Get("openid.ns"))

	split := strings.Split(r.URL.Query().Get("openid.signed"), ",")
	for _, item := range split {
		params.Set("openid."+item, r.URL.Query().Get("openid."+item))
	}
	params.Set("openid.mode", "check_authentication")

	client := urlfetch.Client(c)
	resp, err := client.PostForm(steamLogin, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	response := strings.Split(string(content), "\n")
	if response[0] != "ns:"+openIDNs {
		return "", errors.New("Wrong ns in the response.")
	}
	if strings.HasSuffix(response[1], "false") {
		return "", errors.New("Unable validate openId.")
	}

	openIDURL := r.URL.Query().Get("openid.claimed_id")
	if !validationRegexp.MatchString(openIDURL) {
		return "", errors.New("Invalid steam id pattern.")
	}

	return extractDigitsRegexp.ReplaceAllString(openIDURL, ""), nil
}

// we need to somehow gather URL data from different places to make it work in
// dev _and_ prod.
// http://stackoverflow.com/a/6911635/470509
func requestURI(r *http.Request) *url.URL {
	u := realm(r)

	u.Path = r.URL.Path

	return u
}

func realm(r *http.Request) *url.URL {
	u := new(url.URL)

	if r.TLS != nil {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	u.Host = r.Host

	return u
}

// Mode returns the current openID mode
func Mode(r *http.Request) string {
	return r.URL.Query().Get("openid.mode")
}
