package txtban

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/thehxdev/txtban/config"
	"github.com/thehxdev/txtban/models"
)

const VERSION string = "1.2.2"

type Txtban struct {
	App       *fiber.App
	DB        *models.DB
	ErrLogger *log.Logger
	InfLogger *log.Logger
	Config    *config.TbConfig
}

func Init(configPath string) *Txtban {
	t := &Txtban{
		App:       fiber.New(),
		ErrLogger: log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
		InfLogger: log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		DB:        &models.DB{},
		Config:    &config.TbConfig{},
	}

	config.SetupViper(t.Config, configPath)

	t.ConfigureRoutes()
	t.setupDB(viper.GetString("database.path"))

	return t
}

func (t *Txtban) Run() error {
	t.InfLogger.Println("starting server...")
	listenAddr := fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt("server.port"))
	return t.App.Listen(listenAddr)
}

func (t *Txtban) setupDB(path string) {
	var dbIsNew bool
	if _, err := os.Stat(path); err != nil {
		t.InfLogger.Println("must create a new database file")
		dbIsNew = true
	}

	t.InfLogger.Println("Setting up sqlite3 database...")
	err := t.DB.SetupSqliteDB(path)
	if err != nil {
		t.ErrLogger.Fatal(err)
	}

	if dbIsNew {
		t.InfLogger.Println("creating tables...")
		t.DB.MigrateDB()
	}
}

func (t *Txtban) CloseDB() {
	t.DB.Read.Close()
	t.DB.Write.Close()
}

func (t *Txtban) ConfigureRoutes() {
	t.InfLogger.Println("configuring routes...")
	app := t.App

	// Root
	app.Get("/", rootHandler)

	// User related routes
	app.Post("/useradd", t.useraddHandler)
	app.Get("/whoami", t.whoamiHandler)
	app.Delete("/userdel", t.userdelHandler)
	app.Put("/passwd", t.passwdHandler)

	// Txt related routes
	app.Post("/tee", t.teeHandler)
	app.Get("/ls", t.lsHandler)
	app.Put("/chtxt", t.chtxtHandler)
	app.Delete("/rm", t.rmHandler)
	app.Get("/t/:txtid", t.readHandler)
	app.Put("/mv", t.mvHandler)
	app.Put("/rename", t.renameHandler)
}

func rootHandler(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("txtban service v%s running!", VERSION))
}
