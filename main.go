package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	var cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Start Go Import Redirector",
		Long: "Starts a HTTP server implementing Go Remote Import Paths. " +
			"See https://golang.org/cmd/go/#hdr-Import_path_syntax.",
		Run: func(cmd *cobra.Command, args []string) {

			h := &Handler{
				regexp.MustCompile(viper.GetString("pattern")),
				viper.GetString("meta"),
				viper.GetString("redirect_name"),
				viper.GetString("redirect_to"),
			}
			fmt.Println("Listening on " + viper.GetString("listen") + "...")
			panic(http.ListenAndServe(viper.GetString("listen"), h))
		},
	}

	viper.SetEnvPrefix("grr")
	viper.AutomaticEnv()

	cmdServe.Flags().String("listen", "0.0.0.0:8080", "address to listen on [$GRR_LISTEN]")
	viper.BindPFlag("listen", cmdServe.Flags().Lookup("listen"))

	cmdServe.Flags().String("pattern", "/p/(?P<user>[^/]+)/(?P<project>[^/]+).*",
		"path pattern [$GRR_PATTERN]")
	viper.BindPFlag("pattern", cmdServe.Flags().Lookup("pattern"))

	cmdServe.Flags().String("meta",
		"example.io/p/${user}/${project} git ssh://git@gitlab.com/${user}/${project}.git",
		"meta tag content for go remote import feature [$GRR_META]")
	viper.BindPFlag("meta", cmdServe.Flags().Lookup("meta"))

	cmdServe.Flags().String("redirectName", "Gitlab Project Page",
		"redirect name [$GRR_REDIRECT_NAME]")
	viper.BindPFlag("redirect_name", cmdServe.Flags().Lookup("redirectName"))

	cmdServe.Flags().String("redirectTo", "https://gitlab.com/${user}/${project}",
		"redirect to [$GRR_REDIRECT_TO]")
	viper.BindPFlag("redirect_to", cmdServe.Flags().Lookup("redirectTo"))

	var rootCmd = &cobra.Command{Use: "go-remote-redir"}
	rootCmd.AddCommand(cmdServe)
	rootCmd.Execute()
}
