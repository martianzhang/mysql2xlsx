package common

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/howeyc/gopass"
	ini "gopkg.in/ini.v1"
)

// Config mysql2xlsx config
type Config struct {
	// [mysql] config
	User     string
	Host     string
	Socket   string
	Port     string
	Password string
	Database string
	Charset  string

	// other config
	Query   string // select query
	File    string // storage file abs path
	BOM     bool   // add BOM file header
	Preview int    // preview xlsx file, print first N lines
}

// Cfg global config
var Cfg Config

// ParseFlag parse cmd flags
func ParseFlag() error {

	mysqlHost := flag.String("host", "127.0.0.1", "Connect to host.")
	mysqlSocket := flag.String("socket", "", "The socket file to use for connection.")
	mysqlUser := flag.String("user", "", "User for login if not current user.")
	mysqlPassword := flag.String("password", "", "Password to use when connecting to server. If password is not given it's asked from the tty.")
	mysqlDatabase := flag.String("database", "information_schema", "Database to use.")
	mysqlPort := flag.String("port", "3306", "Port number to use for connection.")
	mysqlCharset := flag.String("charset", "utf8mb4", "mysql default charset")
	mysqlDefaultsExtraFile := flag.String("defaults-extra-file", "", "mysql --defaults-extra-file arg")

	mysqlQuery := flag.String("query", "", "select query")
	filename := flag.String("file", "", `save query result into file, (default "stdout")`)
	var bom *bool
	if runtime.GOOS != "windows" {
		bom = flag.Bool("bom", false, "csv file with UTF8 BOM")
	} else {
		*bom = true
	}
	previewXLSX := flag.Int("preview", 0, "preview xlsx file, print first N lines")

	flag.Parse()

	if *previewXLSX != 0 && *filename != "" {
		Cfg.File = *filename
		Cfg.Preview = *previewXLSX
		return nil
	}

	if *mysqlDefaultsExtraFile != "" {
		err := parseDefaultsExtraFile(*mysqlDefaultsExtraFile)
		if err != nil {
			return err
		}
	}

	if Cfg.User != "" {
		*mysqlUser = Cfg.User
	}
	if *mysqlUser == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if Cfg.Password != "" {
		*mysqlPassword = Cfg.Password
	}

	if Cfg.Charset != "" {
		*mysqlCharset = Cfg.Charset
	}

	if !strings.HasPrefix(strings.ToLower(*mysqlCharset), "utf") {
		*bom = false
	}

	if Cfg.Query != "" {
		*mysqlQuery = Cfg.Query
	}

	if *mysqlPassword == "" {
		fmt.Print("Password:")
		password, err := gopass.GetPasswd()
		if err != nil {
			return err
		}
		*mysqlPassword = strings.TrimSpace(string(password))
	}

	// read query interactive
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
	}

	// test read from file
	if _, err := os.Stat(*mysqlQuery); err == nil {
		buf, err := ioutil.ReadFile(*mysqlQuery)
		if err == nil {
			*mysqlQuery = string(buf)
		}
	}

	*filename = strings.TrimSpace(*filename)
	if *filename == "" {
		*filename = "stdout"
	}

	// use abs path
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(*filename, "/") &&
		(*filename != "" && *filename != "stdout") {
		*filename = pwd + "/" + *filename
	}

	Cfg = Config{
		User:     *mysqlUser,
		Host:     *mysqlHost,
		Socket:   *mysqlSocket,
		Port:     *mysqlPort,
		Password: *mysqlPassword,
		Database: *mysqlDatabase,
		Charset:  *mysqlCharset,

		Query: *mysqlQuery,
		File:  *filename,
		BOM:   *bom,
	}

	return err
}

// parseDefaultsExtraFile parse --defaults-extra-file file
func parseDefaultsExtraFile(file string) error {
	c, err := ini.Load(file)
	if err != nil {
		return err
	}

	// get config from [mysql] section
	Cfg.User = c.Section("client").Key("user").String()
	Cfg.Password = c.Section("client").Key("password").String()
	Cfg.Charset = c.Section("client").Key("default-character-set").String()

	return err
}
