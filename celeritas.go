package celeritas

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const version = "1.0.0"

// Celeritas is the overall type for the Celeritas package. Members that are exported in this
// type are available to any application that uses it.
type Celeritas struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	config   config
	Routes   *chi.Mux
}

type config struct {
	port     string
	renderer string
}

// New reads the .env file, creates our application config, populates the Celeritas type
// with settings based on .env values, and creates necessary folders and files if they don't exist
func (c *Celeritas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath: rootPath,
		folderNames: []string{
			"handlers", "migration", "views", "data", "public", "tmp", "logs", "middleware",
		},
	}

	err := c.Init(pathConfig)
	if err != nil {
		return err
	}

	err = c.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	infoLog, errorLog := c.startLoggers()
	c.InfoLog = infoLog
	c.ErrorLog = errorLog
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = version
	c.RootPath = rootPath
	c.Routes = c.routes().(*chi.Mux)

	c.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	return nil
}

// Init creates necessary folders for our Celeritas application
func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		err := c.CreateDirIfNotExist(fmt.Sprintf("%s/%s", root, path))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Celeritas) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     c.ErrorLog,
		Handler:      c.routes(),
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	c.InfoLog.Printf("Listening on port %s...", os.Getenv("PORT"))
	c.ErrorLog.Fatal(srv.ListenAndServe())
}

func (c *Celeritas) checkDotEnv(root string) error {
	return c.CreateFileIfNotExist(fmt.Sprintf("%s/.env", root))
}

func (c *Celeritas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog, errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
