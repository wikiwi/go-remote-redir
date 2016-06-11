# go-remote-redir
_go-remote-redir_ implements a http server to support [go remote import paths](https://golang.org/cmd/go/#hdr-Remote_import_paths).

## Usage
    Usage:
      robots-disallow [OPTIONS]

    Application Options:
          --listen=        address to listen on (default: 0.0.0.0:8080) [$GRR_LISTEN]
          --pattern=       path pattern (default: /p/(?P<user>[^/]+)/(?P<project>[^/]+).*) [$GRR_PATTERN]
          --meta=          meta tag content for go remote import feature (default: example.io/p/${user}/${project} git ssh://git@gitlab.com/${user}/${project}.git) [$GRR_META]
          --redirect-name= redirect name (default: Gitlab Project Page) [$GRR_REDIRECT_NAME]
          --redirect-to=   redirect to (default: https://gitlab.com/${user}/${project}) [$GRR_REDIRECT_TO]
      -v, --version        show version number

    Help Options:
      -h, --help           Show this help message

## Example
    docker run -p 8080:8080 wikiwi/go-remote-redir

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

