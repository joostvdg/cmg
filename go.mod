module github.com/joostvdg/cmg

// +heroku goVersion go1.13
// +heroku install .

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/getsentry/sentry-go v0.12.0
	github.com/go-errors/errors v1.4.1
	github.com/google/uuid v1.3.0
	github.com/kennygrant/sanitize v1.2.4
	github.com/labstack/echo-contrib v0.11.0
	github.com/labstack/echo/v4 v4.9.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.32.1 // indirect
	github.com/segmentio/backo-go v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.1
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	go.uber.org/automaxprocs v1.5.1
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	gopkg.in/segmentio/analytics-go.v3 v3.1.0
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

go 1.16
