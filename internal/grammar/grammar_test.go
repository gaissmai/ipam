package grammar_test

import (
	"testing"

	g "github.com/gaissmai/ipam/internal/grammar"
)

func TestGrammar(t *testing.T) {
	t.Parallel()

	for name := range g.AllRules() {
		rx := g.MustRx(name)
		t.Log(name, rx)
	}
}

func TestGrammarRules(t *testing.T) {
	t.Parallel()

	tt := []struct {
		rule  string
		input string
		want  bool
	}{
		// ── word ────────────────────────────────────────────────────────────
		{rule: "word", input: "hello", want: true},
		{rule: "word", input: "with-hyphen", want: true},
		{rule: "word", input: "word123", want: true},
		{rule: "word", input: "_underscore", want: true},
		{rule: "word", input: "", want: false},

		// ── cswords ─────────────────────────────────────────────────────────
		{rule: "cswords", input: "single", want: true},
		{rule: "cswords", input: "std,open", want: true},
		{rule: "cswords", input: "a,b-c,d", want: true},
		{rule: "cswords", input: ",leading,comma", want: false},
		{rule: "cswords", input: "trailing,comma,", want: false},
		{rule: "cswords", input: "a,,b", want: false}, // double comma
		{rule: "cswords", input: "-starts,with,hyphen", want: false},
		{rule: "cswords", input: "with space", want: false}, // spaces not allowed

		// ── comment ─────────────────────────────────────────────────────────
		{rule: "comment", input: "# this is a comment", want: true},
		{rule: "comment", input: "#", want: true},
		{rule: "comment", input: "# ", want: true},
		{rule: "comment", input: "not a comment", want: false},
		{rule: "comment", input: "", want: false},
		{rule: "comment", input: "  # indented", want: false}, // ^ anchors at start

		// ── scoping ─────────────────────────────────────────────────────────
		{rule: "scoping", input: "---", want: true},
		{rule: "scoping", input: "----", want: true},
		{rule: "scoping", input: "----------", want: true},
		{rule: "scoping", input: "--", want: false},    // minimum three dashes
		{rule: "scoping", input: "- -", want: false},   // no spaces
		{rule: "scoping", input: "--- x", want: false}, // nothing after dashes

		// ── fqdn ────────────────────────────────────────────────────────────
		{rule: "fqdn", input: "www.example.com", want: true},
		{rule: "fqdn", input: "www.example.com.", want: true}, // trailing dot allowed
		{rule: "fqdn", input: "*.foo.bar", want: true},        // wildcard prefix
		{rule: "fqdn", input: "example.com", want: true},
		{rule: "fqdn", input: "a.b.c.d.e.tld", want: true},
		{rule: "fqdn", input: "localhost", want: false},   // bare hostname, no dot
		{rule: "fqdn", input: "192.168.1.1", want: false}, // TLD must end with letter
		{rule: "fqdn", input: "", want: false},
		{rule: "fqdn", input: "foo bar.com", want: false}, // spaces not allowed in label

		// ── ident (no trailing $, anchored only at start) ───────────────────
		{rule: "ident", input: "VLAN", want: true},
		{rule: "ident", input: "VRF", want: true},
		{rule: "ident", input: "STATUS_OK", want: true},
		{rule: "ident", input: "VLAN10", want: true},
		{rule: "ident", input: "A9", want: true},
		{rule: "ident", input: "A", want: false}, // [A-Z0-9_]+ needs ≥1 more char
		{rule: "ident", input: "lowercase", want: false},
		{rule: "ident", input: "Mixed", want: false},
		{rule: "ident", input: "1NVALID", want: false}, // must start with A-Z

		// ── ip (no trailing $, anchored only at start) ───────────────────────
		{rule: "ip", input: "192.168.1.1", want: true},
		{rule: "ip", input: "192.168.1.1/24", want: true}, // prefix is fine, no trailing $
		{rule: "ip", input: "10.0.0.1", want: true},
		{rule: "ip", input: "2001:db8::1", want: true},
		{rule: "ip", input: "::1", want: true},
		{rule: "ip", input: "not-an-ip", want: false},
		{rule: "ip", input: "", want: false},

		// ── ipv4 (no trailing $) ─────────────────────────────────────────────
		{rule: "ipv4", input: "192.168.1.1", want: true},
		{rule: "ipv4", input: "0.0.0.0", want: true},
		{rule: "ipv4", input: "255.255.255.255", want: true},
		{rule: "ipv4", input: "192.168.1.1/24", want: true}, // CIDR — no trailing $ intended
		{rule: "ipv4", input: "2001:db8::1", want: false},   // IPv6 must not match ipv4
		{rule: "ipv4", input: "1.2.3", want: false},         // only three octets
		{rule: "ipv4", input: "not-an-ip", want: false},

		// ── ipv6 (no trailing $) ─────────────────────────────────────────────
		{rule: "ipv6", input: "2001:db8::1 foo", want: true},
		{rule: "ipv6", input: "::1", want: true},
		{rule: "ipv6", input: "fe80::1", want: true},
		{rule: "ipv6", input: "2001:db8::1/64", want: true}, // prefix — no trailing $ intended
		{rule: "ipv6", input: "192.168.1.1", want: false},   // pure dotted-decimal has no colon
		{rule: "ipv6", input: "not-an-ip", want: false},
	}

	for _, tc := range tt {
		t.Run(tc.rule+"/"+tc.input, func(t *testing.T) {
			t.Parallel()
			got := g.MustRx(tc.rule).MatchString(tc.input)
			if got != tc.want {
				t.Fatalf("rule=%q input=%q: got %v, want %v", tc.rule, tc.input, got, tc.want)
			}
		})
	}
}
