# FAQ

## go get package timeout

```txt
package golang.org/x/sys: unrecognized import path "golang.org/x/sys": https fetch: Get "https://golang.org/x/sys?go-get=1": dial tcp 216.239.37.1:443: i/o timeout
```

```bash
mkir -p $GOPATH/github.com/golang && cd $GOPATH/github.com/golang
git clone https://github.com/golang/sys
ln -s $GOPATH/github.com/golang $GOPATH/golang.org/x
```
