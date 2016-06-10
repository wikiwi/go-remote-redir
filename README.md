# go-remote-redir
_go-remote-redir_ implements a http server to support [go remote import paths](https://golang.org/cmd/go/#hdr-Remote_import_paths).

## Usage
    go-remote-redir serve [flags]

    Flags:
          --listen string         address to listen on [$GRR_LISTEN] (default "0.0.0.0:8080")
          --meta string           meta tag content for go remote import feature [$GRR_META] (default "example.io/p/${user}/${project} git ssh://git@gitlab.com/${user}/${project}.git")
          --pattern string        path pattern [$GRR_PATTERN] (default "/p/(?P<user>[^/]+)/(?P<project>[^/]+).*")
          --redirectName string   redirect name [$GRR_REDIRECT_NAME] (default "Gitlab Project Page")
          --redirectTo string     redirect to [$GRR_REDIRECT_TO] (default "https://gitlab.com/${user}/${project}")

## Example
    docker run -p 8080:8080 wikiwi/go-remote-redir serve

## Output
    curl localhost:8080/user/project?go-get=1
    <html>
            <head>
                    <title>Go Remote Packages</title>
                    <meta name="go-import" content="example.io/p/user/project git ssh://git@gitlab.com/user/project.git">
                    <meta http-equiv="refresh" content="0; url=https://gitlab.com/user/project">
                    <meta name="robots" content="noindex">
            </head>
            <body>
                    You are being automatically redirected to <a href="https://gitlab.com/user/project">Gitlab Project Page</a>.
            </body>
    </html>

## Docker Hub
Automated build is available at the [Docker Hub](https://hub.docker.com/r/wikiwi/go-import-redir).

