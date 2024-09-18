module github.com/joostvdg/cmg

// +heroku goVersion go1.13
// +heroku install .

require (
	github.com/getsentry/sentry-go v0.29.0
	github.com/go-errors/errors v1.5.1
	github.com/google/uuid v1.6.0
	github.com/kennygrant/sanitize v1.2.4
	github.com/labstack/echo-contrib v0.17.1
	github.com/labstack/echo/v4 v4.12.0
	github.com/prometheus/client_golang v1.20.3
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
	github.com/stretchr/testify v1.9.0
	go.uber.org/automaxprocs v1.5.3
	gopkg.in/segmentio/analytics-go.v3 v3.1.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.59.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/segmentio/backo-go v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

exclude github.com/prometheus/client_golang v0.9.1

exclude github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829

exclude github.com/prometheus/client_golang v1.0.0

exclude github.com/prometheus/client_golang v1.3.0

exclude github.com/prometheus/client_golang v1.4.0

exclude github.com/prometheus/client_golang v1.7.1

exclude github.com/prometheus/client_golang v1.10.0

exclude github.com/prometheus/common v0.9.1

exclude github.com/prometheus/common v0.25.0

exclude github.com/prometheus/common v0.26.0

exclude github.com/labstack/echo/v4 v4.3.0

exclude github.com/labstack/echo/v4 v4.5.0

go 1.21

toolchain go1.23.1
