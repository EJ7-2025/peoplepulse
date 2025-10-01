package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
    // No Codespace, o host do banco de dados é 'localhost'
	connStr := "user=admin password=admin dbname=peoplepulse_db host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil { log.Fatal("Falha ao conectar ao DB: ", err) }
	defer db.Close()

	// Apaga as tabelas antigas se existirem, na ordem correta
	db.Exec("DROP TABLE IF EXISTS kpis;")
	db.Exec("DROP TABLE IF EXISTS users;")
	fmt.Println("Tabelas antigas removidas com sucesso.")

	// Cria as tabelas na ordem correta
	_, err = db.Exec(`CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, email VARCHAR(100) UNIQUE NOT NULL, password_hash VARCHAR(255) NOT NULL, position VARCHAR(100), role VARCHAR(50));`)
	if err != nil { log.Fatal("Falha ao criar tabela users: ", err) }

	_, err = db.Exec(`CREATE TABLE kpis (id SERIAL PRIMARY KEY, title VARCHAR(100) NOT NULL, value INT NOT NULL, user_id INT NOT NULL, CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE);`)
	if err != nil { log.Fatal("Falha ao criar tabela kpis: ", err) }
	fmt.Println("Tabelas 'users' e 'kpis' criadas com sucesso.")

	// Insere o usuário de teste
	email := "edgard@peoplepulse.com.br"
	password := []byte("12345")
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	var insertedUserID int
	err = db.QueryRow("INSERT INTO users (name, email, password_hash, position, role) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		"Edgard Masso", email, string(hashedPassword), "Diretor", "diretoria").Scan(&insertedUserID)
	if err != nil { log.Fatal("Falha ao inserir novo usuário: ", err) }
	fmt.Println("Novo usuário 'edgard@peoplepulse.com.br' inserido com sucesso!")
}