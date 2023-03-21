package v1

import (
	"context"
	"crypto/rsa"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/theirish81/kumquat/auth"
	v1 "github.com/theirish81/kumquat/dto/v1"
	"github.com/theirish81/kumquat/internal"
)

type API struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	server     *echo.Echo
	users      auth.Users
}

const v1LoginBasePath = "/v1/login"
const v1SequenceBasePath = "/v1/sequence"
const contextKey = "user"
const authorizePath = "/authorize"
const runPath = "/:id/_run"
const defaultEtcPrivate = "etc/keys/private.pem"
const defaultEtcPublic = "etc/keys/public.pem"

const keyPrivateKeyPath = "PRIVATE_KEY_PATH"
const keyPublicKeyPath = "PUBLIC_KEY_PATH"
const keyTokenDuration = "TOKEN_DURATION"

func NewAPI() (*API, error) {
	priv, err := getPrivateKey()
	if err != nil {
		return nil, err
	}
	pub, err := getPublicKey()
	if err != nil {
		return nil, err
	}
	users, err := auth.LoadUsers(internal.GetUsersPath())
	if err != nil {
		return nil, err
	}
	server := echo.New()
	api := API{privateKey: priv, publicKey: pub, server: server, users: users}
	api.initRoutes()
	return &api, nil
}

func (a API) initRoutes() {
	a.server.Group(v1LoginBasePath, middleware.BasicAuth(func(username string, password string, ctx echo.Context) (bool, error) {
		ux := a.users.Authenticate(username, password)
		if ux != nil {
			ctx.Set(contextKey, ux)
			ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), contextKey, ux)))
			return true, nil
		}
		return false, nil
	})).Any(authorizePath, a.Authorize)
	a.server.Group(v1SequenceBasePath, echojwt.WithConfig(echojwt.Config{
		SigningMethod: jwt.SigningMethodRS512.Alg(),
		SigningKey:    a.publicKey,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(v1.Claims)
		},
	})).Any(runPath, a.SequenceRun)
}

func (a API) Run() {
	_ = http.ListenAndServe(":5000", a.server)
}

func (a API) Authorize(c echo.Context) error {
	user := c.Get(contextKey).(*auth.User)
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, newClaimsFromUser(user))
	if str, err := token.SignedString(a.privateKey); err == nil {
		return c.JSON(200, v1.AccessToken{AccessToken: str})
	} else {
		return c.JSON(http.StatusInternalServerError, v1.Error{Code: 500, Message: "unable to generate token"})
	}

}

func (a API) SequenceRun(c echo.Context) error {
	claims := c.Get(contextKey).(*jwt.Token).Claims.(*v1.Claims)
	sequenceId := c.Param("id")
	if !claims.CanAccess(sequenceId) {
		return c.JSON(403, v1.Error{Code: 403, Message: "cannot access sequence"})
	}

	seqPath := internal.GetSequencePath(sequenceId)
	sequence, err := internal.LoadSequence(seqPath)
	if err != nil {
		return c.JSON(500, v1.Error{Code: 500, Message: "sequence error"})
	}
	// If AcceptParams is set, it means we'll take the JSON body into the scope
	if sequence.AcceptParams {
		// parsing the body into a data structure
		params, err := extractParams(c.Request())
		if err != nil {
			return c.JSON(400, v1.Error{Code: 400, Message: "bad request"})
		}
		// If the parsed data structure does not contain all the fields described in the "Requires" section,
		// then we return an error
		if err := sequence.CheckRequires(params); err != nil {
			return c.JSON(400, v1.Error{Code: 400, Message: err.Error()})
		}
		// otherwise, we set it into the scope
		sequence.Scope.InsertParams(params)
	}
	// Running the sequence
	sequence.Run(c.Request().Context())
	result := sequence.Result()
	return c.JSON(200, result)
}

func extractParams(r *http.Request) (map[string]any, error) {
	if data, err := io.ReadAll(r.Body); err == nil {
		if len(data) > 0 {
			params, err := simplejson.NewJson(data)
			if err != nil {
				return nil, err
			}
			return params.Map()
		}
		return nil, nil
	} else {
		return nil, err
	}
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	keysPath := os.Getenv(keyPrivateKeyPath)
	if len(keysPath) == 0 {
		keysPath = defaultEtcPrivate
	}
	data, err := os.ReadFile(keysPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(data)
}

func getPublicKey() (*rsa.PublicKey, error) {
	keysPath := os.Getenv(keyPublicKeyPath)
	if len(keysPath) == 0 {
		keysPath = defaultEtcPublic
	}
	data, err := os.ReadFile(keysPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(data)
}

func newClaimsFromUser(user *auth.User) v1.Claims {
	tokenDurationString := os.Getenv(keyTokenDuration)
	tokenDuration := 24 * time.Hour
	if len(tokenDurationString) > 0 {
		duration, err := time.ParseDuration(tokenDurationString)
		if err != nil {
			log.Err(err).Msg("could not parse TOKEN_DURATION. Using default")
		}
		tokenDuration = duration
	}
	return v1.Claims{Access: user.Access, Sequences: user.Sequences, Exp: time.Now().Add(tokenDuration).UnixMilli()}
}
