package tracker

import (
	"errors"
	"regexp"
)

var ErrEmptyClientRule = errors.New("client rule requires a peer ID or agent pattern")

type ClientRule struct {
	PeerRegex  *regexp.Regexp
	AgentRegex *regexp.Regexp
}

func CompileClientRule(peerPattern, agentPattern string) (ClientRule, error) {
	if peerPattern == "" && agentPattern == "" {
		return ClientRule{}, ErrEmptyClientRule
	}
	rule := ClientRule{}
	if peerPattern != "" {
		re, err := regexp.Compile(peerPattern)
		if err != nil {
			return ClientRule{}, err
		}
		rule.PeerRegex = re
	}
	if agentPattern != "" {
		re, err := regexp.Compile(agentPattern)
		if err != nil {
			return ClientRule{}, err
		}
		rule.AgentRegex = re
	}
	return rule, nil
}

func (r ClientRule) Match(peerId, userAgent string) bool {
	return (r.PeerRegex == nil || r.PeerRegex.MatchString(peerId)) && (r.AgentRegex == nil || r.AgentRegex.MatchString(userAgent))
}

type ClientException struct {
	Name          string `json:"name"`
	PeerIdPattern string `json:"peerIdPattern"`
	AgentPattern  string `json:"agentPattern"`
	Enabled       bool   `json:"enabled"`
	Comment       string `json:"comment"`
}

type NamedClientRule struct {
	Name string
	Rule ClientRule
}
type ClientWhitelist struct {
	Allow      ClientRule
	Exceptions []NamedClientRule
}

func CompileClientWhitelist(peerPattern, agentPattern string, exceptions []ClientException) (ClientWhitelist, error) {
	allow, err := CompileClientRule(peerPattern, agentPattern)
	if err != nil {
		return ClientWhitelist{}, err
	}
	compiled := ClientWhitelist{Allow: allow}
	for _, exception := range exceptions {
		rule, err := CompileClientRule(exception.PeerIdPattern, exception.AgentPattern)
		if err != nil {
			return ClientWhitelist{}, err
		}
		if exception.Enabled {
			compiled.Exceptions = append(compiled.Exceptions, NamedClientRule{Name: exception.Name, Rule: rule})
		}
	}
	return compiled, nil
}

// Check evaluates only this whitelist and its own exceptions.
func (w ClientWhitelist) Check(peerId, userAgent string) (allowed bool, blockedBy string) {
	if !w.Allow.Match(peerId, userAgent) {
		return false, ""
	}
	for _, exception := range w.Exceptions {
		if exception.Rule.Match(peerId, userAgent) {
			return false, exception.Name
		}
	}
	return true, ""
}
