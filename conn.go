package rediso

import (
	"bufio"
	"net"
	"strconv"
	"time"
)

// Conn is interface
type Conn interface {
	SelectDB(db int) *DB
	Ping() string
	Auth(auth string)
	Info() string
	HostInfo() string
	DB() int
	Exec(query *Query)
	Close()
	Conn() *Client
	ExecQuery(command string, args ...string)
}

// Client is Redis conn.
type Client struct {
	db   *DB
	conn net.Conn
}

// Query is struct.
type Query struct {
	query string
	ret   interface{}
}

// FormatCommand is a function
func (q *Query) FormatCommand(command string, args ...string) {
	q.query = command
	for _, arg := range args {
		q.query = q.query + " " + arg
	}
	q.query = q.query + "\n"
}

type exec interface {
	exec(query *Query)
}

// Config is redis connection config.
type Config struct {
	host    string        // default is `localhost`
	port    int           // default is 6379
	ssl     bool          // use ssl. default is false
	timeout time.Duration // default is 5's
	db      int           // default is 0
}

// Connect is a function to connect redis server and return the conn.
func Connect(config *Config) Conn {
	return NewRedis(config)
}

// DB method return the current DB index.
func (c *Client) DB() int {
	return c.db.index
}

// SelectDB is a function that change db index.
func (c *Client) SelectDB(db int) *DB {
	database := &DB{index: db,
		charset:    "utf-8",
		connection: c,
		query:      new(Query)}
	q := buildCommand("SELECT", strconv.Itoa(db))
	c.Exec(q)
	c.db = database
	return database
}

// Auth is a function to auth to redis server.
func (c *Client) Auth(auth string) {
}

// Ping is a function that ping to redis server.
func (c *Client) Ping() string {
	q := buildCommand("PING")
	c.Exec(q)
	return q.ret.(string)
}

// Close is a function that close connection to redis server.
func (c *Client) Close() {
	c.conn.Close()
}

func buildCommand(command string, args ...string) *Query {
	q := new(Query)
	q.query = command
	for _, arg := range args {
		q.query = q.query + " " + arg
	}
	q.query = q.query + "\n"
	return q
}

// Exec is a function to exec command on redis conn
func (c *Client) Exec(query *Query) {
	// n, err := bufio.NewWriter(c.conn).WriteString(query.query)
	_, err := c.conn.Write([]byte(query.query))
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024*4)
	n, err := bufio.NewReader(c.conn).Read(buf)
	if err != nil {
		panic(err)
	}
	println(string(buf[:n]))
	println("Total:", n, "bytes received.")
	println("------")
	query.ret = string(buf[:n])
}

// ExecQuery is a directly function
func (c *Client) ExecQuery(command string, args ...string) {
	q := buildCommand(command, args...)
	c.Exec(q)
}

// Info is a method that return info of redis server.
func (c *Client) Info() string {
	q := buildCommand("INFO")
	c.Exec(q)
	return q.ret.(string)
}

// HostInfo return the redis server info.
func (c *Client) HostInfo() string {
	info := c.Info()
	return info
}

// Conn return the current client.
func (c *Client) Conn() *Client {
	return c
}

// NewRedis create new redis conn
func NewRedis(config *Config) *Client {
	connection := new(Client)
	host := "localhost"
	port := 6379
	db := 0
	var timeout time.Duration
	timeout = 5
	if config != nil {
		if config.host != "" {
			host = config.host
		}
		if config.port > 1024 {
			port = config.port
		}
		if config.db >= 0 && config.db <= 16 {
			db = config.db
		}
		if config.timeout > 0 {
			timeout = config.timeout
		}
	}
	addr := host + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", addr, time.Second*timeout)
	if err != nil {
		panic(err)
	}
	connection.conn = conn
	connection.SelectDB(db)
	return connection
}

// DB is redis db.
type DB struct {
	index      int
	charset    string
	connection *Client
	query      *Query
}
