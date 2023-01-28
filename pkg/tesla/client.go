package tesla

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

const (
	OwnerApi = "https://owner-api.teslamotors.com/api/1"
	AuthApi  = "https://auth.tesla.com"
)

type Client struct {
	config *Config
	api    *http.Client
	log    *zap.SugaredLogger
}

type ClientRoundTripper struct {
	client *Client
}

func (c *ClientRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	c.client.checkCredentials()

	r.Header.Set("Authorization", "Bearer "+c.client.config.AccessToken)
	r.Header.Set("Content-Type", "application/json")

	return http.DefaultTransport.RoundTrip(r)
}

func NewClient(log *zap.SugaredLogger) *Client {
	c := &Client{
		config: NewConfig("tesla.yaml", log).Load(),
		api:    &http.Client{},
		log:    log,
	}
	c.api.Transport = &ClientRoundTripper{
		client: c,
	}

	return c
}

func (c *Client) checkCredentials() {
	token, _, err := new(jwt.Parser).ParseUnverified(c.config.AccessToken, jwt.MapClaims{})

	if err != nil {
		c.log.Errorf("parsing token: %v\n", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.log.Errorf("extracting token claims: %v\n", err)
		return
	}

	exp := int64(claims["exp"].(float64))
	remaining := time.Unix(exp, 0).Sub(time.Now())

	if remaining.Minutes() < 1 {
		c.log.Infow("Refreshing Credentials", "Remaining", remaining.Hours())
		c.refreshCredentials()
	}
}

func (c *Client) refreshCredentials() {
	jsonData, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     "ownerapi",
		"refresh_token": c.config.RefreshToken,
		"scope":         "openid email offline_access",
	})
	res, err := http.Post(AuthApi+"/oauth2/v3/token", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		c.log.Errorf("request failed: %v\n", err)
		return
	}

	defer res.Body.Close()

	var t TokenRefreshResponse
	err = json.NewDecoder(res.Body).Decode(&t)

	if err != nil {
		c.log.Errorf("unmarshal failed: %v\n", err)
		return
	}

	c.config.AccessToken = t.AccessToken
	c.config.RefreshToken = t.RefreshToken
	c.config.Save()
}

func req[T any](c *Client, method string, path string, body io.Reader) *T {
	req, err := http.NewRequest(method, path, body)

	if err != nil {
		c.log.Errorf("invalid request: %v\n", err)
		return nil
	}

	res, err := c.api.Do(req)

	if err != nil {
		c.log.Errorf("request failed: %v\n", err)
		return nil
	}

	defer res.Body.Close()

	var content T

	data, err := io.ReadAll(res.Body)

	if err != nil {
		c.log.Errorf("reading body failed: %v\n", err)

		return nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	err = dec.Decode(&content)

	if err != nil {
		c.log.Errorf("unmarshal failed: %v\n", err)
		c.log.Debugf("original: %v\n", string(data))

		return nil
	}

	return &content
}

func Get[T any](c *Client, path string) *T {
	return req[T](c, "GET", OwnerApi+path, nil)
}

func Post[T any](c *Client, path string, body io.Reader) *T {
	return req[T](c, "POST", OwnerApi+path, body)
}
