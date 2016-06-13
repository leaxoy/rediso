package main

import "github.com/leaxoy/rediso"

func main() {
	conn := rediso.Connect(nil)
	conn.Exec(rediso.SET, "name", "leaxoy")
	conn.Exec(rediso.GET, "name")
	s := conn.Ping() // should return "+PONG"
	println(s)
	db := conn.DB() // should return 0
	println(db)
	conn.SelectDB(1) // change db to 1
	db = conn.DB()   // should return 1
	println(db)
	println(conn.HostInfo())
	conn.Close()
}
