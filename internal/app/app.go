package app

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/iamsorryprincess/url-shortener/internal/config"
	"github.com/iamsorryprincess/url-shortener/internal/server"
	"github.com/iamsorryprincess/url-shortener/internal/service"
	"github.com/iamsorryprincess/url-shortener/internal/storage"
	"github.com/iamsorryprincess/url-shortener/internal/worker"
	"github.com/iamsorryprincess/url-shortener/pkg/hash"
	_ "github.com/jackc/pgx/stdlib"
)

func Run() {
	configuration, confErr := config.ParseConfiguration()

	if confErr != nil {
		log.Fatal(confErr)
		return
	}

	var db *sql.DB
	var urlService *service.URLService
	var batchWorker *worker.Worker

	if configuration.StoragePath != "" && configuration.DBConnectionString != "" {
		var err error
		db, err = initDB(configuration.DBConnectionString)

		if err != nil {
			log.Fatal(err)
			return
		}

		defer db.Close()
		postgresqlStorage, err := storage.NewPostgresqlStorage(db)

		if err != nil {
			log.Fatal(err)
			return
		}

		urlService = service.NewURLService(postgresqlStorage, configuration.BaseURL)
		batchWorker = worker.NewWorker(postgresqlStorage)
	} else if configuration.DBConnectionString != "" {
		var err error
		db, err = initDB(configuration.DBConnectionString)

		if err != nil {
			log.Fatal(err)
			return
		}

		defer db.Close()
		postgresqlStorage, err := storage.NewPostgresqlStorage(db)

		if err != nil {
			log.Fatal(err)
			return
		}

		urlService = service.NewURLService(postgresqlStorage, configuration.BaseURL)
		batchWorker = worker.NewWorker(postgresqlStorage)
	} else if configuration.StoragePath != "" {
		fileStorage, file, err := storage.NewFileStorage(configuration.StoragePath)

		if err != nil {
			log.Fatal(err)
			return
		}

		defer file.Close()
		urlService = service.NewURLService(fileStorage, configuration.BaseURL)
		batchWorker = worker.NewWorker(fileStorage)
	} else {
		inMemoryStorage := storage.NewInMemoryStorage()
		urlService = service.NewURLService(inMemoryStorage, configuration.BaseURL)
		batchWorker = worker.NewWorker(inMemoryStorage)
	}

	keyManager, err := hash.NewGcmKeyManager()

	if err != nil {
		log.Fatal(err)
		return
	}

	batchWorker.Start(configuration.WorkersCount, configuration.WorkerPoolSize)
	defer batchWorker.Stop()
	httpServer := server.NewServer(configuration, urlService, keyManager, db, batchWorker)
	log.Fatal(httpServer.Run())
}

func initDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
