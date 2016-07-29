# Copyright (C) 2016 wikiwi.io
#
# This software may be modified and distributed under the terms
# of the MIT license. See the LICENSE file for details.

### GitHub settings ###
GITHUB_USER ?= wikiwi
GITHUB_REPO ?= go-remote-redir

### Docker settings ###
DOCKER_REPO    ?= wikiwi/go-remote-redir
LATEST_VERSION ?=

### GO settings ###
GO_PACKAGE  ?= github.com/${GITHUB_USER}/${GITHUB_REPO}

# Glide options
GLIDE_OPTS ?=
GLIDE_GLOBAL_OPTS ?=

### CI Settings ###
# Set branch with most current HEAD of master e.g. master or origin/master.
# E.g. Gitlab doesn't pull the master branch but fetches it to origin/master.
MASTER_BRANCH ?= master

### Build Tools ###
GO ?= go
GLIDE ?= glide
GIT ?= git
DOCKER ?= docker
GITHUB_RELEASE ?= github-release

### Environment ###
HAS_GLIDE := $(shell command -v ${GLIDE};)
HAS_GIT := $(shell command -v ${GIT};)
HAS_GO := $(shell command -v ${GO};)
GOOS := $(shell ${GO} env GOOS)
GOARCH := $(shell ${GO} env GOARCH)
BINARY := go-remote-redir

# Load versioning logic.
include Makefile.versioning

# Docker Image info.
IMAGE := ${DOCKER_REPO}:${BUILD_REF}

# Show build info.
info:
	@echo "Version: ${BUILD_VERSION}"
	@echo "Image:   ${IMAGE}"
	@echo "Tags:    ${TAGS}"

.PHONY: build
ifneq (${GOOS}, "windows")
build: bin/${GOOS}/${GOARCH}/${BINARY}
else
build: bin/${GOOS}/${GOARCH}/${BINARY}.exe
endif

.PHONY: build-cross
build-cross: bin/linux/amd64/${BINARY} bin/freebsd/amd64/${BINARY} bin/darwin/amd64/${BINARY} bin/windows/amd64/${BINARY}.exe
	$(NOOP)

.PHONY: build-for-docker
build-for-docker: bin/linux/amd64/${BINARY}

# docker-build will build the docker image.
.PHONY: docker-build
docker-build: build-for-docker
	${DOCKER} build --pull -t ${IMAGE} \
		--build-arg "BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"`" \
		--build-arg "VCS_REF=${GIT_SHA}" \
		--build-arg "VCS_VERSION=${BUILD_VERSION}" \
		--build-arg "VCS_MESSAGE=$$(git log --oneline -n1 --pretty=%B | head -n1)" \
		--build-arg "BUILD_URL=$$(test -n "$${TRAVIS_BUILD_ID}" && echo https://travis-ci.org/${GITHUB_USER}/${GITHUB_REPO}/builds/$${TRAVIS_BUILD_ID})" \
		.

# docker-push will push all tags to the repository
.PHONY: docker-push
docker-push: ${TAGS:%=docker-push-%}
docker-push-%:
	${DOCKER} tag ${IMAGE} ${DOCKER_REPO}:$* && docker push ${DOCKER_REPO}:$*

.PHONY: has-tags
has-tags:
ifndef TAGS
	@echo No tags set for this build
	false
endif

# clean deletes build artifacts from the project.
.PHONY: clean
clean:
	rm -rf bin
	rm .coverprofile

# test will start the project test suites.
.PHONY: test
test:
	echo Running unit tests
	go test -i && go test

.PHONY: test-with-coverage
test-with-coverage:
	go test -coverprofile=.coverprofile

.PHONY: github-release
github-release:
ifdef IS_DIRTY
	$(error Current trunk is marked dirty)
endif
ifndef IS_RELEASE
	@echo "Skipping release as this commit is not tagged as one"
else
	${GITHUB_RELEASE} release -u "${GITHUB_USER}" -r "${GITHUB_REPO}" -t "${GIT_TAG}" -n "${GIT_TAG}" $$(test -n "${VERSION_STAGE}" && echo --pre-release) || true
endif

# bootstrap will install project dependencies.
.PHONY: bootstrap
bootstrap:
ifndef HAS_GO
	$(error You must install Go)
endif
ifndef HAS_GIT
	$(error You must install Git)
endif
ifndef HAS_GLIDE
	${GO} get -u github.com/Masterminds/glide
endif
	${GLIDE} ${GLIDE_GLOBAL_OPTS} install ${GLIDE_OPTS}

include Makefile.build

