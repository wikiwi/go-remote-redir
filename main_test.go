package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var fixPattern = "/(?P<user>[^/]+)/(?P<project>[^/]+).*"
var fixMeta = "example.io/${user}/${project} git ssh://git@gitlab.com/${user}/${project}.git"
var fixRedirectName = "Test"
var fixRedirectTo = "https://test.com/${user}/${project}"

func getFreePort() string {
	l, _ := net.Listen("tcp", "localhost:0")
	defer l.Close()
	return strings.Split(l.Addr().String(), ":")[1]
}

func runServer() {
	main()
	panic("main routine exited")
}

func TestServe(t *testing.T) {
	testScenarios := []struct {
		inPattern, inMeta, inRedirectName, inRedirectTo string
		requestPath                                     string
		statusCode                                      int
		goImport, httpEquiv                             string
	}{
		{
			inPattern:      "/(?P<user>[^/]+)/(?P<project>[^/]+).*",
			inMeta:         "example.io/${user}/${project} git ssh://git@gitlab.com/${user}/${project}.git",
			inRedirectName: "Test",
			inRedirectTo:   "https://test.com/${user}/${project}",
			requestPath:    "/user/project?go-get=1",
			statusCode:     200,
			goImport:       "example.io/user/project git ssh://git@gitlab.com/user/project.git",
			httpEquiv:      "0; url=https://test.com/user/project",
		},
		{
			statusCode:  404,
			requestPath: "/user/project",
		},
	}
	for _, x := range testScenarios {
		addr := "localhost:" + getFreePort()
		os.Args = []string{"go-remote-redir", "--listen", addr, "--pattern", x.inPattern, "--meta",
			x.inMeta, "--redirect-name", x.inRedirectName, "--redirect-to", x.inRedirectTo}
		go runServer()

		time.Sleep(time.Second)

		resp, err := http.Get("http://" + addr + x.requestPath)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != x.statusCode {
			t.Fatalf("code was %d != %d", resp.StatusCode, x.statusCode)
		}

		if resp.StatusCode != 200 {
			continue
		}

		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			t.Fatal(err)
		}

		content, _ := doc.Find("meta[name=go-import]").Attr("content")
		if content != x.goImport {
			t.Errorf("go-import meta was %q, expected %q", content, x.goImport)
		}
		content, _ = doc.Find("meta[http-equiv=refresh]").Attr("content")
		if content != x.httpEquiv {
			t.Errorf("http-equiv meta was %q, expected %q", content, x.httpEquiv)
		}
		content, _ = doc.Find("meta[name=robots]").Attr("content")
		if content != "noindex" {
			t.Errorf("robots meta was %q, expected %q", content, "noindex")
		}
	}
}

func TestVersion(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	output := make(chan string)
	go func() {
		bytes, _ := ioutil.ReadAll(r)
		output <- strings.TrimSpace(string(bytes))
	}()

	os.Args = []string{"go-remote-redir", "--version"}
	main()
	w.Close()
	os.Stdout = stdout

	content := <-output
	if content != version {
		t.Fatalf("%q != %q", content, version)
	}
}
