package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/howeyc/gopass"
	"github.com/tealeg/xlsx"
)

type config struct {
	// MySQL config
	User     string
	Host     string
	Port     string
	Password string
	Database string
	Charset  string

	// other config
	Query string // select query
	File  string // storage file abs path
}

var cfg config

func parseFlag() error {
	mysqlHost := flag.String("h", "localhost", "mysql host")
	mysqlUser := flag.String("u", "", "mysql user name")
	mysqlPassword := flag.String("p", "", "mysql password")
	mysqlDatabase := flag.String("d", "", "mysql database name")
	mysqlPort := flag.String("P", "3306", "mysql port")
	mysqlCharset := flag.String("c", "utf8mb4", "mysql default charset")
	mysqlQuery := flag.String("q", "", "select query")
	fileName := flag.String("f", "", "xlsx file name")
	flag.Parse()

	if *mysqlPassword == "" {
		fmt.Print("Password:")
		password, err := gopass.GetPasswd()
		if err != nil {
			return err
		}
		*mysqlPassword = strings.TrimSpace(string(password))
	}

	if *mysqlQuery == "" {
		// allow line separator, sql end with ';'
		fmt.Println("Query (end with '; + <Enter>'):")
		reader := bufio.NewReader(os.Stdin)
		for {
			tmpSql, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			*mysqlQuery = *mysqlQuery + tmpSql
			tmpSql = strings.TrimSpace(tmpSql)
			if tmpSql[len(tmpSql)-1] == ';' {
				break
			}
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if *fileName == "" {
		*fileName = pwd + "/temp.xlsx"
	}

	if !strings.HasPrefix(*fileName, "/") {
		*fileName = pwd + "/" + *fileName
	}

	cfg = config{
		User:     *mysqlUser,
		Host:     *mysqlHost,
		Port:     *mysqlPort,
		Password: *mysqlPassword,
		Database: *mysqlDatabase,
		Charset:  *mysqlCharset,

		Query: *mysqlQuery,
		File:  *fileName,
	}

	return nil
}

func getData() (*sql.Rows, error) {
	var res *sql.Rows

	// init database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return res, err
	}
	defer db.Close()

	res, err = db.Query(cfg.Query)
	return res, err
}

func save2xlsx(filePath string, rows *sql.Rows) error {
	file := xlsx.NewFile()
	// create new sheet
	sheet, err := file.AddSheet("result")
	if err != nil {
		return err
	}

	// set table header with column name
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	sheetHeader := sheet.AddRow()
	for _, name := range columns {
		cell := sheetHeader.AddCell()
		cell.Value = name
	}

	// set every table rows
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		rows.Scan(scanArgs...)
		sheetRow := sheet.AddRow()
		for _, v := range values {
			cell := sheetRow.AddCell()
			cell.Value = string(v)
		}
	}

	// save to file
	return file.Save(filePath)
}

func main() {
	// parse config
	if err := parseFlag(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// execute sql and get all result rows
	rows, err := getData()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// save result into xlsx files
	err = save2xlsx(cfg.File, rows)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("save data into file: '%s'\n", cfg.File)
}
