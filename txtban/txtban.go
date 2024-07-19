package txtban

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/thehxdev/txtban/config"
	"github.com/thehxdev/txtban/models"
)

const VERSION string = "1.2.3"

type Txtban struct {
	Server    *http.Server
	DB        *models.DB
	ErrLogger *log.Logger
	InfLogger *log.Logger
	Config    *config.TbConfig
}

func Init(configPath string) *Txtban {
	t := &Txtban{
		ErrLogger: log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
		InfLogger: log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		DB:        &models.DB{},
		Config:    &config.TbConfig{},
	}

	config.SetupViper(t.Config, configPath)

	t.Server = &http.Server{
		Addr:    net.JoinHostPort(viper.GetString("server.address"), strconv.Itoa(viper.GetInt("server.port"))),
		Handler: t.configureRoutes(),
	}

	t.setupDB(viper.GetString("database.path"))

	return t
}

func (t *Txtban) Run() error {
	t.InfLogger.Println("starting server...")
	return t.Server.ListenAndServe()
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

func (t *Txtban) configureRoutes() http.Handler {
	t.InfLogger.Println("configuring routes...")
	r := chi.NewRouter()

	// Root
	r.Get("/", rootHandler)

	// User related routes
	r.Post("/useradd", t.useraddHandler)
	r.Get("/whoami", t.whoamiHandler)
	r.Delete("/userdel", t.userdelHandler)
	r.Put("/passwd", t.passwdHandler)

	// Txt related routes
	r.Post("/tee", t.teeHandler)
	r.Get("/ls", t.lsHandler)
	r.Put("/chtxt", t.chtxtHandler)
	r.Delete("/rm", t.rmHandler)
	r.Get("/t/{txtid}", t.readHandler)
	r.Put("/mv", t.mvHandler)
	r.Put("/rename", t.renameHandler)

	return r
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("txtban service v%s running!", VERSION)
	w.Write([]byte(msg))
}
