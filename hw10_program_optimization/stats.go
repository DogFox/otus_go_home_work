package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
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

	regex, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	bufReader := bufio.NewReader(r)
	// Читаем построчно
	for {
		var user User
		line, err := bufReader.ReadString('\n')
		if err != nil {
			// Проверяем на конец файла
			if err == io.EOF {
				if err = user.UnmarshalJSON([]byte(line)); err != nil {
					return nil, err
				}
				if regex.MatchString(user.Email) {
					result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] += 1
				}
				break
			}
			return nil, err
		}
		if err = user.UnmarshalJSON([]byte(line)); err != nil {
			return nil, err
		}
		if regex.MatchString(user.Email) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] += 1
		}

	}

	return result, nil
}
