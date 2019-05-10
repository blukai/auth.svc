all: push

APP?=auth.svc
USERSPACE?=blukai
PROJECT?=github.com/${USERSPACE}/${APP}
GOOS?=linux

REPO=$(shell git config --get remote.origin.url)

ifndef COMMIT
	COMMIT := git-$(shell git rev-parse --short HEAD)
endif

build:
	GOOS=${GOOS} go build \
	-ldflags "-X main.Commit=${COMMIT} -X main.Repo=${REPO}" \
	-o ${APP} \
	./cmd/auth

container: build
	docker build --pull -t $(APP):$(COMMIT) .

clean:
	rm -f ${APP}
