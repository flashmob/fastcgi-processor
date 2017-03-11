# FastCGI processor for the [Go-guerrilla](https://github.com/flashmob/go-guerrilla) package.

FastCGI is an optimized CGI protocol implementation used by web servers to execute scripts or other programs.

Well, it's NOT only for web servers. Yes, you've read that right ;-) 


## About

This package is a _Processor_ for the Go-Guerrilla default `Backend` interface implementation. Typical use for this
package is if you would like to add the ability to deliver emails to your FastCGI backend, using Go-Guerrilla's 
built-in _gateway_ backend. 

Just like a web server would hand over the HTTP request to a FastCGI backend, this plugin
allows you to hand over the processing of an email to your FastCGI backend, such as php-fpm.

The reason why you would do this is because perhaps your codebase is in a scripting language such as PHP,
so there's no need to learn Go, becoming easier for you to maintain, no need to re-compile to change, use your favourite 
framework / library / IDE, etc.

Also, there's no overhead of a web server - it goes straight to your script.


## Usage

Import `"github.com/flashmob/fastcgi-processor"` to your Go-guerrilla project.  
Assuming you have imported the go-guerrilla package already, and all dependencies.

Then, when [using go-guerrilla as a package](https://github.com/flashmob/go-guerrilla/wiki/Using-as-a-package), use something like this

```go


cfg := &AppConfig{
    LogFile:      "stderr",
    AllowedHosts: []string{"example.com"},
    BackendConfig: backends.BackendConfig{
        "save_process" : "HeadersParser|Debugger|FastCGI",
        "validate_process" : "FastCGI",
        "fcgi_script_filename_save" : "/home/path/to/save.php",
        "fcgi_script_filename_validate" : "/home/path/to/validate.php",
        "fcgi_connection_type" : "unix",
        "fcgi_connection_address" : "/tmp/php-fpm.sock"
    },
}
d := Daemon{Config: cfg}
d.AddProcessor("FastCGI", fastcgi_processor.Processor)

d.Start()

// .. keep the server busy..

```


This will let Go-Guerrilla know about your FastCGI processor. Note that here we've
added `FastCGI` to the end of the `save_process` config option, then used the `d.AddProcessor` api
 call to register it. Then configured other settings.

See the configuration section for how to configure. 

## Configuration

The following values are required in your `backend_config` section

```json
"backend_config":{
  "fcgi_script_filename_save" : "/home/path/to/save.php",
  "fcgi_script_filename_validate" : "/home/path/to/validate.php",
  "fcgi_connection_type" : "unix",
  "fcgi_connection_address" : "/tmp/php-fpm.sock"
  // .. other config values
}           


```

`fcgi_connection_type` type can be `unix` or `tcp`. 
`fcgi_connection_address` is a path to a unix socket descriptor, or IP address with tcp port eg. "127.0.0.1:8000"

If `fcgi_connection_address` using the unix socket descriptor, make sure your program has 
permissions for writing to it. The permissions will be tested during initialization.

Don't forget to add `FastCGI` to the end of your `save_process` config option, eg:

`"save_process": "HeadersParser|Debugger|Hasher|Header|FastCGI",`

also, add `FastCGI` to the end of your `validate_process` config option if you want to use the validate script, eg:

`"validate_process": "FastCGI",`

# Scripting

## Validate Recipient Email

A single parameter comes to to your recipient validating script via HTTP GET.

* `rcpt_to` - the email address that we want to validate

Output:

Please echo the string *PASS* and nothing else if validation passed.
Otherwise return anything you wish.

## Save Mail

The parameters comes to to your saving script via a HTTP POST.

The following parameters will be sent:

- `remote_ip` - remote ip address of the client that we got the email from (not the sender)
- `subject` - the subject of the email (if available)
- `tls_on` - boolean, represented as string "true" or "false", was the connection a TLS connection?
- `helo` - hello sent by the client when connecting to us
- `mail_from` - string of the From email address, could be blank to indicate a bounce
- `body` - the raw email body, along with the headers. Please make sure your Fast CGI gateway can handle large enough POST

Output: 

Please echo the string `SAVED` if successful.

## Example

See MailDiranasaurus - it uses this package as an example, https://github.com/flashmob/MailDiranasaurus

## Credits

This package depends on Shen Sheng's [Go Fastcgi client](https://github.com/tomasen/fcgi_client) package.

`go get github.com/sloonz/go-maildir`

## Tips

Your FastCGI script should timeout well before 30 seconds, preferably finish under 1 second.


 
