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
Usage of ./mysql2xlsx:
  -bom
        csv file with UTF8 BOM
  -charset string
        mysql default charset (default "utf8mb4")
  -database string
        mysql database name (default "information_schema")
  -defaults-extra-file string
        mysql --defaults-extra-file arg
  -file string
        save query result into file, (default "stdout")
  -host string
        mysql host (default "127.0.0.1")
  -password string
        mysql password
  -port string
        mysql port (default "3306")
  -query string
        select query
  -user string
        mysql user name
```
