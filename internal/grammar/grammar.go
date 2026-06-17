// Package grammar defines patterns and rules for IPAM domain specific language.
package grammar

import (
	"regexp"

	"github.com/gaissmai/grammar"
)

// g is ready to use after init().
//
//nolint:gochecknoglobals
var g *grammar.Grammar

// init compiles the grammar rules.
func init() {
	g = grammar.New("IPAM")

	for name, rawPattern := range AllRules() {
		if err := g.Add(name, rawPattern); err != nil {
			panic(err)
		}
	}

	if err := g.Compile(); err != nil {
		panic(err)
	}
}

// MustRx returns the regexp for rule from the underlying grammar.
func MustRx(name string) *regexp.Regexp {
	rx, err := g.Rx(name)
	if err != nil {
		panic(err)
	}

	return rx
}

// AllRules to build the IPAM tokens.
func AllRules() map[string]string {
	return map[string]string{
		"label":   label,
		"word":    word,
		"cswords": cswords,
		"fqdn":    fqdn,
		"ident":   ident,
		"comment": comment,
		"scoping": scoping,
		"ip":      ip,
		"ipv4":    ipv4,
		"ipv6":    ipv6,
	}
}

// ###########################################################################
//                      IPAM GRAMMAR
// slimmed down version, just some helper regexps for the lexer state machine
// ###########################################################################

const (

	// # comment til EOL
	comment = `^ # .* $`

	scoping = `^ --- -* $`

	// label is a field of a FQDN
	label = `(?: [- \w]+ )` // hyphen and word chars, the rest is done by the parser

	// word
	word = `\w [-\w]*`

	// comma-separated words
	cswords = `^ ${word} (?: , ${word} )* $`

	// full qualified domain name, com  com. www.google.com. *.foo.bar
	// the length restriction parsing etc. is done by the parser
	// IDNA not supported, just plain ASCII
	fqdn = `^
           (?:
             (?: \*\. )?           // alias can start with wildcard followed by dot
             (?: ${label} \. )+    // one or more (label followed by dot)
             ${label}? [a-zA-Z]+   // TLD ends with at least one ASCII letter
             \.?                   // last dot optional
          $)`

	// keywords, attributes and flags are uppercase
	ident = `^ (?: [A-Z] [A-Z0-9_]+ )` // missing $ is intended

	// simple and fast patterns for lexing, true parsing with netip.ParseAddr()
	ip = `^ (?: ${ipv4} | ${ipv6} ) ` // missing $ is intended

	ipv4 = `^ (?:
              \d{1,3} \.
              \d{1,3} \.
              \d{1,3} \.
              \d{1,3}
            )` // missing $ is intended

	ipv6 = `^ (?:                // xdigits and colons, no dot
             [[:xdigit:] :]+ :
             [[:xdigit:] :]*
            )` // missing $ is intended
)
