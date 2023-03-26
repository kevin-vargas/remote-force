package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"remote-force/server/entity"
	"remote-force/server/interceptor"
	"remote-force/server/store"
	"strings"
)

type provider struct {
	client     *http.Client
	tokenStore store.Token
}

const (
	user_email_url = "https://api.github.com/user/emails"
)

type userEmailResponse struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func (p *provider) GetByID(id string) (usr entity.User, err error) {
	o, ok, err := p.tokenStore.Get(id)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("not found user token")
		return
	}
	req, err := http.NewRequest(http.MethodGet, user_email_url, nil)
	if err != nil {
		return
	}
	// TODO: replace for capita letter
	tokenType := strings.Title(o.TokenType)
	token := tokenType + " " + o.AccessToken
	req.Header.Set("Authorization", token)
	var emails []userEmailResponse
	res, err := p.client.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("invalid status code %d", res.StatusCode)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&emails); err != nil {
		return
	}
	for _, email := range emails {
		if email.Primary {
			return entity.User{
				Email: email.Email,
			}, nil
		}
	}
	err = errors.New("not primary email on github")
	return
}

func New(ts store.Token) interceptor.UserProvider {
	return &provider{
		tokenStore: ts,
	}
}
