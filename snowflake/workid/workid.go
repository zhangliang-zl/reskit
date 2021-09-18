package workid

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"net"
)

var (
	createSql = "CREATE TABLE  if not exists `workid` (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`ip` varchar(64) NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `idx_ip` (`ip`)" +
		") ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4"
)

type Factory struct {
	mysqlDsn string
}

func (f *Factory) Get() (int64, error) {

	if f.mysqlDsn == "" {
		return 0, errors.New("mysql dsn is nil")
	}

	db, err := sql.Open("mysql", f.mysqlDsn)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// create table if not exist
	if _, err := db.Exec(createSql); err != nil {
		return 0, err
	}

	ip, err := f.privateIPv4()
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	ipStr := ip.To4().String()
	// select or insert
	var id int64
	err = tx.QueryRow("select id from `workid` where `ip`=?  for update ", ipStr).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return 0, err
	}

	if err == sql.ErrNoRows {
		result, err := tx.Exec("insert into workid(`ip`) values(?)", ipStr)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		id, err = result.LastInsertId()
		if err != nil {
			return 0, tx.Rollback()
		}
	}

	return id, tx.Commit()
}

func (f *Factory) privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipNet, ok := a.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}

		ip := ipNet.IP.To4()
		if f.isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func (f *Factory) isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func NewFactory(dsn string) *Factory {
	return &Factory{
		mysqlDsn: dsn,
	}
}
