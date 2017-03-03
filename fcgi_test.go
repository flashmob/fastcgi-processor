package fcgi_processor

import (
	"github.com/flashmob/go-guerrilla/mail"
	"net/url"
	"strings"
	"testing"
)

// change baseDir to location of your save.php and validate.php scripts
var baseDir string = "/vagrant/projects/golang/src/github.com/flashmob/fastcgi-processor/examples/"

// fcgiType is "unix" or "tcp"
var fcgiType string = "unix"

// fcgiAddr is a path to a unix socket descriptor, or IP address with tcp port eg. "127.0.0.1:8000"
var fcgiAddr string = "/var/run/php/php7.0-fpm.sock"

// Test GET method to validate recipient
// This test requires a functioning php-fpm daemon
func TestGet(t *testing.T) {

	c := &fcgiConfig{
		ScriptFileNameNameSave:     baseDir + "/save.php",
		ScriptFileNameNameValidate: baseDir + "validate.php",
		ConnectionType:             fcgiType,
		ConnectionAddress:          fcgiAddr,
	}

	f, err := newFastCGIProcessor(c)
	if err != nil {
		t.Error("could not newFastCGIPorcessor", err)
	}

	q := url.Values{}
	q.Add("rcpt_to", "test@moo.com")
	result, err := f.get(c.ScriptFileNameNameValidate, q)

	if strings.Index(string(result), "PASSED") != 0 {
		t.Error("save did not return PASSED, it got:", string(result))
	}
}

// Test email saving using POST
// This test requires a functioning php-fpm daemon
func TestPost(t *testing.T) {
	c := &fcgiConfig{
		ScriptFileNameNameSave:     baseDir + "save.php",
		ScriptFileNameNameValidate: baseDir + "validate.php",
		ConnectionType:             fcgiType,
		ConnectionAddress:          fcgiAddr,
	}

	f, err := newFastCGIProcessor(c)

	if err != nil {
		t.Error("could not newFastCGIProcessor", err)
	}

	q := url.Values{}
	q.Add("rcpt_to", "test@moo.com")

	e := &mail.Envelope{
		RemoteIP: "127.0.0.1",
		QueuedId: "abc12345",
		Helo:     "helo.example.com",
		MailFrom: mail.Address{"test", "example.com"},
		TLS:      true,
	}
	e.PushRcpt(mail.Address{"test", "example.com"})
	e.Data.WriteString("Subject:Test\n\nThis is a test.")

	result, err := f.postSave(e)

	if err != nil {
		t.Error("postSave error", err)
	}

	if strings.Index(string(result), "SAVED") != 0 {
		t.Error("save did not return SAVED, it got:", string(result))
	}
}
