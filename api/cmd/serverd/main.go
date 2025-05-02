package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"omg/api/cmd/serverd/router"
	"omg/api/internal/authenticate"
	"omg/api/internal/controller/orders"
	"omg/api/internal/controller/products"
	"omg/api/internal/controller/system"
	"omg/api/internal/controller/users"
	"omg/api/internal/repository"
	"omg/api/internal/repository/generator"
	"omg/api/internal/ws"
	"omg/api/pkg/app"
	"omg/api/pkg/db/pg"
	"omg/api/pkg/env"
	"omg/api/pkg/httpserv"

	"github.com/friendsofgo/errors"
)

func main() {
	ctx := context.Background()

	appCfg := app.Config{
		ProjectName:      env.GetAndValidateF("PROJECT_NAME"),
		AppName:          env.GetAndValidateF("APP_NAME"),
		SubComponentName: env.GetAndValidateF("PROJECT_COMPONENT"),
		Env:              app.Env(env.GetAndValidateF("APP_ENV")),
		Version:          env.GetAndValidateF("APP_VERSION"),
		Server:           env.GetAndValidateF("SERVER_NAME"),
		ProjectTeam:      os.Getenv("PROJECT_TEAM"),
	}
	if err := appCfg.IsValid(); err != nil {
		log.Fatal(err)
	}

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Exiting...")
}

func run(ctx context.Context) error {
	log.Println("Starting app initialization")
	dbOpenConns, err := strconv.Atoi(env.GetAndValidateF("DB_POOL_MAX_OPEN_CONNS"))
	if err != nil {
		return errors.WithStack(fmt.Errorf("invalid db pool max open conns: %w", err))
	}
	dbIdleConns, err := strconv.Atoi(env.GetAndValidateF("DB_POOL_MAX_IDLE_CONNS"))
	if err != nil {
		return errors.WithStack(fmt.Errorf("invalid db pool max idle conns: %w", err))
	}

	conn, err := pg.NewPool(env.GetAndValidateF("DB_URL"), dbOpenConns, dbIdleConns)
	if err != nil {
		return err
	}

	defer conn.Close()

	rtr, err := initRouter(ctx, conn)

	log.Println("App initialization completed")

	httpserv.NewServer(rtr.Handler()).Start(ctx)

	return nil
}

func initRouter(
	ctx context.Context,
	dbConn pg.BeginnerExecutor) (router.Router, error) {
	if err := generator.InitSnowflakeGenerators(); err != nil {
		return router.Router{}, err
	}

	return router.New(
		ctx,
		strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		os.Getenv("GQL_INTROSPECTION_ENABLED") == "true",
		system.New(repository.New(dbConn)),
		products.New(repository.New(dbConn)),
		users.New(repository.New(dbConn)),
		orders.New(repository.New(dbConn)),
		authenticate.NewAuthService(repository.New(dbConn), os.Getenv("AUTH_SECRET_KEY")),
		ws.NewHub(),
	), nil
}
