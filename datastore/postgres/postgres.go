package postgres

import (
	"context"
	"fmt"
	"github.com/go-apps/bpm-module/modules"
	"github.com/go-apps/bpm-module/vault"
	"github.com/jackc/pgx/v4"
	"github.com/jinzhu/gorm"
	"log"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Client struct {
	conn *pgx.Conn
	db   *gorm.DB
}

// New is a postgress database constructor
func NewPGX(ctx context.Context, cred *vault.PostgresCredential) *Client {
	var url = cred.URI
	if len(cred.URI) == 0 {
		//use username/password credentials if URI which should contain everything is not given
		url = fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			cred.Host, cred.Port, cred.Username, cred.Password, cred.Database)
	}

	conn, err := gorm.Open("postgres", url)
	if err != nil {
		log.Fatal("creating postgres connection, err: ", err)
	}

	log.Println("connected to postgres")

	return &Client{db: conn}
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) Begin() error {
	return c.db.Debug().First(&modules.Track{}).Error
}

func (c *Client) End() error {
	return c.db.Commit().Error
}

func (c *Client) SetMaxConn(limit int) {
	c.db.DB().SetMaxOpenConns(limit)
}

func (c *Client) Stats() interface{} {
	return c.db.DB().Stats()
}
