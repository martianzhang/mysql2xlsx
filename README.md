# Introduction

Save SQL query result into `xlsx`, `csv` file or print as `ascii` table.

## Install

```bash
go get -u github.com/martianzhang/mysql2xlsx
```

## Cross platform compile

```bash
# GOOS: darwin linux windows
# GOARCH: amd64
GOOS={GOOS} GOARCH={GOARCH} go build
```

## Example

```bash
./mysql2xlsx  -h 127.0.0.1 -P 3306 -d database -u root -f result.xlsx -q "select * from tbl"
Password:<hidden input>
save data into file: '/path/to/mysql2xlsx/result.xlsx'
```

## Get Usage

```bash
./mysql2xlsx --help
Usage of ./mysql2xlsx:
  -P string
        mysql port (default "3306")
  -c string
        mysql default charset (default "utf8mb4")
  -d string
        mysql database name
  -defaults-extra-file string
        mysql --defaults-extra-file arg
  -f string
        save query result into file, default to stdout
  -h string
        mysql host (default "localhost")
  -p string
        mysql password
  -q string
        select query
  -u string
        mysql user name
```
