package rediso

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

// Conn is interface
type Conn interface {
	Auth(auth string) bool
	Conn() *Client
	DB() int
	Exec(query *Query)
	ExecQuery(command string, args ...string) interface{}
	SelectDB(db int) *DB
	Info() string
	HostInfo() string
	Ping() string
	Close()
}

// Client is Redis conn.
type Client struct {
	db   *DB
	conn net.Conn
	conf *Config
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
	auth    string        // default is null
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
func (c *Client) Auth(auth string) bool {
	q := buildCommand("Auth", auth)
	c.Exec(q)
	return q.ret.(string) == "+OK\n"
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
	query.ret = string(buf[:n])
}

// ExecQuery is a directly function
func (c *Client) ExecQuery(command string, args ...string) interface{} {
	q := buildCommand(command, args...)
	c.Exec(q)
	return q.ret
}

// Info is a method that return info of redis server.
func (c *Client) Info() string {
	q := buildCommand("INFO")
	c.Exec(q)
	return q.ret.(string)
}

// HostInfo return the redis server info.
func (c *Client) HostInfo() string {
	auth := c.conf.auth
	if c.conf.auth == "" {
		auth = "nil"
	}
	return fmt.Sprintf("Host:\t%s\nPort:\t%d\nDB:\t%d\nAuth:\t%s", c.conf.host, c.conf.port, c.conf.db, auth)
}

// Conn return the current client.
func (c *Client) Conn() *Client {
	return c
}

// NewRedis create new redis conn
func NewRedis(config *Config) *Client {
	if config == nil {
		config = new(Config)
		config.host = "localhost"
		config.port = 6379
		config.db = 0
		config.timeout = 5
		config.ssl = false
		config.auth = ""
	}
	addr := config.host + ":" + strconv.Itoa(config.port)
	conn, err := net.DialTimeout("tcp", addr, time.Second*config.timeout)
	if err != nil {
		panic(err)
	}
	client := new(Client)

	client.conn = conn
	client.conf = config
	if config.auth != "" {
		ok := client.Auth(config.auth)
		if !ok {
			panic("Auth failed")
		}
	}
	client.SelectDB(config.db)
	return client
}

// DB is redis db.
type DB struct {
	index      int
	charset    string
	connection *Client
	query      *Query
}
