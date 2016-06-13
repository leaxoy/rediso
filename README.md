# This is a redis client written in golang.

## Basic usage.

```language:go
func main() {
	conn := rediso.Connect(nil)
	conn.ExecQuery(rediso.SET, "name", "leaxoy")
	conn.ExecQuery(rediso.GET, "name")
	s := conn.Ping() // should return "+PONG"
	println(s)
	db := conn.DB() // should return 0
	println(db)
	conn.SelectDB(1) // change db to 1
	db = conn.DB()   //should return 1
	println(db
  conn.Close()
}
```

`Config` is redis config.
```
{
  host    string        // default is `localhost`
  port    int           // default is 6379
  ssl     bool          // use ssl. default is false
  timeout time.Duration // default is 5's
  db      int           // default is 0
}
```
