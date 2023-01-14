package main

import (
	"errors"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"io"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	frontendContentsPath = "../public"
	imagesPath = "../public/images"
	mysqlErrNumDuplicateEntry   = 1062
)

var (
	db                  *sqlx.DB
	mySQLConnectionData *MySQLConnectionEnv
)

type MySQLConnectionEnv struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

type Account struct {
	AccountID   int    `db:"account_id"`
	LoginName string `db:"login_name"`
	ShadowPassword string `db:"shadow_password"`
}
type Accountlist []Account

type Event struct {
	EventID   int    `db:"event_id"`
	AccountID   int    `db:"account_id"`
	Title string `db:"title"`
	Description  string `db:"description"`
	EventDate time.Time  `db:"event_date"`
}
type Eventlist []Event
type EventDatail struct {
	Event	*Event		`json:"event"`
	Persons	*Personlist	`json:"persons"`
}
type NewEvent struct {
	Title  string   `json:"title"`
	Description string `json:"description"`
	EventDate      int64  `json:"event_date"`
}
type Person struct {
	PersonId   int    `db:"person_id"`
	FirstName string `db:"first_name"`
	LastName string `db:"last_name"`
}
type Personlist []Person

func getEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}

func NewMySQLConnectionEnv() *MySQLConnectionEnv {
	return &MySQLConnectionEnv{
		Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
		Port:     getEnv("MYSQL_PORT", "3306"),
		User:     getEnv("MYSQL_USER", "poc_webapp"),
		DBName:   getEnv("MYSQL_DBNAME", "poc_webapp"),
		Password: getEnv("MYSQL_PASS", "poc_webapp"),
	}
}

func (mc *MySQLConnectionEnv) ConnectDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Asia%%2FTokyo", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)
	fmt.Println(dsn)
	return sqlx.Open("mysql", dsn)
}

func main() {
	// Echo instance
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/api/test", hello)
	e.GET("/api/accounts", accounts)

	e.GET("/api/events", getEventList)
	e.GET("/api/events/:event_id", getEvent)
	e.POST("/api/events", postEvent)
	e.DELETE("/api/events/:event_id", deleteEvent)

	e.GET("/api/persons", getPersonList)
//	e.GET("/api/persons/:person_id", getPerson)
//	e.POST("/api/persons", postPerson)
//	e.DELETE("/api/persons/:person_id", deletePerson)

	e.POST("/api/images", uploadImage)
//	e.DELETE("/api/images/:image_id", deleteImage)

	e.GET("/", getIndex)
	e.GET("/home", getIndex)
	e.Static("/assets", frontendContentsPath+"/assets")
	e.Static("/images", imagesPath)

	mySQLConnectionData = NewMySQLConnectionEnv()

	var err error
	db, err = mySQLConnectionData.ConnectDB()
	if err != nil {
		e.Logger.Fatalf("failed to connect db: %v", err)
		return
	}
	db.SetMaxOpenConns(10)
	defer db.Close()

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func accounts(c echo.Context) error {
	rows, err := db.Queryx(`select account_id, login_name, shadow_password from accounts`)
	if err != nil {
		c.Logger().Errorf("failed to query: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	var account Account
	var accountlist Accountlist
	for rows.Next() {
		err := rows.StructScan(&account) //sqlのrows.Scanの代わりにsqlxのrows.StructScanを使う
		if err != nil {
			c.Logger().Errorf("failed to purse query responce: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		accountlist = append(accountlist, account)
	}
	return c.JSON(http.StatusOK, accountlist)
}
// GET api/events/
// イベントリストの取得
func getEventList(c echo.Context) error {
	limit := c.QueryParam("limit")
	if limit != "" {
		limit = fmt.Sprintf(" limit "+limit)
	}
	offset := c.QueryParam("offset")
	if offset != "" {
		offset = fmt.Sprintf(" offset "+offset)
	}
	rows, err := db.Queryx(`select event_id, title, description, event_date from events` + limit + offset)
	if err != nil {
		c.Logger().Errorf("failed to query: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	var event Event
	var eventlist Eventlist
	for rows.Next() {
		err := rows.StructScan(&event) //sqlのrows.Scanの代わりにsqlxのrows.StructScanを使う
		if err != nil {
			c.Logger().Errorf("failed to purse query responce: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		eventlist = append(eventlist, event)
	}

	return c.JSON(http.StatusOK, eventlist)
}
// GET api/events/{event_id}
// 個々のイベントの取得（参加者や画像URL等の詳細情報付き）
func getEvent(c echo.Context) error {
	eventID := c.Param("event_id")

	var event Event
	err := db.Get(&event, "SELECT * FROM `events` WHERE `event_id` = ?",eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.String(http.StatusNotFound, "not found: event")
		}
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	rows, err := db.Queryx("select * from persons where person_id in (select distinct(person_id) from event_person_tagging where event_id =?);",eventID)
	if err != nil {
		c.Logger().Errorf("failed to query: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	var person Person
	var personlist Personlist
	for rows.Next() {
		err := rows.StructScan(&person) //sqlのrows.Scanの代わりにsqlxのrows.StructScanを使う
		if err != nil {
			c.Logger().Errorf("failed to purse query responce: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		personlist = append(personlist, person)
	}

	var res EventDatail
	res = EventDatail{
		Event:	&event,
		Persons:	&personlist,
	}

	return c.JSON(http.StatusOK, res)
}
// POST /api/events
// Eventを登録
func postEvent(c echo.Context) error {
	event := new(NewEvent)
	if err := c.Bind(event);err != nil {
		// error handling
	}
	title:= event.Title
	description:=event.Description
	eventDate:=event.EventDate
	c.Logger().Errorf("info: %v", event)

	tx, err := db.Beginx()
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO `events`"+
		"	(`event_id`, `account_id`, `title`, `description`, `event_date`) VALUES (default,1, ?, ?, ?)",
		title, description , eventDate)
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)

		if ok && mysqlErr.Number == uint16(mysqlErrNumDuplicateEntry) {
			return c.String(http.StatusConflict, "duplicated: event")
		}

		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var id int
	err = tx.Get(&id,"SELECT LAST_INSERT_ID()")
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	err = tx.Commit()
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	ret := map[string]int{"EventID": id}
	return c.JSON(http.StatusCreated, ret)
}
// DELETE api/events/{event_id}
// イベントの削除（※まずは単純削除：ToDo タグなどのBindの掃除）
func deleteEvent(c echo.Context) error {
	eventID := c.Param("event_id")

	tx, err := db.Beginx()
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer tx.Rollback()
	_, err = tx.Exec( "DELETE FROM `events` WHERE `event_id` = ?",eventID)
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	err = tx.Commit()
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
    return c.NoContent(http.StatusNoContent)
}

func getPersonList(c echo.Context) error {
	limit := c.QueryParam("limit")
	if limit != "" {
		limit = fmt.Sprintf(" limit "+limit)
	}
	offset := c.QueryParam("offset")
	if offset != "" {
		offset = fmt.Sprintf(" offset "+offset)
	}
	rows, err := db.Queryx("select * from persons" + limit + offset)
	if err != nil {
		c.Logger().Errorf("failed to query: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	var person Person
	var personlist Personlist
	for rows.Next() {
		err := rows.StructScan(&person) //sqlのrows.Scanの代わりにsqlxのrows.StructScanを使う
		if err != nil {
			c.Logger().Errorf("failed to purse query responce: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		personlist = append(personlist, person)
	}
	return c.JSON(http.StatusOK, personlist)
}

func uploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "file missing")
	}

	image_name := file.Filename
	mime_type := file.Header.Get("Content-Type")
	size := file.Size
	if size > 1000000 {
		return c.String(http.StatusBadRequest, "file exceed 1MByte")
	}
	src, err := file.Open()

	tx, err := db.Beginx()
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO `images` (`image_id`, `image_name`, `mime_type`) VALUES (default, ?, ?)", image_name, mime_type)
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)

		if ok && mysqlErr.Number == uint16(mysqlErrNumDuplicateEntry) {
			return c.String(http.StatusConflict, "duplicated: image")
		}

		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var image_id int
	err = tx.Get(&image_id,"SELECT LAST_INSERT_ID()")
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	fmt.Println(image_id)

    // サーバー側に保存するために空ファイルを作成
	dst_file_name := imagesPath + fmt.Sprintf("/%d.png", image_id)
    dst, err := os.Create(dst_file_name);
    if err != nil {
		c.Logger().Errorf("cant open file: %v", dst_file_name)
		return c.NoContent(http.StatusInternalServerError)
    }
    defer dst.Close();
	if _, err = io.Copy(dst, src); err != nil {
		c.Logger().Errorf("cant save file: %v", dst_file_name)
		return c.NoContent(http.StatusInternalServerError)
	}

	err = tx.Commit()
	if err != nil {
		c.Logger().Errorf("db error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.String(http.StatusOK, "uploading!")
}

func getIndex(c echo.Context) error {
	return c.File(frontendContentsPath + "/index.html")
}
