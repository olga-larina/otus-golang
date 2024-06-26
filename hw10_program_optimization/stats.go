package hw10programoptimization

import (
	"bufio"
	"io"
	"log/slog"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countUserDomains(r, domain)
}

func countUserDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	scanner := bufio.NewScanner(r)
	domainFirstLevel := "." + domain
	user := User{}

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			slog.Error("failed reading line", "error", err)
			return nil, err
		}

		matched, domain := domainByFirstLevel(user.Email, domainFirstLevel)
		if matched {
			result[domain]++
		}
	}

	return result, nil
}

func domainByFirstLevel(email string, domainFirstLevel string) (bool, string) {
	matched := strings.HasSuffix(email, domainFirstLevel)
	if !matched {
		return false, ""
	}
	splitted := strings.SplitN(email, "@", 2)
	if len(splitted) < 2 {
		return false, ""
	}
	return true, strings.ToLower(splitted[1])
}
