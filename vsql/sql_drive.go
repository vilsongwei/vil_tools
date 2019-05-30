package vsql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

type SqlDriver interface {
	Connector(user, passwd, addr, dbname, version string, maxIdle, maxOpen int) error
	BreakOff()
}

type SqlDrive struct {
	SqlCnf *MySqlCnf
	Driver SqlDriver
	Lock   sync.RWMutex
}

func (sel *SqlDrive) Open() error {
	// 连接默认数据库
	if err := sel.Driver.Connector(sel.SqlCnf.UserName, sel.SqlCnf.Password, sel.SqlCnf.Address, sel.SqlCnf.Default, sel.SqlCnf.Version, sel.SqlCnf.MaxIdleConns, sel.SqlCnf.MaxOpenConns); err != nil {
		return err
	}
	// 连接其他指定的数据库
	for _, dbv := range sel.SqlCnf.Databases {
		if err := sel.Driver.Connector(sel.SqlCnf.UserName, sel.SqlCnf.Password, sel.SqlCnf.Address, dbv, sel.SqlCnf.Version, sel.SqlCnf.MaxIdleConns, sel.SqlCnf.MaxOpenConns); err != nil {
			return err
		}
	}
	return nil
}

func (sel *SqlDrive) SqlConnStr(version, user, passwd, addr, dbname string) (cs string) {
	switch version {
	case "5":
		cs = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, addr, dbname)
	case "8":
		cs = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true", user, passwd, addr, dbname)
	default:
		cs = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, addr, dbname)
	}
	return
}
