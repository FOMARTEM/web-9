package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (h *Handlers) GetCount(c echo.Context) error {
	count, err := h.dbProvider.SelectQuery()
	if err != nil {
		return c.String(http.StatusInternalServerError, "error"+err.Error())
	}
	return c.String(http.StatusOK, strconv.Itoa(count))
}

func (h *Handlers) PostCount(c echo.Context) error {
	countStr := c.FormValue("count")
	count, err := strconv.Atoi(countStr)

	if err != nil {
		return c.String(http.StatusBadRequest, "error Это не число!")
	}

	err = h.dbProvider.InsertQuery(count)

	if err != nil {
		return c.String(http.StatusInternalServerError, "error"+err.Error())
	}
	return c.NoContent(http.StatusCreated)
}

func (dp *DatabaseProvider) SelectQuery() (int, error) {
	var count int

	row := dp.db.QueryRow("SELECT number FROM counts")
	err := row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dp *DatabaseProvider) InsertQuery(count int) error {
	_, err := dp.db.Exec("UPDATE counts SET number = number + $1", count)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	address := "127.0.0.1:8081"

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
	e.GET("/count", h.GetCount)
	e.POST("/count", h.PostCount)

	if err := e.Start(address); err != nil {
		log.Fatal(err)
	}
}
