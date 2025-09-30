package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() {
	connStr := "user=admin password=admin dbname=peoplepulse_db host=172.23.118.68 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Falha ao conectar ao DB: ", err)
	}
	defer db.Close()

	db.Exec("TRUNCATE TABLE kpis RESTART IDENTITY CASCADE;")
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE;")
	fmt.Println("Tabelas limpas com sucesso.")

	email := "edgard@peoplepulse.com.br"
	password := []byte("12345")

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Falha ao gerar hash: ", err)
	}

	var insertedUserID int
	err = db.QueryRow("INSERT INTO users (name, email, password_hash, position, role) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		"Edgard Masso", email, string(hashedPassword), "Diretor", "diretoria").Scan(&insertedUserID)
	if err != nil {
		log.Fatal("Falha ao inserir novo usuário: ", err)
	}

	fmt.Println("Novo usuário 'edgard@peoplepulse.com.br' inserido com sucesso!")
}
