package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
)

type App struct {
	router      http.Handler
	mySQLconfig mysql.Config
	config      Config
	db          *sql.DB
}

func New(config Config) *App {
	app := &App{
		config:      config,
		mySQLconfig: config.MySQLConfig,
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	var err error
	a.db, err = sql.Open("mysql", a.mySQLconfig.FormatDSN())
	if err != nil {
		return fmt.Errorf("could not open mysql database: %w", err)
	}

	defer func() {
		if err := a.db.Close(); err != nil {
			fmt.Println("failed to close mysql", err)
		}
	}()

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("server failed to start: %w", err)
		}

		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
