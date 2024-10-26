package core

import (
	"database/sql"
	"errors"
	"github.com/leijeng/huo-core/common/consts"
	"github.com/leijeng/huo-core/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"path"
	"time"
)

func dbInit() {

	if Cfg.DBCfg.DSN != "" {
		logMode := config.GetLogMode(Cfg.DBCfg.LogMode)
		initDb(Cfg.DBCfg.Driver, Cfg.DBCfg.DSN, Cfg.DBCfg.Prefix, consts.DB_DEF, logMode, Cfg.DBCfg.SlowThreshold,
			Cfg.DBCfg.MaxIdleConns, Cfg.DBCfg.MaxOpenConns, Cfg.DBCfg.MaxLifetime, Cfg.DBCfg.Singular, Cfg.Logger.Color(), Cfg.DBCfg.IgnoreNotFound)
	}
	for key, dbc := range Cfg.DBCfg.DBS {
		if !dbc.Disable {
			var logMode logger.LogLevel
			if dbc.LogMode != "" {
				logMode = config.GetLogMode(dbc.LogMode)
			} else {
				logMode = config.GetLogMode(Cfg.DBCfg.LogMode)
			}
			prefix := dbc.Prefix
			if prefix == "" && Cfg.DBCfg.Prefix != "" {
				prefix = Cfg.DBCfg.Prefix
			}
			slow := dbc.SlowThreshold
			if slow < 1 && Cfg.DBCfg.SlowThreshold > 0 {
				slow = Cfg.DBCfg.SlowThreshold
			}
			singular := Cfg.DBCfg.Singular
			maxIdle := dbc.MaxIdleConns
			if maxIdle < 1 {
				maxIdle = Cfg.DBCfg.GetMaxIdleConns()
			}

			maxOpen := dbc.MaxOpenConns
			if maxOpen < 1 {
				maxOpen = Cfg.DBCfg.GetMaxOpenConns()
			}

			maxLifetime := dbc.MaxLifetime
			if maxLifetime < 1 {
				maxLifetime = Cfg.DBCfg.GetMaxLifetime()
			}
			driver := dbc.Driver
			if driver == "" && Cfg.DBCfg.Driver != "" {
				driver = Cfg.DBCfg.Driver
			}
			ignoreNotFound := dbc.IgnoreNotFound
			if !ignoreNotFound && Cfg.DBCfg.IgnoreNotFound {
				ignoreNotFound = Cfg.DBCfg.IgnoreNotFound
			}
			initDb(driver, dbc.DSN, prefix, key, logMode, slow, maxIdle, maxOpen, maxLifetime, singular, Cfg.Logger.Color(), ignoreNotFound)
		}
	}

}

func initDb(driver, dns, prefix, key string, logMode logger.LogLevel, slow, maxIdle, maxOpen, maxLifetime int, singular, color, ignoreNotFound bool) {
	var db *gorm.DB
	var err error
	switch driver {
	case Mysql.String():
		db, err = gorm.Open(mysql.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound))
	case Pgsql.String():
		db, err = gorm.Open(postgres.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound))
	case Sqlite.String():
		db, err = gorm.Open(sqlite.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound))
	case Mssql.String():
		db, err = gorm.Open(sqlserver.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound))
	// case "oracle":
	// 	db, err = gorm.Open(oracle.Open(dbc.DSN), &gorm.Config{})
	// case "clickhouse":
	// 	db, err = gorm.Open(clickhouse.Open(dbc.DSN), &gorm.Config{})
	default:
		err = errors.New("db err")
	}
	if err != nil {
		log.Printf("db:%s dns:%s err:%s", key, dns, err.Error())
		panic(err)
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		log.Printf("db:%s dns:%s err:%s", key, dns, err.Error())
		panic(err)
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(maxLifetime))
	SetDb(key, db)
}

func GetGromLogCfg(logMode logger.LogLevel, prefix string, slowThreshold int, singular, color, ignoreNotFound bool) *gorm.Config {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		//DisableForeignKeyConstraintWhenMigrating: true,
	}

	filePath := path.Join(Cfg.Logger.Director, "%Y-%m-%d", "sql.log")
	w, _ := GetWriter(filePath)
	slow := time.Duration(slowThreshold) * time.Millisecond
	_default := logger.New(log.New(w, "", log.LstdFlags), logger.Config{
		SlowThreshold:             slow,
		Colorful:                  color,
		IgnoreRecordNotFoundError: ignoreNotFound,
	})

	config.Logger = _default.LogMode(logMode)

	return config
}

// var (
// 	gplusInit = true
// )

func SetDb(key string, db *gorm.DB) {
	lock.Lock()
	defer lock.Unlock()
	dbs[key] = db
	// //集成gplus
	// if gplusInit {
	// 	gplusInit = false
	// 	gplus.Init(db)
	// }
}

// GetDb 获取所有map里的db数据
func Dbs() map[string]*gorm.DB {
	// lock.RLock()
	// defer lock.RUnlock()
	return dbs
}

func Db(name string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	if db, ok := dbs[name]; !ok || db == nil {
		log.Printf("db:%s not init", name)
		panic("db not init")
		// Log.Error("DB not found", zap.Error(errors.New(name+" DB not found")))
		// return nil
	} else {
		return db
	}
}

// 获取默认的（master）db
func DB() *gorm.DB {
	return Db(consts.DB_DEF)
}