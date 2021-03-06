package Config

import (
	"errors"
	"fmt"
	"github.com/jasonlvhit/gocron"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"github.com/jinzhu/gorm"
)

type Service interface {
	Start()
}

type App struct {
	Config     *Configuration
	Db         *gorm.DB
	Scheduler  *gocron.Scheduler
	Services   []Service
}

func NewApp() *App {
	config := SetupConfiguration()
	app := App{

		Config:    config,
		Scheduler: gocron.NewScheduler(),
	}
	if config.DbConfig.DbDriver != "" {
		var connectionString string = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", config.DbConfig.Host, config.DbConfig.Port, config.DbConfig.UserName, config.DbConfig.TableName, config.DbConfig.Password)
		log.Println("Using Connection string", connectionString)
		db, err := gorm.Open(config.DbConfig.DbDriver, connectionString)
		if err != nil {
			log.Println("Db open error:", err.Error())
		}
		err = db.DB().Ping()
		if err != nil {
			log.Println("Db open error:", err.Error())
		}

		app.Db = db
	}

	return &app
}

func (app *App) AddService(svc Service) {
	app.Services = append(app.Services, svc)
}

func (app *App) Run() {

	for _, service := range app.Services {
		service.Start()
	}

	go func() {
		log.Println("Starting Scheduler")
		<-app.Scheduler.Start()
	}()
	if app.Db != nil {
		defer app.Db.Close()
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	signal.Notify(sigs, os.Kill)
	log.Println("pushed global error handler")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
			debug.PrintStack()
			// find out exactly what the error was and set err
			var err error
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			if err != nil {
				// sendMeMail(err)
				log.Println("Panic occured", err.Error())
			}
		}
	}()
	go func() {
		log.Println("Waiting for SIGTERM to exit cleanly")
		sig := <-sigs
		log.Println("recived signal", sig)
		done <- true
	}()

	<-done

}
