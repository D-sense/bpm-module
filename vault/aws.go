package vault

import (
	"context"
	"fmt"
	"github.com/go-apps/bpm-module/modules"
	"github.com/jackc/pgx/v4"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
)


// Config is global object that holds all application level variables.
var Config appConfig

type PostgresCredential struct {
	URI      string `long:"uri" json:"uri"`
	Username string `long:"username" default:"username" json:"username"`
	Password string `long:"password" default:"password" json:"password"`
	Database string `long:"database" default:"database" json:"database"`
	Host     string `long:"host" default:"host" json:"host"`
	Port     string `long:"port" default:"port" json:"port"`
}

type Client struct {
	conn *pgx.Conn
	db   *gorm.DB
}

// NewPGX is a postgres database constructor
func (c *Client) NewPGX(ctx context.Context, cred *PostgresCredential) Client {
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

	return Client{db: conn}
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) Begin() error {
	// return c.db.Begin().Error
	return c.db.Debug().First(&modules.Track{}).Error
}

func (c *Client) End() error {
	return c.db.Commit().Error
	// return c.db.Debug().First(&core.User{}).Error
}

func (c *Client) SetMaxConn(limit int) {
	c.db.DB().SetMaxOpenConns(limit)
}

func (c *Client) Stats() interface{} {
	return c.db.DB().Stats()
}

func GetCredentials() (*PostgresCredential, error) {
	if err := LoadConfig("../../vault"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	return &PostgresCredential{
		URI:      Config.Url,
		Username: Config.Database.AppName,
		Password: Config.Database.Password ,
		Database: Config.Database.Username,
		Host:     Config.Server.Host,
		Port:     Config.Server.Port,
	}, nil
}

type appConfig struct {
	Url string `yaml:"url"`

	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	Database struct {
		AppName     string `yaml:"app_name`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.SetEnvPrefix("gbedu")
	v.AutomaticEnv()

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}
	return v.Unmarshal(&Config)
}