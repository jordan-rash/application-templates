package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	services "github.com/nats-io/nats.go/micro"
)

var (
	VERSION string = "0.0.1"

	upstreamServiceSubjects = []string{"todo.*.create", "todo.*.read", "todo.*.update", "todo.*.delete"}
)

type MyContext struct {
	ctx        context.Context
	NatsServer string
	NatsNkey   string
	NatsJwt    string
}

func main() {
	mCtx, err := setMyContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting context: %v\n", err)
		os.Exit(1)
	}

	nc, err := nats.Connect(mCtx.NatsServer,
		nats.Name("todo-apigateway"),
		nats.UserJWTAndSeed(mCtx.NatsJwt, mCtx.NatsNkey),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to NATS server: %v\n", err)
		return
	}
	setupSignalHandlers(nc)

	fmt.Fprintln(os.Stdout, "Starting TODO API Gateway service")
	_, err = services.AddService(nc, services.Config{
		Name:    "TodoCreate",
		Version: VERSION,
		Endpoint: &services.EndpointConfig{
			Subject: "api.create",
			Handler: services.HandlerFunc(createHandler(nc)),
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding service: %v\n", err)
		return
	}
	_, err = services.AddService(nc, services.Config{
		Name:    "TodoRead",
		Version: VERSION,
		Endpoint: &services.EndpointConfig{
			Subject: "api.read",
			Handler: services.HandlerFunc(readHandler(nc)),
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding service: %v\n", err)
		return
	}
	_, err = services.AddService(nc, services.Config{
		Name:    "TodoUpdate",
		Version: VERSION,
		Endpoint: &services.EndpointConfig{
			Subject: "api.update",
			Handler: services.HandlerFunc(updateHandler(nc)),
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding service: %v\n", err)
		return
	}
	_, err = services.AddService(nc, services.Config{
		Name:    "TodoDelete",
		Version: VERSION,
		Endpoint: &services.EndpointConfig{
			Subject: "api.delete",
			Handler: services.HandlerFunc(deleteHandler(nc)),
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding service: %v\n", err)
		return
	}
	_, err = services.AddService(nc, services.Config{
		Name:    "TodoHealthCheck",
		Version: VERSION,
		Endpoint: &services.EndpointConfig{
			Subject: "api.healthz",
			Handler: services.HandlerFunc(healthCheckHandler(nc)),
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding service: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, "TODO API Gateway successfully started")
	<-mCtx.ctx.Done()
	fmt.Fprintln(os.Stdout, "TODO API Gateway successfully exited")
}

func setMyContext() (*MyContext, error) {
	var found bool
	mc := new(MyContext)
	mc.ctx = context.Background()

	mc.NatsServer, found = os.LookupEnv("NEX_HOSTSERVICES_NATS_SERVER")
	if !found {
		return mc, fmt.Errorf("NEX_HOSTSERVICES_NATS_SERVER not set")
	}
	mc.NatsNkey, found = os.LookupEnv("NEX_HOSTSERVICES_NATS_USER_SEED")
	if !found {
		return mc, fmt.Errorf("NEX_HOSTSERVICES_NATS_USER_SEED not set")
	}
	mc.NatsJwt, found = os.LookupEnv("NEX_HOSTSERVICES_NATS_USER_JWT")
	if !found {
		return mc, fmt.Errorf("NEX_HOSTSERVICES_NATS_USER_JWT not set")
	}
	return mc, nil
}

func setupSignalHandlers(nc *nats.Conn) {
	go func() {
		signal.Reset(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		for {
			switch s := <-c; {
			case s == syscall.SIGTERM || s == os.Interrupt || s == syscall.SIGQUIT:
				fmt.Fprintf(os.Stdout, "Caught signal [%s], requesting clean shutdown", s.String())
				_ = nc.Drain()
				os.Exit(0)

			default:
				_ = nc.Drain()
				os.Exit(0)
			}
		}
	}()
}

type healthCheck struct {
	Errors error `json:"errors"`
}

func healthCheckHandler(nc *nats.Conn) func(req services.Request) {
	return func(req services.Request) {
		var errs error
		for _, s := range upstreamServiceSubjects {
			_, err := nc.Request(s, []byte("ping"), time.Second)
			errs = errors.Join(errs, err)
		}
		err := req.RespondJSON(healthCheck{Errors: errs})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to health check: %v\n", err)
		}
	}
}
