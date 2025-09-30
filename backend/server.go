package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("sua_chave_secreta_super_secreta")

type User struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"-"`
	Position     sql.NullString `json:"position"`
	Role         sql.NullString `json:"role"`
}
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Claims struct {
	UserID int    `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
type KPI struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Value int    `json:"value"`
}

func Login(c *gin.Context) {
	var creds LoginCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida"})
		return
	}
	connStr := "user=admin password=admin dbname=peoplepulse_db host=172.23.118.68 sslmode=disable"
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()
	var user User
	err := db.QueryRow("SELECT id, name, email, password_hash, position, role FROM users WHERE email = $1", creds.Email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Position, &user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{UserID: user.ID, Role: user.Role.String, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expirationTime)}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível gerar o token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": user.Role.String})
}

// VERSÃO DE TESTE DA FUNÇÃO GetKPIs - SEM SEGURANÇA
func GetKPIs(c *gin.Context) {
	log.Println("--- Rota de Teste /kpis ACESSADA ---")
	// Apenas conecta e busca os KPIs do usuário com ID=1 (nosso usuário de teste)
	connStr := "user=admin password=admin dbname=peoplepulse_db host=172.23.118.68 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("!!! ERRO no teste: Falha ao abrir conexão com DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao conectar ao DB"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, value FROM kpis WHERE user_id = $1", 1) // ID do usuário fixo em 1
	if err != nil {
		log.Println("!!! ERRO no teste: Falha ao executar query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar KPIs"})
		return
	}
	defer rows.Close()

	var kpis []KPI
	for rows.Next() {
		var k KPI
		if err := rows.Scan(&k.ID, &k.Title, &k.Value); err != nil {
			log.Println("!!! ERRO no teste: Falha ao processar dados:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados"})
			return
		}
		kpis = append(kpis, k)
	}
	log.Println(">>> Sucesso no teste: KPIs encontrados e enviados.")
	c.IndentedJSON(http.StatusOK, kpis)
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/login", Login)
	router.GET("/kpis", GetKPIs)
	router.Run(":8080")
}
