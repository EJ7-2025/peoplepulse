package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("sua_chave_secreta_super_secreta")
// DENTRO DO CODESPACE, 'localhost' funciona para se comunicar com o container Docker
var connStr = "user=admin password=admin dbname=peoplepulse_db host=localhost sslmode=disable"

// As structs (User, LoginCredentials, Claims, KPI) permanecem as mesmas...
type User struct { ID int; Name string; Email string; PasswordHash string; Position sql.NullString; Role sql.NullString }
type LoginCredentials struct { Email string `json:"email"`; Password string `json:"password"` }
type Claims struct { UserID int; Role string; jwt.RegisteredClaims }
type KPI struct { ID int `json:"id"`; Title string `json:"title"`; Value int `json:"value"` }

// NOVA FUNÇÃO para inicializar o banco de dados
func initializeDatabase() {
	db, err := sql.Open("postgres", connStr)
	if err != nil { log.Fatal("Falha ao conectar ao DB para inicialização: ", err) }
	defer db.Close()

	// Cria a tabela de usuários se ela não existir
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			position VARCHAR(100),
			role VARCHAR(50)
		);
	`)
	if err != nil { log.Fatal("Falha ao criar tabela users: ", err) }

	// Cria a tabela de KPIs se ela não existir
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS kpis (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			value INT NOT NULL,
			user_id INT NOT NULL,
			CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil { log.Fatal("Falha ao criar tabela kpis: ", err) }

	// Verifica se o usuário de teste já existe
	var userCount int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", "edgard@peoplepulse.com.br").Scan(&userCount)
	if userCount == 0 {
		log.Println("Usuário de teste não encontrado, criando...")
		password := []byte("12345")
		hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

		var insertedUserID int
		err = db.QueryRow("INSERT INTO users (name, email, password_hash, position, role) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			"Edgard Masso", "edgard@peoplepulse.com.br", string(hashedPassword), "Diretor", "diretoria").Scan(&insertedUserID)
		if err != nil { log.Fatal("Falha ao inserir usuário de teste: ", err) }

		// Insere os KPIs para o usuário de teste
		_, err = db.Exec("INSERT INTO kpis (title, value, user_id) VALUES ('Resolução de Tickets', 85, $1), ('Commits no Repositório', 60, $1)", insertedUserID)
		if err != nil { log.Fatal("Falha ao inserir KPIs de teste: ", err) }
		log.Println("Usuário de teste e KPIs criados com sucesso.")
	} else {
		log.Println("Banco de dados já inicializado.")
	}
}

func Login(c *gin.Context) {
	var creds LoginCredentials
	if err := c.ShouldBindJSON(&creds); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida"}); return }
	db, _ := sql.Open("postgres", connStr); defer db.Close()
	var user User
	err := db.QueryRow("SELECT id, name, email, password_hash, position, role FROM users WHERE email = $1", creds.Email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Position, &user.Role)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"}); return }
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"}); return }
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{UserID: user.ID, Role: user.Role.String, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expirationTime)}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": user.Role.String})
}

func GetKPIs(c *gin.Context) {
    // (O código desta função permanece o mesmo)
}

func main() {
    // A INICIALIZAÇÃO OCORRE AQUI
    initializeDatabase()

	router := gin.Default(); router.Use(cors.Default()); router.POST("/login", Login); router.GET("/kpis", GetKPIs); router.Run(":8080")
}