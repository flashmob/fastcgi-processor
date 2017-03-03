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

So, say good-bye to Web Services. Say hello to Email Services! A new paradigm is here to take the interweb by storm. Well, maybe once a more catchy name is given...

## Usage

Import `"github.com/flashmob/fastcgi-processor"` to your Go-guerrilla project. Import `"github.com/flashmob/go-guerrilla/backends"` 
assuming your have done already, assuming you have imported the go-guerrilla package already.

Somewhere at the top of your code, maybe in your `init()` function, add

`backends.Svc.AddProcessor("FastCGI", fastcgi_processor.Processor)`

This will let Go-Guerrilla know about your FastCGI processor.

See the configuration section for how to configure. Send your configuration to Go-Guerrilla's backends.New() function.


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

Don't forget to add `FastCGI` to the end of your `process_stack` config option, eg:

`"process_stack": "HeadersParser|Debugger|Hasher|Header|FastCGI",`


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


 
