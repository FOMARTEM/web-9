package main

import (
	"database/sql"
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

func (h *Handlers) GetQuery(c echo.Context) error {
	name := c.QueryParam("name")

	if name == "" {
		return c.String(http.StatusBadRequest, "Hello, stranger!")
	}

	name, err := h.dbProvider.SelectName(name)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.String(http.StatusInternalServerError, "Such user does not exist!")
		}

		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Hello, "+name+"!")
}

func (h *Handlers) PostQuery(c echo.Context) error {
	name := c.QueryParam("name")

	if name == "" {
		return c.String(http.StatusBadRequest, "Hello, stranger!")
	}

	err := h.dbProvider.InsertName(name)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "name added")
}

func (dp *DatabaseProvider) SelectName(name string) (string, error) {
	var nam string

	row := dp.db.QueryRow("SELECT name FROM usernames WHERE name = ($1)", name)
	err := row.Scan(&nam)

	if err != nil {
		return "", err
	}

	return nam, nil
}

func (dp *DatabaseProvider) InsertName(name string) error {
	_, err := dp.db.Exec("INSERT INTO usernames (name) VALUES ($1)", name)

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
	e.GET("/api/user/get", h.GetQuery)
	e.POST("/api/user/post", h.PostQuery)

	address := "127.0.0.1:8081"

	if err := e.Start(address); err != nil {
		log.Fatal(err)
	}
}
