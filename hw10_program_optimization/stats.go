package hw10programoptimization

import (
	"bufio"
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
	bufReader := bufio.NewScanner(r)
	for bufReader.Scan() {
		var user User
		if err := user.UnmarshalJSON(bufReader.Bytes()); err != nil {
			return nil, err
		}
		if strings.Contains(user.Email, domainSuffix) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
