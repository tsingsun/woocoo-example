package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tsingsun/woocoo"
	"github.com/tsingsun/woocoo/contrib/gql"
	"github.com/tsingsun/woocoo/contrib/telemetry"
	"github.com/tsingsun/woocoo/contrib/telemetry/otelweb"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/pkg/log"
	"github.com/tsingsun/woocoo/web"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func main() {
	app := woocoo.New()
	otelCnf := app.AppConfiguration().Sub("otel")
	telemetry.NewConfig(otelCnf,
		telemetry.WithTracerProviderOptions(zipkinTracer(otelCnf)...),
		telemetry.WithMeterProviderOptions(prometheusProvider(otelCnf)...),
	)

	httpSvr := web.New(web.WithConfiguration(app.AppConfiguration().Sub("web")),
		web.RegisterMiddleware(otelweb.NewMiddleware()),
		web.RegisterMiddleware(gql.New()),
		web.WithGracefulStop(),
	)
	r := httpSvr.Router().Engine
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	r.GET("/abort", func(c *gin.Context) {
		c.Abort()
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	app.RegisterServer(httpSvr)
	if err := app.Run(); err != nil {
		log.Error(err)
	}
}

func zipkinTracer(cnf *conf.Configuration) (opts []trace.TracerProviderOption) {
	exporter, err := zipkin.New(cnf.String("traceExporterEndpoint"))
	if err != nil {
		panic(err)
	}
	return append(opts, trace.WithBatcher(exporter))
}

func prometheusProvider(cnf *conf.Configuration) (opts []metric.Option) {
	exporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}
	return []metric.Option{
		metric.WithReader(exporter),
	}
}
