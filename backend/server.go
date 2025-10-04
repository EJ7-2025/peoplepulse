package main
import ("database/sql";"net/http";"strings";"time";"github.com/gin-contrib/cors";"github.com/gin-gonic/gin";"github.com/golang-jwt/jwt/v5";_ "github.com/lib/pq";"golang.org/x/crypto/bcrypt")
var jwtKey=[]byte("sua_chave_secreta_super_secreta")
var connStr="user=admin password=admin dbname=peoplepulse_db host=localhost sslmode=disable"
type User struct{ID int;Name string;Email string;PasswordHash string `json:"-"`;Position sql.NullString;Role sql.NullString}
type LoginCredentials struct{Email string `json:"email"`;Password string `json:"password"`}
type Claims struct{UserID int;Role string;jwt.RegisteredClaims}
type KPI struct{ID int `json:"id"`;Title string `json:"title"`;Value int `json:"value"`}
func Login(c *gin.Context){
var creds LoginCredentials
if err:=c.ShouldBindJSON(&creds);err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"Requisição inválida"});return}
db,_:=sql.Open("postgres",connStr);defer db.Close()
var user User
err:=db.QueryRow("SELECT id, name, email, password_hash, position, role FROM users WHERE email = $1",creds.Email).Scan(&user.ID,&user.Name,&user.Email,&user.PasswordHash,&user.Position,&user.Role)
if err!=nil{c.JSON(http.StatusUnauthorized,gin.H{"error":"Email ou senha inválidos"});return}
if err:=bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(creds.Password));err!=nil{c.JSON(http.StatusUnauthorized,gin.H{"error":"Email ou senha inválidos"});return}
expirationTime:=time.Now().Add(24*time.Hour);claims:=&Claims{UserID:user.ID,Role:user.Role.String,RegisteredClaims:jwt.RegisteredClaims{ExpiresAt:jwt.NewNumericDate(expirationTime)}}
token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims);tokenString,_:=token.SignedString(jwtKey)
c.JSON(http.StatusOK,gin.H{"token":tokenString,"role":user.Role.String})
}
func GetKPIs(c *gin.Context){
authHeader:=c.GetHeader("Authorization");if authHeader==""{c.JSON(http.StatusUnauthorized,gin.H{"error":"Cabeçalho de autorização não fornecido"});return}
tokenString:=strings.TrimPrefix(authHeader,"Bearer ");claims:=&Claims{}
token,err:=jwt.ParseWithClaims(tokenString,claims,func(token *jwt.Token)(interface{},error){return jwtKey,nil});if err!=nil||!token.Valid{c.JSON(http.StatusUnauthorized,gin.H{"error":"Token inválido"});return}
db,_:=sql.Open("postgres",connStr);defer db.Close()
rows,err:=db.Query("SELECT id, title, value FROM kpis WHERE user_id = $1",claims.UserID);if err!=nil{c.JSON(http.StatusInternalServerError,gin.H{"error":"Erro ao buscar KPIs"});return}
defer rows.Close();var kpis[]KPI
for rows.Next(){var k KPI;if err:=rows.Scan(&k.ID,&k.Title,&k.Value);err!=nil{c.JSON(http.StatusInternalServerError,gin.H{"error":"Erro ao processar dados dos KPIs"});return};kpis=append(kpis,k)}
c.IndentedJSON(http.StatusOK,kpis)
}
func main(){
router:=gin.Default();router.Use(cors.Default());router.POST("/login",Login);router.GET("/kpis",GetKPIs);router.Run(":8080")
}