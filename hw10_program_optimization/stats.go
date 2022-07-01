package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var b strings.Builder
	b.WriteString(".")
	b.WriteString(domain)
	dDomain := b.String()

	reader := bufio.NewReader(r)

	result := make(DomainStat)

	for {
		user := User{}
		line, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if err = easyjson.Unmarshal(line, &user); err != nil {
			return nil, err
		}
		if strings.Contains(user.Email, dDomain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
