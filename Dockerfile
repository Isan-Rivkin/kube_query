FROM golang:1.15-alpine AS build_kq

RUN apk add --update alpine-sdk git make && \
	git config --global http.https://gopkg.in.followRedirects true 

WORKDIR /app

COPY . .

RUN make build-linux && \
	mv /app/kq_linux /app/kq

RUN apk add --update alpine-sdk make 

WORKDIR /app

COPY . .

FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_kq /app/kq /bin/kq

ENTRYPOINT ["/bin/kq"]
