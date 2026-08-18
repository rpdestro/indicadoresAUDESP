package handlers

import (
	"fmt"
	"net/http"

	// O pacote de banco de dados foi removido temporariamente.
	// Nós o colocaremos de volta quando formos fazer o INSERT real.
	"github.com/xuri/excelize/v2"
)

// ImportarPlanilha recebe um arquivo Excel e processa os dados da LOA para o banco
func ImportarPlanilha(w http.ResponseWriter, r *http.Request) {
	// 1. Configura o recebimento do arquivo
	err := r.ParseMultipartForm(10 << 20) // Limite de upload: 10MB
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("planilha_loa")
	if err != nil {
		http.Error(w, "Erro ao ler o arquivo enviado", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Abre o Excel em memória
	f, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Erro ao abrir a planilha. Certifique-se que é formato XLSX.", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// 3. Lê as linhas da aba principal (Ajustaremos o nome da aba depois)
	rows, err := f.GetRows("Planilha1")
	if err != nil {
		http.Error(w, "Aba de dados não encontrada no arquivo", http.StatusBadRequest)
		return
	}

	// 4. Loop de varredura das linhas
	for indice, row := range rows {
		// Pula o cabeçalho (linha 0)
		if indice == 0 {
			continue 
		}
		
		// Lógica de mapeamento (Será ajustada conforme a planilha real)
		// Supondo que row[0] = Programa, row[1] = Ação, row[2] = Und Medida
		if len(row) >= 3 {
			programa := row[0]
			acao := row[1]
			undMedida := row[2]
			
			// Aqui simulamos o INSERT. Quando formos integrar com o banco,
			// voltaremos a importar "api-audesp/database" lá no topo.
			fmt.Printf("Pronto para inserir - Programa: %s | Ação: %s | Medida: %s\n", programa, acao, undMedida)
		}
	}

	// Retorna sucesso para o frontend
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensagem": "Planilha processada com sucesso!"}`))
}