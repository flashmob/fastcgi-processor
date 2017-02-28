package fcgi_processor

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/tomasen/fcgi_client"
	"github.com/flashmob/go-guerrilla/response"
	"io/ioutil"

	"strconv"

)


type fcgiConfig struct {
	// full path to script for the save mail task
	// eg. /home/user/scripts/save.php
	ScriptFileNameNameSave     string `json:"fcgi_script_filename_save"`
	// full path to script for recipient validation
	// eg /home/user/scripts/val_rcpt.php
	ScriptFileNameNameValidate string `json:"fcgi_script_filename_validate"`
	// "tcp" or "unix"
	ConnectionType string `json:"fcgi_connection_type"`
	// where to Dial, eg "/tmp/php-fpm.sock" for unix-socket or "127.0.0.1:9000" for tcp
	ConnectionAddress string `json:"fcgi_connection_address"`
}



type FastCGIProcessor struct {

	config  *fcgiConfig
	client *fcgiclient.FCGIClient
}


func newFastCGIProcessor(config *fcgiConfig) (*FastCGIProcessor, error) {
	p := &FastCGIProcessor{}
	p.config = config
	err := p.connect()
	if err != nil {
		backends.Log().Debug("FastCgi error", err)
		return p, err
	}
	return p, err
}

func (f *FastCGIProcessor) connect() (err error) {
	f.client, err = fcgiclient.Dial(f.config.ConnectionType, f.config.ConnectionAddress)
	return err
}

// get sends a get query to script with q query values
func (f *FastCGIProcessor) get(script string, q url.Values) (result []byte, err error) {

	env := make(map[string]string)
	env["SCRIPT_FILENAME"] = script
	env["SERVER_SOFTWARE"] = "Go-guerrilla fastcgi"
	env["REMOTE_ADDR"] = "127.0.0.1"
	env["QUERY_STRING"] = q.Encode()

	resp, err := f.client.Get(env)
	if err != nil {
		backends.Log().Debug("FastCgi Get failed", err)
		return result, err
	}

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		backends.Log().Debug("FastCgi Read Body failed", err)
		return result, err
	}

	return result, nil

}

func (f *FastCGIProcessor) postSave(e *mail.Envelope) (result []byte, err error) {
	env := make(map[string]string)
	env["SCRIPT_FILENAME"] = f.config.ScriptFileNameNameSave
	env["SERVER_SOFTWARE"] = "Go-guerrilla fastcgi"
	env["REMOTE_ADDR"] = "127.0.0.1"
	data := url.Values{}

	for i := range e.RcptTo {
		data.Set(fmt.Sprintf("rcpt_to_%d", i), e.RcptTo[i].String())
	}
	data.Set("remote_ip", e.RemoteIP)
	data.Set("subject", e.Subject)
	data.Set("tls_on", strconv.FormatBool(e.TLS))
	data.Set("helo", e.Helo)
	data.Set("mail_from", e.MailFrom.String())
	data.Set("body", e.String())

	resp, err := f.client.PostForm(env, data)
	if err != nil {
		return result, err
	}

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		backends.Log().Debug("FastCgi Read Body failed", err)
		return result, err
	}

	return

	/*
	todo: figure out how we can call directly and use a reader for efficiency, eg.

		r := io.MultiReader(
			bytes.NewReader([]byte("---------------------------974767299852498929531610575\r\n")),
			// ..url encoded data here
			// ..a boundary here
			e.NewReader(),

		/*


		f.client.Post(env, "multipart/form-data; boundary=---------------------------974767299852498929531610575", e.NewReader(), e.Len())
	*/
}


var Processor = func() backends.Decorator {

	// The following initialization is run when the program first starts

	// config will be populated by the initFunc
	var (
		p *FastCGIProcessor
	)
	// initFunc is an initializer function which is called when our processor gets created.
	// It gets called for every worker
	initializer := backends.InitializeWith(func(backendConfig backends.BackendConfig) error {
		configType := backends.BaseConfig(&fcgiConfig{})
		bcfg, err := backends.Svc.ExtractConfig(backendConfig, configType)

		if err != nil {
			return err
		}
		c := bcfg.(*fcgiConfig)
		p, err = newFastCGIProcessor(c)
		if err != nil {
			return err
		}
		p.config = c
		return nil
	})
	// register our initializer
	backends.Svc.AddInitializer(initializer)



	return func(c backends.Processor) backends.Processor {
		// The function will be called on each email transaction.
		// On success, it forwards to the next step in the processor call-stack,
		// or returns with an error if failed
		return backends.ProcessWith(func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
			if task == backends.TaskValidateRcpt {
				// Check the recipients for each RCPT command.
				// This is called each time a recipient is added,
				// validate only the _last_ recipient that was appended
				if size := len(e.RcptTo); size > 0 {
					v := url.Values{}
					v.Set("rcpt_to", e.RcptTo[len(e.RcptTo)-1].String())
					result, err := p.get(p.config.ScriptFileNameNameValidate, v)
					if err != nil {
						backends.Log().Debug("FastCgi error", err)
						return backends.NewResult(
							response.Canned.FailNoSenderDataCmd),
							backends.StorageNotAvailable
					}
					if string(result[0:4]) == "PASS" {
						// validation passed
						return c.Process(e, task)
					} else {
						// validation failed
						backends.Log().Debug("FastCgi Read Body failed", err)
						return backends.NewResult(
							response.Canned.FailNoSenderDataCmd),
							backends.StorageNotAvailable
					}


					return c.Process(e, task)

				}
				return c.Process(e, task)
			} else if task == backends.TaskSaveMail {
				for i := range e.RcptTo {
					// POST to FCGI
					resp, err := p.postSave(e)
					if err != nil {

					} else if strings.Index(string(resp), "PASS") == 0 {
						return c.Process(e, task)
					} else {
						backends.Log().WithError(err).Error("Could not save email")
						return backends.NewResult(fmt.Sprintf("554 Error: could not save email for [%s]", e.RcptTo[i])), err
					}
				}
				// continue to the next Processor in the decorator chain
				return c.Process(e, task)
			} else {
				return c.Process(e, task)
			}

		})
	}
}
