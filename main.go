package main

import (
	"fmt"
	"log"
	"net/http"

	"api-audesp/database"
	"api-audesp/handlers"

	"github.com/rs/cors"
)

func main() {
	// 1. Inicializa a conexão com o PostgreSQL no Supabase
	database.Connect()

	// 2. Cria o roteador padrão do Go
	mux := http.NewServeMux()

	// 3. Rotas da API (O Cérebro do sistema)
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

	// >>> NOVO: Servidor de Arquivos Estáticos (O Rosto do sistema) <<<
	// Avisa ao Go para pegar os arquivos HTML da pasta "public" e entregar no navegador
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	// 4. Configuração de CORS (Permite a comunicação segura)
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	// 5. Sobe o servidor na porta 8080
	porta := ":8080"
	fmt.Printf("🚀 Servidor da API e Frontend escutando na porta %s\n", porta)
	
	err := http.ListenAndServe(porta, handler)
	if err != nil {
		log.Fatal("❌ Erro ao iniciar o servidor: ", err)
	}
}