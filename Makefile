DEPENDS = github.com/yosida95/recvknocking

build: gitpull.go depends
	GOPATH=${PWD} go build -o gitpull ./gitpull.go

depends:
	for DEPEND in $(DEPENDS); do \
		GOPATH=${PWD} go get $$DEPEND; \
	done
