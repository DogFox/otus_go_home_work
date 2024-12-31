package hw10programoptimization

import (
	json "encoding/json"
	"errors"
	"io"
	"strings"
)

//go:generate easyjson -all stats.go
//easyjson:json
type User struct {
	ID       int    `json:"Id"`       //nolint:tagliatelle
	Name     string `json:"Name"`     //nolint:tagliatelle
	Username string `json:"Username"` //nolint:tagliatelle
	Email    string `json:"Email"`    //nolint:tagliatelle
	Phone    string `json:"Phone"`    //nolint:tagliatelle
	Password string `json:"Password"` //nolint:tagliatelle
	Address  string `json:"Address"`  //nolint:tagliatelle
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return getUsersWithDomains(r, domain)
}

func getUsersWithDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainSuffix := "." + domain
	var user User
	decoder := json.NewDecoder(r)
	for {
		if err := decoder.Decode(&user); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if strings.Contains(user.Email, domainSuffix) {
			domainPart := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[domainPart]++
		}
	}

	return result, nil
}
