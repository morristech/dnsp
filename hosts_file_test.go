package dnsp_test

import (
	"bytes"
	"testing"

	"github.com/gophergala/dnsp"
)

func TestHostReader(t *testing.T) {
	t.Parallel()

	i, exp := 0, []string{
		"foo.com",
		"bar.net",
		`^.*\.xxx$`,
		"blocked.com",
		"blocked.net",
		"blocked.org",
		"6.blocked.info",
	}

	(&dnsp.HostsReader{
		Reader: bytes.NewBufferString(`
# Host names, one per line:
foo.com
bar.net  # with comment

# Regular expressions:
*.xxx

# Hosts file lines:
127.0.0.1 blocked.com
127.0.0.1 blocked.net blocked.org
::1 6.blocked.info

1.2.3.4 not-blocked.com
		`)}).ReadFunc(func(host string, rx bool) {
		if i > len(exp) {
			t.Errorf("unexpected host read: %q", host)
			return
		}

		if exp[i] != host {
			t.Errorf("expected %q, got %q", exp[i], host)
		}

		if exp, act := (host[0] == '^'), rx; exp != act {
			t.Errorf("expected (%q, %v), got %v", host, exp, act)
		}

		i++
	})

	if i < len(exp) {
		t.Errorf("hosts files not read: %q", exp[i:])
	}
}
