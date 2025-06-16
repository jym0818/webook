package main

func main() {
	server := InitServer()
	server.Run(":8080")
}
