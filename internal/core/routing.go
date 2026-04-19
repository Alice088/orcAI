package core

import (
	"errors"
	"strings"
)

type AgentRouter interface {
	Route(input string) (string, error)
}

type KeywordRouter struct {
	Keywords      map[string]string
	FallbackRoute *string
}

func NewKeywordRouter(keywords map[string]string, fallbackRoute *string) AgentRouter {
	return &KeywordRouter{
		Keywords:      keywords,
		FallbackRoute: fallbackRoute,
	}
}

func (kr KeywordRouter) Route(input string) (string, error) {
	strs := strings.Fields(input)

	for _, str := range strs {
		if v, ok := kr.Keywords[str]; ok {
			return v, nil
		}
	}

	if kr.FallbackRoute != nil {
		return *kr.FallbackRoute, nil
	}

	return "", errors.New("no routes")
}
