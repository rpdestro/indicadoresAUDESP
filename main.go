cat << 'EOF' > main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"api-audesp/database"
	"api-audesp/handlers"

	"github.com/rs/cors"
)

func main() {
	// >>> DIAGNÓSTICO EDUCATIVO: Fazendo o Go nos contar o que ele enxerga <<<
	arquivos, err := os.ReadDir("public")
	if err != nil {
		fmt.Println("⚠️ ALERTA: A pasta 'public' NÃO foi encontrada pelo Golang:", err)
	} else {
		fmt.Println("📁 Pasta 'public' ENCONTRADA! Arquivos vistos pelo servidor:")
		for _, arq := range arquivos {
			fmt.Println("  -", arq.Name())
		}
	}

	// 1. Inicializa a conexão com o PostgreSQL no Supabase
	database.Connect()

	// 2. Cria o roteador padrão do Go
	mux := http.NewServeMux()

	// 3. Rotas da API
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "online", "mensagem": "API AUDESP rodando perfeitamente!"}`))
	})
	mux.HandleFunc("/api/indicadores", handlers.SalvarIndicador)
	mux.HandleFunc("/api/planejamento", handlers.BuscarPlanejamento)
	mux.HandleFunc("/api/dashboard", handlers.ListarDashboard)
	mux.HandleFunc("/api/compliance", handlers.PainelConformidade)
	mux.HandleFunc("/api/detalhes", handlers.DetalhesSecretaria)

	// 4. Servidor de Arquivos Estáticos (Refinado)
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)

	// 5. Configuração de CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	// 6. Sobe o servidor na porta 8080
	porta := ":8080"
	fmt.Printf("🚀 Servidor da API e Frontend escutando na porta %s\n", porta)
	
	err = http.ListenAndServe(porta, handler)
	if err != nil {
		log.Fatal("❌ Erro ao iniciar o servidor: ", err)
	}
}
EOF
