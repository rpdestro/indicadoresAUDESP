package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Importa o driver do Postgres silenciosamente
)

// DB é a variável global que segurará o "pool" de conexões com o banco
var DB *sql.DB

// Connect inicializa a conexão com o banco de dados
func Connect() {
	// 1. Carrega as variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. O sistema tentará usar variáveis de ambiente do SO.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ ERRO FATAL: A variável DATABASE_URL não foi definida no arquivo .env.")
	}

	// 2. Abre a conexão usando o driver "postgres"
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ ERRO FATAL ao tentar configurar a conexão:", err)
	}

	// 3. Testa efetivamente se o banco responde (Ping)
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ ERRO FATAL: Banco de dados não respondeu ao Ping:", err)
	}

	fmt.Println("✅ Conexão com o banco Supabase estabelecida com sucesso!")
}
