package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inception/internal/payment"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {
	initConfig()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("cannot open an SQLite memory database: %v", err)
	}
	defer db.Close()
	fmt.Println("init")
	initDB(db)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	initRoute(e, db)

	go StartServer(e)
	waitForGracefulShutdown(e)
}

func StartServer(e *echo.Echo) {
	if err := e.Start(":8080"); err != nil {
		panic("shutdown server")
	} else {
		fmt.Println("start server")
	}
}

func initDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE transactions (amount INTEGER, currency TEXT,token TEXT,status TEXT);")
	if err != nil {
		log.Fatalf("cannot create schema: %v", err)
	}

}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config File Not Found")
		} else {
			panic(fmt.Errorf("error config file : %s", err))
		}
	}
}

func initRoute(e *echo.Echo, db *sql.DB) {
	payment := payment.NewHandler(payment.NewService(payment.NewRepository(db)))

	e.GET("/health", health)

	e.POST("/payment", payment.Payment)
	e.POST("/inquiry", payment.Inquiry)
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func waitForGracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
