package signed

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Signed{})
}

// Middleware implements an HTTP handler that writes the
// visitor's IP address to a file or stream.
type Signed struct {
	Secret string `json:"secret,omitempty"`
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Signed) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.matchers.signed",
		New: func() caddy.Module { return new(Signed) },
	}
}

func (s *Signed) Provision(ctx caddy.Context) error {
	s.logger = ctx.Logger()
	return nil
}

func (s *Signed) Validate() error {
	if s.Secret == "" {
		return fmt.Errorf("secret is required")
	}
	return nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (s *Signed) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		args := d.RemainingArgs()

		switch len(args) {
		case 1:
			s.Secret = args[0]
		default:
			return d.Err("unexpected number of arguments")
		}
	}

	return nil
}

func (s *Signed) Match(r *http.Request) bool {
	queryParameters := r.URL.Query()
	if queryParameters == nil {
		// no query params set, deny request
		return false
	}
	token := strings.ToLower(queryParameters.Get("token"))
	if token == "" {
		// token not set, deny request
		return false
	}
	modifiedURL := r.URL
	if r.TLS == nil {
		modifiedURL.Scheme = "http"
	} else {
		modifiedURL.Scheme = "https"
	}
	modifiedURL.Host = r.Host
	modifiedQuery := modifiedURL.Query()
	modifiedQuery.Del("token")
	modifiedURL.RawQuery = modifiedQuery.Encode()
	correctToken := generateToken(modifiedURL.String())
	if token == correctToken {
		// Token is correct!
		return true
	}
	return false
}

func generateToken(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	return strings.ToLower(sha)
}

// Interface guards
var (
	_ caddy.Provisioner        = (*Signed)(nil)
	_ caddy.Module             = (*Signed)(nil)
	_ caddyhttp.RequestMatcher = (*Signed)(nil)
	_ caddyfile.Unmarshaler    = (*Signed)(nil)
)
