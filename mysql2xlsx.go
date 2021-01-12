package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/howeyc/gopass"
	"github.com/tealeg/xlsx"
	ini "gopkg.in/ini.v1"
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
	mysqlDefaultsExtraFile := flag.String("defaults-extra-file", "", "mysql --defaults-extra-file")

	mysqlQuery := flag.String("q", "", "select query")
	fileName := flag.String("f", "", "xlsx file name")
	flag.Parse()

	if *mysqlDefaultsExtraFile != "" {
		cfg, err := parseDefaultsExtraFile(*mysqlDefaultsExtraFile)
		if err != nil {
			return err
		}
		if cfg.Password != "" {
			*mysqlPassword = cfg.Password
		}
		if cfg.User != "" {
			*mysqlUser = cfg.User
		}
	}

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
			line, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			*mysqlQuery = *mysqlQuery + line
			line = strings.TrimSpace(line)
			if len(line) > 1 && line[len(line)-1] == ';' {
				break
			}
		}
	} else {
		if _, err := os.Stat(*mysqlQuery); err == nil {
			buf, err := ioutil.ReadFile(*mysqlQuery)
			if err == nil {
				*mysqlQuery = string(buf)
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

func save2csv(filePath string, rows *sql.Rows) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// set table header with column name
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	w.Write(columns)

	// set every table rows
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		columns := make([]sql.NullString, len(columns))
		cols := make([]interface{}, len(columns))
		for i := range columns {
			cols[i] = &columns[i]
		}

		if err := rows.Scan(cols...); err != nil {
			return err
		}

		values := make([]string, len(columns))
		for i, col := range columns {
			values[i] = col.String
		}
		w.Write(values)
	}
	return err
}

type mysqlConfig struct {
	User     string `toml:"mysql"`
	Password string `toml:"password"`
}

// parseDefaultsExtraFile dealwith --defaults-extra-file arg
func parseDefaultsExtraFile(file string) (mysqlConfig, error) {
	var cfg mysqlConfig
	c, err := ini.Load(file)
	if err != nil {
		return cfg, err
	}
	cfg.User = c.Section("mysql").Key("user").String()
	cfg.Password = c.Section("mysql").Key("password").String()
	return cfg, err
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
	if strings.HasSuffix(strings.ToLower(cfg.File), ".csv") {
		err = save2csv(cfg.File, rows)
	} else {
		err = save2xlsx(cfg.File, rows)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("save data into file: '%s'\n", cfg.File)
}
