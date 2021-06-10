package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"keboola-as-code/src/model"
)

func (a StorageApi) WithToken(token *model.Token) *StorageApi {
	a.token = token
	a.http.resty.SetHeader("X-StorageApi-Token", token.Token)
	return &a
}

func (a *StorageApi) Token() *model.Token {
	if a.token == nil {
		panic(fmt.Errorf("token is not set"))
	}

	return a.token
}

func (a *StorageApi) ProjectId() int {
	if a.token == nil {
		panic(fmt.Errorf("token is not set"))
	}
	return a.token.ProjectId()
}

func (a *StorageApi) ProjectName(token string) string {
	if a.token == nil {
		panic(fmt.Errorf("token is not set"))
	}
	return a.token.ProjectName()
}

func (a *StorageApi) GetToken(token string) (*model.Token, error) {
	if res, err := a.GetTokenR(token).Send(); err != nil {
		return nil, err
	} else {
		return res.Result().(*model.Token), nil
	}
}

func (a *StorageApi) GetTokenR(token string) *resty.Request {
	return a.http.
		R(resty.MethodGet, "/tokens/verify").
		SetHeader("X-StorageApi-Token", token).
		SetResult(&model.Token{})
}
