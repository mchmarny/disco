ARG BUILD_BASE=golang:1.19.4
ARG FINAL_BASE=cgr.dev/chainguard/alpine-base
ARG VERSION=v0.0.1-default
ARG USER=disco

# BUILD
FROM $BUILD_BASE as builder
WORKDIR /src/
COPY . /src/
ARG VERSION
ENV VERSION=$VERSION GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
    -ldflags="-w -s -X main.version=${VERSION} -extldflags '-static'" \
    -a -mod vendor -o server ./cmd/server/main.go

# RUN
FROM $FINAL_BASE
ARG VERSION
LABEL disco.version="${VERSION}"
COPY --from=builder /src/server /app/
WORKDIR /app
# packages
RUN apk add --update curl bash ca-certificates python3
# gcloud
ENV CLOUDSDK_INSTALL_DIR /gcloud/
RUN curl -sSL https://sdk.cloud.google.com | bash
ENV PATH $PATH:/gcloud/google-cloud-sdk/bin/:/app/bin/
# aquasecurity
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh \
    | sh -s -- -b /usr/local/bin
# debug 
RUN echo $PATH 
RUN ls -al 
# automator
ENTRYPOINT ["./server"]