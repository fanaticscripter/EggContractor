FROM --platform=$BUILDPLATFORM golang:1.16-buster AS builder
WORKDIR /src

RUN curl -sSL https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" >/etc/apt/sources.list.d/yarn.list && \
    curl -sSL https://deb.nodesource.com/setup_14.x | bash - && \
    apt-get update && \
    apt-get install -y --no-install-recommends nodejs yarn
COPY package.json yarn.lock /src/
RUN yarn install
COPY postcss.config.js tailwind.config.js webpack.*.js /src/
COPY css /src/css
COPY js /src/js
COPY templates /src/templates
RUN yarn webpack --config webpack.prod.js && \
    yarn postcss --env=production css/app.css -o static/app.css

COPY . /src/
ARG BUILD
ARG GIT_COMMIT
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags \"-static\" \
    -X github.com/fanaticscripter/EggContractor/web.AppBuild=$BUILD \
    -X github.com/fanaticscripter/EggContractor/web.GitCommit=$GIT_COMMIT"

FROM --platform=$BUILDPLATFORM scratch
WORKDIR /
COPY --from=builder /src/EggContractor /
COPY --from=builder /src/migrations /migrations
COPY --from=builder /src/static /static
COPY --from=builder /src/templates /templates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# The cgo DNS resolver could be problematic in the scratch image.
ENV GODEBUG=netdns=go
ENV EGGCONTRACTOR_CONFIG_FILE=/config.toml
ENTRYPOINT ["/EggContractor"]
CMD ["serve"]
