package hw10programoptimization

import (
	json "encoding/json"
	"io"
	"strings"
)

//go:generate easyjson -all stats.go
//easyjson:json
type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
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
			if err == io.EOF {
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
