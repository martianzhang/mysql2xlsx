# Introduction

Save SQL query result into `xlsx`, `csv` file or print as `ascii` table.

## Install

```bash
git clone github.com/martianzhang/mysql2xlsx
cd mysql2xlsx && make build
```

## Cross platform compile

```bash
# GOOS: darwin linux windows
# GOARCH: amd64
GOOS={GOOS} GOARCH={GOARCH} go build
```

## Example

```bash
./mysql2xlsx  --host 127.0.0.1 -port 3306 --databasee database --user root --file result.xlsx --query "select * from tbl"
Password:<hidden input>
```

## Get Usage

```bash
./mysql2xlsx --help
  -bom
    	csv file with UTF8 BOM
  -charset string
    	mysql default charset (default "utf8mb4")
  -database string
    	Database to use. (default "information_schema")
  -defaults-extra-file string
    	mysql --defaults-extra-file arg
  -file string
    	save query result into file, (default "stdout")
  -host string
    	Connect to host. (default "127.0.0.1")
  -password string
    	Password to use when connecting to server. If password is not given it's asked from the tty.
  -port string
    	Port number to use for connection. (default "3306")
  -preview int
    	preview result file, print first N lines
  -query string
    	select query
  -socket string
    	The socket file to use for connection.
  -user string
    	User for login if not current user.
```

## Other use case

```bash
# preview xlsx file
mysql2xlsx -preview 10 -file test.xlsx
```

```bash
# simple mysql client
mysql2xlsx --defaults-xtra-file my.cnf -query "select 1"
+---+
| 1 |
+---+
| 1 |
+---+
```
