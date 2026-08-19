package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"api-audesp/database"
	"api-audesp/handlers"

	"github.com/rs/cors"
)

//go:embed public
var publicFiles embed.FS

func main() {
	// 1. Inicializa a conexão com o PostgreSQL
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

	// >>> CÓDIGO CORRIGIDO: SERVIDOR DE ARQUIVOS EMBUTIDO <<<
	// Extrai a subpasta "public" do sistema de arquivos embutido no executável
	publicFS, err := fs.Sub(publicFiles, "public")
	if err != nil {
		log.Fatal("❌ Erro ao embutir pasta public: ", err)
	}

	// Serve os arquivos HTML diretamente da memória (Resolve o 404 definitivamente)
	mux.Handle("/", http.FileServer(http.FS(publicFS)))

	// 4. Configuração de CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	// 5. Sobe o servidor
	porta := ":8080"
	fmt.Printf("🚀 Servidor da API e Frontend escutando na porta %s\n", porta)
	
	if err := http.ListenAndServe(porta, handler); err != nil {
		log.Fatal("❌ Erro ao iniciar o servidor: ", err)
	}
}
