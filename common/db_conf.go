package common

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatasourceConfig struct {
	SourceName string      `json:"sourceName"`
	DBType     string      `json:"type"`
	DBConfig   DBConfigure `json:"config"`
}

type _DatasourceConfig DatasourceConfig

func (c *DatasourceConfig) UnmarshalJSON(b []byte) error {
	d := _DatasourceConfig{}
	switch jsoniter.Get(b, "type").ToString() {
	case "mysql":
		d.DBConfig = &MysqlDBConfig{}
	default:
		return errors.New("unsupported ds type")
	}
	if err := jsoniter.Unmarshal(b, &d); err != nil {
		return err
	}
	*c = (DatasourceConfig)(d)
	return nil
}

type DBConfigure interface {
	ConnectDB() *gorm.DB
}

type MysqlDBConfig struct {
	Address   string `json:"address"`
	Database  string `json:"database"`
	User      string `json:"user"`
	Password  string `json:"password"`
	ExtParams string `json:"extParams"`
}

func (c *MysqlDBConfig) ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", c.User, c.Password, c.Address, c.Database)
	if len(c.ExtParams) > 0 {
		dsn = dsn + "?" + c.ExtParams
	}
	log.Infof(">>>>>>> Connect Mysql dsn: %s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true,
		Logger:      logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Errorf("Connect to Mysql [%s] failed, config:[%v]", c.Database, c)
	}
	log.Infof(">>>>>>> Connect Mysql db [%s] succeed!", c.Database)
	return db
}
