package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"

	"github.com/jessevdk/go-flags"
)

var version = "0.0.1-dev"
var tmpl = template.Must(template.New("").Parse(`
<html>
	<head>
		<title>Go Remote Redirect</title>
		<meta name="go-import" content="{{.MetaImport}}">
		<meta http-equiv="refresh" content="0; url={{.RedirectTo}}">
		<meta name="robots" content="noindex">
	</head>
	<body>
		You are being automatically redirected to <a href="{{.RedirectTo}}">{{.RedirectName}}</a>.
	</body>
</html>
`))

var opts struct {
	Listen       string `long:"listen" default:"0.0.0.0:8080" env:"GRR_LISTEN" description:"address to listen on"`
	Pattern      string `long:"pattern" default:"/p/(?P<user>[^/]+)/(?P<project>[^/]+).*" env:"GRR_PATTERN" description:"path pattern"`
	MetaImport   string `long:"meta" default:"example.io/p/${user}/${project} git ssh://git@gitlab.com/${user}/${project}.git" env:"GRR_META" description:"meta tag content for go remote import feature"`
	RedirectName string `long:"redirect-name" default:"Gitlab Project Page" env:"GRR_REDIRECT_NAME" description:"redirect name"`
	RedirectTo   string `long:"redirect-to" default:"https://gitlab.com/${user}/${project}" env:"GRR_REDIRECT_TO" description:"redirect to"`
	Version      bool   `long:"version" short:"v" description:"show version number"`
}

type Handler struct {
	PathPattern  *regexp.Regexp
	MetaImport   string
	RedirectName string
	RedirectTo   string
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if req.URL.Query().Get("go-get") != "1" || !h.PathPattern.MatchString(path) {
		http.NotFound(rw, req)
		return
	}
	tmpl.Execute(rw, struct {
		MetaImport   string
		RedirectName string
		RedirectTo   string
	}{
		h.PathPattern.ReplaceAllString(path, h.MetaImport),
		h.PathPattern.ReplaceAllString(path, h.RedirectName),
		h.PathPattern.ReplaceAllString(path, h.RedirectTo),
	})
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "robots-disallow"
	_, err := parser.Parse()
	if err != nil {
		if e2, ok := err.(*flags.Error); ok && e2.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println(version)
	} else {
		h := &Handler{
			regexp.MustCompile(opts.Pattern),
			opts.MetaImport,
			opts.RedirectName,
			opts.RedirectTo,
		}
		fmt.Println("Listening on " + opts.Listen + "...")
		panic(http.ListenAndServe(opts.Listen, h))
	}
}
