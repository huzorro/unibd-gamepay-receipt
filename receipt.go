package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func moReview(r *http.Request, w http.ResponseWriter, db *sql.DB, log *log.Logger) (int, string) {
	rowsN := url.QueryEscape(r.URL.Query().Get("rows"))
	stmtOut, err := db.Prepare("SELECT spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time FROM mo_receipt LIMIT ?")
	defer stmtOut.Close()
	rows, err := stmtOut.Query(rowsN)
	if err != nil {
		panic(err.Error())
	}
	var rowArray []string
	for rows.Next() {
		var spid string
		var srctermid string
		var linkid string
		var citycode string
		var cmd string
		var desttermid string
		var fee string
		var serviceid string
		var time string
		err = rows.Scan(&spid, &srctermid, &linkid, &citycode, &cmd, &desttermid, &fee, &serviceid, &time)
		if err != nil {
			panic(err.Error)
		}
		log.Printf("spid=%s, srctermid=%s, linkid=%s, citycode=%s, cmd=%s, desttermid=%s, fee=%s, serviceid=%s, time=%s", spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time)
		rowStr := fmt.Sprintf("spid=%s, srctermid=%s, linkid=%s, citycode=%s, cmd=%s, desttermid=%s, fee=%s, serviceid=%s, time=%s", spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time)
		rowArray = append(rowArray, rowStr)
	}
	return http.StatusOK, strings.Join(rowArray, "<BR/>")
}
func moReceipt(r *http.Request, w http.ResponseWriter, db *sql.DB, log *log.Logger) (int, string) {
	spid, _ := url.QueryUnescape(r.URL.Query().Get("Spid"))
	srctermid, _ := url.QueryUnescape(r.URL.Query().Get("Src"))
	linkid, _ := url.QueryUnescape(r.URL.Query().Get("Linkid"))
	citycode, _ := url.QueryUnescape(r.URL.Query().Get("CityCode"))
	cmd, _ := url.QueryUnescape(r.URL.Query().Get("Cmd"))
	desttermid, _ := url.QueryUnescape(r.URL.Query().Get("Dest"))
	fee, _ := url.QueryUnescape(r.URL.Query().Get("Fee"))
	serviceid, _ := url.QueryUnescape(r.URL.Query().Get("Svcid"))
	time, _ := url.QueryUnescape(r.URL.Query().Get("Time"))
	//?Spid=901077&Src=134111&Linkid=112212121221&CityCode=0010&Cmd=wz*345&Dest=10669501&Fee=100&Svcid=xxyz&Time=20140225165755
	// CREATE TABLE mo_receipt (
	// 	id int(11) NOT NULL AUTO_INCREMENT,
	// 	spid varchar(6) NOT NULL DEFAULT '',
	// 	srctermid varchar(32) NOT NULL DEFAULT '',
	// 	citycode varchar(4) NOT NULL DEFAULT '',
	// 	desttermid varchar(30) NOT NULL DEFAULT '',
	// 	linkid varchar(30) NOT NULL DEFAULT '',
	// 	cmd varchar(150) NOT NULL DEFAULT '',
	// 	fee varchar(10) NOT NULL DEFAULT '',
	// 	serviceid varchar(18) NOT NULL DEFAULT '',
	// 	status varchar(50) NOT NULL DEFAULT '',
	// 	time varchar (20) NOT NULL DEFAULT '',
	// 	logtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	// 	PRIMARY KEY (id),
	// 	UNIQUE KEY `src_linkid` (`srctermid`,`linkid`)
	// )ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8
	stmtIn, err := db.Prepare("INSERT INTO mo_receipt (spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIn.Close()
	// _, err = stmtIn.Exec(spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time)
	res, err := stmtIn.Exec(spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time)

	if err != nil {
		panic(err.Error())
	}
	rowId, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	log.Printf("<%d> INSERT INTO mo_receipt (spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", rowId, spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time)
	// log.Printf("<> INSERT INTO mo_receipt (spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", spid, srctermid, linkid, citycode, cmd, desttermid, fee, serviceid, time)

	return http.StatusOK, "resultCode=0"
}
func mrReceipt(r *http.Request, w http.ResponseWriter, db *sql.DB, log *log.Logger) (int, string) {
	spid, _ := url.QueryUnescape(r.URL.Query().Get("Spid"))
	srctermid, _ := url.QueryUnescape(r.URL.Query().Get("Src"))
	linkid, _ := url.QueryUnescape(r.URL.Query().Get("Linkid"))
	status, _ := url.QueryUnescape(r.URL.Query().Get("Status"))
	cmd, _ := url.QueryUnescape(r.URL.Query().Get("Cmd"))
	//?Spid=901077&Src=134111&Linkid=112212121221&Status=DELIVRD&Cmd=wz*345
	// 	CREATE TABLE mr_receipt (
	// 	id int(11) NOT NULL AUTO_INCREMENT,
	// 	spid varchar(6) NOT NULL DEFAULT '',
	// 	srctermid varchar(32) NOT NULL DEFAULT '',
	// 	linkid varchar(30) NOT NULL DEFAULT '',
	// 	cmd varchar(150) NOT NULL DEFAULT '',
	// 	status varchar (30) NOT NULL DEFAULT '',
	// 	logtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	// 	PRIMARY KEY (id),
	// 	UNIQUE KEY `src_linkid` (`srctermid`,`linkid`)
	// )ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8
	stmtIn, err := db.Prepare("INSERT INTO mr_receipt (spid, srctermid, linkid, status, cmd) VALUES(?, ?, ?, ?, ?) ")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIn.Close()
	res, err := stmtIn.Exec(spid, srctermid, linkid, status, cmd)
	// _, err = stmtIn.Exec(spid, srctermid, linkid, status, cmd)

	if err != nil {
		panic(err.Error())
	}
	rowId, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	log.Printf("<%d> INSERT INTO mr_receipt (spid, srctermid, linkid, status, cmd) VALUES('%s', '%s', '%s', '%s', '%s')", rowId, spid, srctermid, linkid, status, cmd)
	// log.Printf("<> INSERT INTO mr_receipt (spid, srctermid, linkid, status, cmd) VALUES('%s', '%s', '%s', '%s', '%s')", spid, srctermid, linkid, status, cmd)
	return http.StatusOK, "resultCode=0"
}
func main() {
	mtn := martini.Classic()
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/receipt?charset=utf8")
	db.SetMaxOpenConns(10)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	mtn.Map(db)
	logger := log.New(os.Stdout, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)
	mtn.Map(logger)

	mtn.Get("/moReceiver", moReceipt)
	mtn.Get("/mrReceiver", mrReceipt)
	mtn.Get("/moReview", moReview)
	// mtn.Get("/mrReview", mrReivew)
	http.ListenAndServe(":10086", mtn)
	// mtn.Run()
}

// 	CREATE TABLE tmp_mo_receipt (
// 		id int(11) NOT NULL AUTO_INCREMENT,
// 		spid varchar(6) NOT NULL DEFAULT '',
// 		srctermid varchar(32) NOT NULL DEFAULT '',
// 		citycode varchar(4) NOT NULL DEFAULT '',
// 		desttermid varchar(30) NOT NULL DEFAULT '',
// 		linkid varchar(30) NOT NULL DEFAULT '',
// 		cmd varchar(150) NOT NULL DEFAULT '',
// 		fee varchar(10) NOT NULL DEFAULT '',
// 		serviceid varchar(18) NOT NULL DEFAULT '',
// 		status varchar(50) NOT NULL DEFAULT '',
// 		time varchar (20) NOT NULL DEFAULT '',
// 		logtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 		PRIMARY KEY (id),
// 		UNIQUE KEY `src_linkid` (`srctermid`,`linkid`)
// 	)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8

// DELIMITER //
// CREATE TRIGGER t_afterupdate_on_mo_receipt
// AFTER update ON mo_receipt
// FOR EACH ROW
// BEGIN
//      insert into tmp_mo_receipt(spid, srctermid, citycode, desttermid, linkid, cmd, fee, serviceid, status, time, logtime) values(NEW.spid, NEW.srctermid, NEW.citycode, NEW.desttermid, NEW.linkid, NEW.cmd, NEW.fee, NEW.serviceid, NEW.status, NEW.time, NEW.logtime);
// END//

// UPDATE mo_receipt a, mr_receipt b SET a.status = b.status WHERE a.srctermid = b.srctermid AND a.linkid = b.linkid

// DELIMITER //
// CREATE TRIGGER t_afterinsert_on_mr_receipt
// AFTER insert ON mr_receipt
// FOR EACH ROW
// BEGIN
// 	 update mo_receipt SET status = NEW.status WHERE srctermid = NEW.srctermid AND linkid = NEW.linkid;
// END//

// DELIMITER ;

// //接口文档

// ?feetermid=计费手机号伪码&status=状态报告&linkid=linkid&desttermid=用户手机号伪码&cmdword=业务代码&sendtime=计费时间&msgid=消息id
