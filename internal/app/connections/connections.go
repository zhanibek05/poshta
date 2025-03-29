package connections

import (
	"poshta/internal/app/config"
	"log"
	_"github.com/go-sql-driver/mysql" 
	"github.com/jmoiron/sqlx"
)

type Connections struct {
	DB *sqlx.DB
}

func NewConnections(cfg *config.Config) (*Connections, error) {
	db, err := sqlx.Open("mysql", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to MySQL successfully")
	return &Connections{DB: db}, nil
}

func (c *Connections) Close() {
	if c.DB != nil {
		c.DB.Close()
		log.Println("Database connection closed")
	}
}
