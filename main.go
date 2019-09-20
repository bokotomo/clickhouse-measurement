package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	. "sample/util"

	"github.com/kshvakov/clickhouse"
)

var (
	connectCH    *sql.DB
	connectMysql *sql.DB
)

func driverCH() error {
	var err error

	connectCH, err = sql.Open("clickhouse", "tcp://server:9000?debug=true")
	if err != nil {
		return err
	}

	// デバック用
	if err := connectCH.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return err
	}

	return nil
}

func driverMysql() error {
	var err error

	connectMysql, err = sql.Open("mysql", "dev:dev@tcp(db:3388)/dev?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return err
	}

	return nil
}

func sample() error {
	// create
	_, err := connectCH.Exec(`
		CREATE TABLE IF NOT EXISTS example (
			country_code FixedString(2),
			os_id        UInt8,
			browser_id   UInt8,
			categories   Array(Int16),
			action_day   Date,
			action_time  DateTime
		) engine=Memory
	`)
	if err != nil {
		return err
	}

	// インサートトランザクション
	var (
		tx, _   = connectCH.Begin()
		stmt, _ = tx.Prepare("INSERT INTO example (country_code, os_id, browser_id, categories, action_day, action_time) VALUES (?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		if _, err := stmt.Exec(
			"RU",
			10+i,
			100+i,
			clickhouse.Array([]int16{1, 2, 3}),
			time.Now(),
			time.Now(),
		); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	// select
	rows, err := connectCH.Query("SELECT country_code, os_id, browser_id, categories, action_day, action_time FROM example")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			country               string
			os, browser           uint8
			categories            []int16
			actionDay, actionTime time.Time
		)
		if err := rows.Scan(&country, &os, &browser, &categories, &actionDay, &actionTime); err != nil {
			return err
		}
		log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_day: %s, action_time: %s", country, os, browser, categories, actionDay, actionTime)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// drop
	if _, err := connectCH.Exec("DROP TABLE example"); err != nil {
		return err
	}

	return nil
}

func insert(sizeIn int, j int) error {
	// インサートトランザクション
	var (
		tx, _   = connectCH.Begin()
		stmt, _ = tx.Prepare("INSERT INTO example2 (a,b,c,d,e,f,g,h,i,j) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()
	for i := 0; i < sizeIn; i++ {
		a := j*sizeIn + i
		if _, err := stmt.Exec(
			"RU",
			10+a,
			1000+a,
			0.1112,
			time.Now(),
			time.Now(),
			RandString[rand.Intn(200)],
			RandString[rand.Intn(200)],
			RandString[rand.Intn(200)],
			RandString[rand.Intn(200)],
		); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// サンプルデータインサート
// データ容量＋データ速度の検証をする。
// 1GB＋100,000,000レコードで１０カラムぐらいで検証
func sample2_1() error {
	// create
	_, err := connectCH.Exec(`
		CREATE TABLE IF NOT EXISTS example2 (
			a   FixedString(2),
			b   Int32,
			c   Int32,
			d   Float32,
			e   Date,
			f   DateTime,
			g   String,
			h   String,
			i   String,
			j   String
		) engine=Memory
	`)
	if err != nil {
		return err
	}

	size := 100000000
	sizeOut := 10000
	sizeIn := size / sizeOut
	for j := 0; j < sizeOut; j++ {
		if err := insert(sizeIn, j); err != nil {
			return err
		}
		fmt.Println("done", j)
	}

	return nil
}

// サンプルデータインサート
// データ容量＋データ速度の検証をする。
// 1GB＋100,000,000レコードで１０カラムぐらいで検証
func sample2_2() error {

	return nil
}

func showRandString() {
	for i := 0; i < 200; i++ {
		fmt.Println("\"" + RandString6(20) + "\",")
	}
}

func main() {
	if err := driverCH(); err != nil {
		log.Fatal(err)
	}
	if err := driverMysql(); err != nil {
		log.Fatal(err)
	}
	// if err := sample(); err != nil {
	// 	log.Fatal(err)
	// }
	// rand.Seed(time.Now().UnixNano())
	// var a string
	// for i := 0; i < 1000; i++ {
	// 	for j := 0; j < 100000; j++ {
	// 		a = a + "," + RandString[rand.Intn(200)]
	// 	}
	// 	a = ""
	// 	fmt.Println("done", i)
	// }
	// fmt.Println("OK")
	// // fmt.Println()

	// connectCH.Exec("DROP TABLE example2")
	// if err := sample2_1(); err != nil {
	// 	log.Fatal(err)
	// }

	if err := sample2_2(); err != nil {
		log.Fatal(err)
	}
}
