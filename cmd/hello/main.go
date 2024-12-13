package main

import (
	"database/sql"
	_ "encoding/json"
	_ "flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "sandbox"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

func (h *Handlers) GetHello(c echo.Context) error {
	msg, err := h.dbProvider.SelectHello()
	if err != nil {
		return c.String(http.StatusInternalServerError, "error"+err.Error())
	}
	return c.String(http.StatusOK, "Hello"+" "+msg+"!")
}

func (h *Handlers) PostHello(c echo.Context) error {
	input := struct {
		Msg string `json:"msg"`
	}{}
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "error"+err.Error())
	}
	if err := h.dbProvider.InsertHello(input.Msg); err != nil {
		return c.String(http.StatusInternalServerError, "error"+err.Error())
	}
	return c.String(http.StatusCreated, "Created")
}

func (dp *DatabaseProvider) SelectHello() (string, error) {
	var msg string
	row := dp.db.QueryRow("SELECT message FROM hello ORDER BY RANDOM() LIMIT 1")
	err := row.Scan(&msg)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (dp *DatabaseProvider) InsertHello(msg string) error {
	_, err := dp.db.Exec("INSERT INTO hello (message) VALUES ($1)", msg)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	dp := DatabaseProvider{db: db}

	h := Handlers{dbProvider: dp}

	e := echo.New()
	e.GET("/get", h.GetHello)
	e.POST("/post", h.PostHello)

	address := "127.0.0.1:8081"

	if err := e.Start(address); err != nil {
		log.Fatal(err)
	}
}
