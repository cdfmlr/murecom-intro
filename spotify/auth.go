package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// AuthResp is the response of Auth
type AuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Auth authenticate using a client credentials flow
func Auth(clientId, clientSecret string) (*AuthResp, error) {
	authUrl := "https://accounts.spotify.com/api/token"

	// post form
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	form := []byte(data.Encode())

	//form, err := json.Marshal(map[string]string{"grant_type": "client_credentials"})
	//if err != nil {
	//	return nil, AuthError{
	//		Field:   "json marshal req form",
	//		Err:     err,
	//		Context: form,
	//	}
	//}
	//fmt.Printf("form: %#v\n", string(form))

	// construct request
	req, err := http.NewRequest(http.MethodPost,
		authUrl,
		bytes.NewReader(form))

	if err != nil {
		return nil, AuthError{
			Field:   "new request",
			Err:     err,
			Context: req,
		}
	}

	// header
	req.Header.Set("Authorization",
		"Basic "+base64.StdEncoding.EncodeToString(
			[]byte(clientId+":"+clientSecret)))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	//fmt.Printf("req: %#v\n", req)

	// do request
	c := http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return nil, AuthError{
			Field:   "do request",
			Err:     err,
			Context: resp,
		}
	}
	if resp.StatusCode != 200 {
		return nil, AuthError{
			Field:   fmt.Sprint("resp status ", resp.StatusCode),
			Err:     err,
			Context: resp.Status,
		}
	}

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, AuthError{
			Field:   "read resp body",
			Err:     err,
			Context: resp,
		}
	}

	// parse json
	var respJson AuthResp
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return &respJson, AuthError{
			Field:   "json unmarshal",
			Err:     err,
			Context: string(body),
		}
	}

	return &respJson, nil
}

// AuthError Spotify 认证失败
type AuthError struct {
	Field   string      // 失败的地方
	Err     error       // 遇到的错误
	Context interface{} // 附加的信息
}

func (a AuthError) Error() string {
	return fmt.Sprintf("auth_failed(%v): %v: %#v", a.Field, a.Err, a.Context)
}

type authHolder struct {
	ClientID     string
	ClientSecret string
	token        oauth2.Token
	mutex        sync.RWMutex
}

func (a *authHolder) GetToken() *oauth2.Token {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	expiry := a.token.Expiry.Add(-60 * time.Second)
	if a.token.AccessToken != "" && time.Now().Before(expiry) {
		return &a.token
	}
	// expired
	a.mutex.RUnlock()
	a.refreshToken()
	a.mutex.RLock()

	return &a.token
}

func (a *authHolder) refreshToken() {
	var err error
	for tries := 30; tries > 0; tries-- {
		authResult, err := Auth(a.ClientID, a.ClientSecret)
		if err == nil { // success
			a.mutex.Lock()
			a.token = oauth2.Token{
				AccessToken:  authResult.AccessToken,
				TokenType:    authResult.TokenType,
				RefreshToken: "",
				Expiry:       time.Now().Add(time.Duration(authResult.ExpiresIn) * time.Second),
			}
			a.mutex.Unlock()
			return
		} else { // failed
			fmt.Println("authHolder.refreshToken: auth failed (to retry): ", err)
			time.Sleep(1 * time.Second)
		}
	}
	// TODO: All retries failed
	fmt.Println("authHolder.refreshToken: All retries failed: ", err)
}

// AuthHolder global spotify token holder 单例
var AuthHolder = &authHolder{}

func NewSpotifyClient() *spotify.Client {
	c := spotify.Authenticator{}.NewClient(AuthHolder.GetToken())
	c.AutoRetry = true
	return &c
}

func InitAuth() {
	AuthHolder.ClientID = Config.ClientID
	AuthHolder.ClientSecret = Config.ClientSecret
	AuthHolder.refreshToken()
}
