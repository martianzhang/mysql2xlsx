package common

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"gopkg.in/ini.v1"
)

// Config mysql2xlsx config
type Config struct {
	// [mysql] config
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

// ParseFlag parse cmd flags
func ParseFlag() (Config, error) {
	var cfg Config

	mysqlHost := flag.String("h", "localhost", "mysql host")
	mysqlUser := flag.String("u", "", "mysql user name")
	mysqlPassword := flag.String("p", "", "mysql password")
	mysqlDatabase := flag.String("d", "", "mysql database name")
	mysqlPort := flag.String("P", "3306", "mysql port")
	mysqlCharset := flag.String("c", "utf8mb4", "mysql default charset")
	mysqlDefaultsExtraFile := flag.String("defaults-extra-file", "", "mysql --defaults-extra-file arg")

	mysqlQuery := flag.String("q", "", "select query")
	fileName := flag.String("f", "", "save query result into file, default to stdout")
	flag.Parse()

	if *mysqlDefaultsExtraFile != "" {
		cfg, err := parseDefaultsExtraFile(*mysqlDefaultsExtraFile)
		if err != nil {
			return cfg, err
		}
		if cfg.Password != "" {
			*mysqlPassword = cfg.Password
		}
		if cfg.User != "" {
			*mysqlUser = cfg.User
		}
		if cfg.Charset != "" {
			*mysqlCharset = cfg.Charset
		}
	}

	if *mysqlPassword == "" {
		fmt.Print("Password:")
		password, err := gopass.GetPasswd()
		if err != nil {
			return cfg, err
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
				return cfg, err
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

	// use abs path
	pwd, err := os.Getwd()
	if err != nil {
		return cfg, err
	}
	if !strings.HasPrefix(*fileName, "/") &&
		(*fileName != "" && *fileName != "stdout") {
		*fileName = pwd + "/" + *fileName
	}

	cfg = Config{
		User:     *mysqlUser,
		Host:     *mysqlHost,
		Port:     *mysqlPort,
		Password: *mysqlPassword,
		Database: *mysqlDatabase,
		Charset:  *mysqlCharset,

		Query: *mysqlQuery,
		File:  *fileName,
	}

	return cfg, err
}

// parseDefaultsExtraFile parse --defaults-extra-file file
func parseDefaultsExtraFile(file string) (Config, error) {
	var cfg Config
	c, err := ini.Load(file)
	if err != nil {
		return cfg, err
	}

	// get config from [mysql] section
	cfg.User = c.Section("mysql").Key("user").String()
	cfg.Password = c.Section("mysql").Key("password").String()
	cfg.Charset = c.Section("mysql").Key("default-character-set").String()

	return cfg, err
}
