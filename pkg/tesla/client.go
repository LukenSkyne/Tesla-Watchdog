package tesla

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
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

func (c *Client) FirstTimeSetup() bool {
	if c.config.MainVehicle != "" {
		return false
	}

	vehicleList, err := c.GetVehicleList()

	if err != nil {
		c.log.Fatalf("unable to perform first time setup: %v", err)
	}

	if len(vehicleList.Response) < 1 {
		c.log.Fatal("no vehicles found")
	}

	c.config.MainVehicle = vehicleList.Response[0].IdString
	c.config.Save()

	c.log.Info("your main vehicle has been selected automatically, but you can change it in the config and restart")

	for _, s := range vehicleList.Response {
		c.log.Infof("%v = %v", s.DisplayName, s.Id)
	}

	return true
}

func (c *Client) checkCredentials() {
	token, _, err := new(jwt.Parser).ParseUnverified(c.config.AccessToken, jwt.MapClaims{})

	if err != nil {
		c.log.Errorf("parsing token: %v", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.log.Errorf("extracting token claims: %v", err)
		return
	}

	exp := int64(claims["exp"].(float64))
	remaining := time.Unix(exp, 0).Sub(time.Now())

	if remaining.Minutes() < 1 {
		c.log.Info("refreshing credentials")
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
		c.log.Errorf("[refreshing credentials] request failed: %v", err)
		return
	}

	defer res.Body.Close()

	var t TokenRefreshResponse
	err = json.NewDecoder(res.Body).Decode(&t)

	if err != nil {
		c.log.Errorf("[refreshing credentials] unmarshal failed: %v", err)
		return
	}

	c.config.AccessToken = t.AccessToken
	c.config.RefreshToken = t.RefreshToken
	c.config.Save()
}

func req[T any](c *Client, method string, path string, body io.Reader) (*T, error) {
	req, err := http.NewRequest(method, path, body)

	if err != nil {
		return nil, err
	}

	res, err := c.api.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var content T

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(bytes.NewReader(data)).Decode(&content)

	if err != nil {
		return nil, errors.Wrap(err, string(data))
	}

	return &content, nil
}

func Get[T any](c *Client, path string) (*T, error) {
	return req[T](c, "GET", OwnerApi+path, nil)
}

func Post[T any](c *Client, path string, body io.Reader) (*T, error) {
	return req[T](c, "POST", OwnerApi+path, body)
}
