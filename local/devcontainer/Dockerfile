FROM alpine:3.17.0

RUN apk add --no-cache \
	go \
	git \
	make \
	npm \
	openssh

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2

RUN go install github.com/cweill/gotests/gotests@v1.6.0 \
	&& go install github.com/fatih/gomodifytags@v1.16.0 \
	&& go install github.com/josharian/impl@v1.1.0 \
	&& go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0 \
	&& go install github.com/go-delve/delve/cmd/dlv@latest \
	&& go install honnef.co/go/tools/cmd/staticcheck@latest \
	&& go install golang.org/x/tools/gopls@latest \
	&& go install github.com/rhysd/actionlint/cmd/actionlint@latest \
	&& go install golang.org/x/vuln/cmd/govulncheck@latest \
	&& go install github.com/google/pprof@latest

RUN npm i -g nodemon

ENV PATH="$PATH:/root/go/bin"