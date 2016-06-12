package main

import "bufio"
import "net"
import "strconv"
import "sync"
import "time"
import "os"

// import "strings"

// Conn is interface
type Conn interface {
	SelectDB(db int) *DB
	Ping()
	Auth(auth string)
	Close()
}

type pool struct {
	size  uint
	conns []*Connection
}

// Query is struct.
type Query struct {
	q   string
	ret string
}

func formatCommand(q *Query, args ...string) {

}

// FormatCommand is a function
func (q *Query) FormatCommand(command string, args ...string) {
	q.q = command
	for _, arg := range args {
		q.q = q.q + " " + arg
	}
	q.q = q.q + "\r\n"
}

type exec interface {
	exec(query *Query)
}

// Connection is Redis Connection.
type Connection struct {
	host  string
	port  string
	ssl   bool
	db    int
	user  string
	conn  net.Conn
	mutex *sync.Mutex
}

// Config is redis connection config.
type Config struct {
	host    string        // default is `localhost`
	port    string        // default is 6379
	ssl     bool          // use ssl. default is false
	timeout time.Duration // default is 5's
}

// Connect is a function to connect redis server and return the Connection.
func Connect(config *Config) *Connection {
	connection := new(Connection)
	host := "localhost"
	port := 6379
	db := 0

	addr := "127.0.0.1:6379"
	addr = host + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	connection.conn = conn
	connection.SelectDB(db)
	return connection
}

// SelectDB is a function that change db index.
func (c *Connection) SelectDB(db int) *DB {
	c.db = db
	database := &DB{index: db,
		charset:    "utf-8",
		connection: c,
		query:      &Query{q: "", ret: ""}}
	database.query.FormatCommand("SELECT", "1")
	database.exec()
	return database
}

// Auth is a function to auth to redis server.
func (c *Connection) Auth(auth string) {
	c.user = auth
}

// Ping is a function that ping to redis server.
func (c *Connection) Ping() {
	_, err := c.conn.Write([]byte("PING\r\n"))
	if err != nil {
		panic(err)
	}
	msg, err := bufio.NewReader(c.conn).ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	println(string(msg))
}

// Close is a function that close connection to redis server.
func (c *Connection) Close() {
	c.conn.Close()
}

// DB is redis db.
type DB struct {
	index      int
	charset    string
	connection *Connection
	query      *Query
}

func (db *DB) exec() {
	// println(db.query.q)
	n, err := db.connection.conn.Write([]byte(db.query.q))
	if err != nil {
		panic(err)
	}
	println("\tsend", n)
	meg, err := bufio.NewReader(db.connection.conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	println("\t", string(meg))
}

// Set is a function to set.
func (db *DB) Set(key, val string) {
	db.connection.mutex.Lock()
	defer db.connection.mutex.Unlock()
	db.query.FormatCommand("SET", key, val)
	db.exec()
}

// Get is a function can get a val.
func (db *DB) Get(key string) string {
	db.query.FormatCommand("GET", key)
	db.exec()
	return ""
}

func main() {
	name := "lily"
	if len(os.Args) == 2 {
		name = os.Args[1]
	}
	config := new(Config)
	conn := Connect(config)
	db := conn.SelectDB(1)
	db.Set("name", name)
	db.Get("name")
}
