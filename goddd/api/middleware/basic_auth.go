package middleware

import (
	"crypto/sha256"
	"encoding/base64"
	"goddd/api/httperror"
	"goddd/domain/services"
	"goddd/pkg/errors"

	"github.com/gin-gonic/gin"
)

const AuthorizationBasic = "Basic"

// Middleware for basic authentification
// Header set to `Basic username:password`.
func BasicAuth(serv *services.UserService) gin.HandlerFunc {
	// lambda for credential check
	checkCredential := func(c *gin.Context, username, pwd string) error {
		if username == "" || pwd == "" {
			// missing credentials
			return errors.MissingAuthCredentials.New("missing authentification credentials (username and password required)")
		}
		// retrieve user
		u, err := serv.GetUserByUsername(username)
		if err != nil {
			// user not found in repository
			if errors.GetID(err) == errors.RepoItemNotFound {
				return errors.UnAuthorized.Wrap(err, "invalid credentials")
			}
			// other error (internal)
			return errors.InternalError.Wrap(err, "authentication failure")
		}

		h := sha256.New()
		h.Write([]byte(pwd))

		// hashed password should match the stored one
		if base64.StdEncoding.EncodeToString(h.Sum(nil)) != u.Password {
			// password does not match
			return errors.UnAuthorized.Wrap(err, "invalid credentials")
		}
		// set user
		c.Set("user-id", u.Id)
		return nil
	}
	// handler
	return func(c *gin.Context) {
		username, pwd, ok := c.Request.BasicAuth()
		if !ok {
			// invalid basic authen
			httperror.Abort(c, errors.InvalidAuthProtocol.New("missing or invalid authorization header"))
			return
		}
		if err := checkCredential(c, username, pwd); err != nil {
			// authen failure
			httperror.Abort(c, err)
			return
		}
		c.Next()
	}
}
