package main

import (
	"whatsapp/db"
	routes "whatsapp/internal/routes"
)

func main() {
	db.GetDB()
	r := routes.SetUpRouter()
	r.Run(":8080")
}
