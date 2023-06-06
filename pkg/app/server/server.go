package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	//"syscall"

	"nostr-relay/pkg/app/session"
	"nostr-relay/pkg/config"
	"nostr-relay/pkg/db"
	"nostr-relay/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	gormlogger "gorm.io/gorm/logger"
)

type Server struct {
	Config *config.Config
	srv    *http.Server
}

func initLogger() {
	//log輸出為json格式
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//輸出設定為標準輸出(預設為stderr)
	logrus.SetOutput(os.Stdout)
	//設定要輸出的log等級
	logrus.SetLevel(logrus.DebugLevel)
}

func New(cfg *config.Config) *Server {

	return &Server{
		Config: cfg,
	}

}

func (t *Server) Serve() {

	t.init()
	addr := ":" + t.Config.Server.Port

	log.Printf("======= Server start to listen (%s) and serve =======\n", addr)
	r := Router(t)
	// r.Run(addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	t.srv = srv

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Prepare to Shutdown Server ...")
	t.shutdown()

	log.Printf("======= Server Exit =======\n")
	//CloseLogger()
}

func (t *Server) init() {
	initLogger()
	t.initDB()
	fmt.Println("Automigrate", t.Config.DB.Automigrate)
	if t.Config.DB.Automigrate {
		t.migration() //remove for production
	}

}

func (t *Server) migration() { //remove for production
	fmt.Println("Run Migration --> Start")
	db.GetMainDB().AutoMigrate(&models.RelayEvent{})
	fmt.Println("Run Migration --> Done")

}

func (t *Server) initDB() {

	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormlogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,           // Disable color
		},
	)

	l := db.Logger{
		Logger: newLogger,
	}

	dbConn, err := db.GormOpen(&t.Config.DB, &l)
	if err != nil {
		logrus.Fatal(err)
	}
	db.SetMainDB(dbConn)

	sqlDB, err := dbConn.DB()
	if err != nil {
		logrus.Fatal(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		logrus.Fatal(err)
	}
}

func (t *Server) Shutdown(c *gin.Context) {
	t.shutdown()
}

func (t *Server) shutdown() error {
	if t.srv == nil {
		log.Fatal("Shutdown: server srv == nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := t.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server Shutdown")

	log.Println("ForEachSession Close")

	session.ForEachSession(func(s session.SessionF) {
		s.Close()
	})
	session.WaitGroup.Wait()
	log.Println("ForEachSession Close Done")

	return nil
}
