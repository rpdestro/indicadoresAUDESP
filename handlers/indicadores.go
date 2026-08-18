package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"api-audesp/database"
	"api-audesp/models"
)

// 1. SalvarIndicador recebe os dados, salva o registro e processa as imagens
func SalvarIndicador(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "Erro ao processar formulário: "+err.Error(), http.StatusBadRequest)
		return
	}

	codExecutora := r.FormValue("und_executora")
	programa := r.FormValue("programa")
	acao := r.FormValue("acao")
	unidadeMedida := r.FormValue("und_medida")
	anoVigente, _ := strconv.Atoi(r.FormValue("ano_vigente"))
	quadrimestre, _ := strconv.Atoi(r.FormValue("quadrimestre"))
	metaPrevista, _ := strconv.ParseFloat(r.FormValue("meta_prevista"), 64)
	metaRealizada, _ := strconv.ParseFloat(r.FormValue("meta_realizada"), 64)
	localExecucao := r.FormValue("local_execucao")
	dataInicio := r.FormValue("data_inicio")
	dataFim := r.FormValue("data_fim")
	justificativa := r.FormValue("justificativa")

	var idMonitoramento string
	query := `
		INSERT INTO monitoramento_audesp 
		(cod_executora, programa, acao, unidade_medida, ano_vigente, quadrimestre, meta_prevista, meta_realizada, local_execucao, data_inicio, data_fim, justificativa)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	
	err = database.DB.QueryRow(query, codExecutora, programa, acao, unidadeMedida, anoVigente, quadrimestre, metaPrevista, metaRealizada, localExecucao, dataInicio, dataFim, justificativa).Scan(&idMonitoramento)
	if err != nil {
		fmt.Println("Erro no Banco de Dados:", err)
		http.Error(w, "Falha ao salvar relatório no banco de dados.", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["fotos"]
	if len(files) > 0 {
		os.MkdirAll("uploads", os.ModePerm)

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				continue 
			}
			defer file.Close()

			nomeUnico := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
			caminhoDestino := filepath.Join("uploads", nomeUnico)

			dst, err := os.Create(caminhoDestino)
			if err != nil {
				continue
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				continue
			}

			queryAnexo := `INSERT INTO anexos_monitoramento (id_monitoramento, caminho_arquivo, nome_original) VALUES ($1, $2, $3)`
			_, err = database.DB.Exec(queryAnexo, idMonitoramento, caminhoDestino, fileHeader.Filename)
			if err != nil {
				fmt.Println("Erro ao gravar anexo no banco:", err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"mensagem": "Indicador e fotos salvos com sucesso!"}`))
}

// 2. BuscarPlanejamento devolve as ações da LOA baseadas na Unidade Executora
func BuscarPlanejamento(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	executora := r.URL.Query().Get("executora")
	if executora == "" {
		http.Error(w, "Código da executora é obrigatório", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query("SELECT programa, acao, und_medida FROM planejamento_loa WHERE cod_executora = $1", executora)
	if err != nil {
		http.Error(w, "Erro ao buscar planejamento", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []models.PlanejamentoLOA
	for rows.Next() {
		var p models.PlanejamentoLOA
		if err := rows.Scan(&p.Programa, &p.Acao, &p.UnidadeMedida); err != nil {
			continue
		}
		lista = append(lista, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

// 3. ListarDashboard busca todos os indicadores salvos
func ListarDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT id, cod_executora, acao, meta_prevista, meta_realizada, local_execucao, ano_vigente, quadrimestre 
		FROM monitoramento_audesp 
		ORDER BY id DESC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Erro ao buscar dados do dashboard", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []models.DashboardItem = []models.DashboardItem{} 

	for rows.Next() {
		var item models.DashboardItem
		err := rows.Scan(&item.ID, &item.CodExecutora, &item.Acao, &item.MetaPrevista, &item.MetaRealizada, &item.LocalExecucao, &item.AnoVigente, &item.Quadrimestre)
		if err != nil {
			continue
		}
		lista = append(lista, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

// 4. PainelConformidade gera os faróis de compliance
func PainelConformidade(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	anoStr := r.URL.Query().Get("ano")
	if anoStr == "" {
		anoStr = "2026"
	}
	quadStr := r.URL.Query().Get("quadrimestre")
	if quadStr == "" {
		quadStr = "1"
	}

	ano, _ := strconv.Atoi(anoStr)
	quadrimestre, _ := strconv.Atoi(quadStr)

	query := `
		SELECT 
			SUBSTR(p.cod_executora, 1, 4) AS cod_secretaria,
			COUNT(DISTINCT p.id) AS acoes_exigidas,
			COUNT(DISTINCT m.id) AS acoes_entregues
		FROM planejamento_loa p
		LEFT JOIN monitoramento_audesp m 
			ON p.cod_executora = m.cod_executora 
			AND p.acao = m.acao 
			AND m.ano_vigente = $1 
			AND m.quadrimestre = $2
		GROUP BY SUBSTR(p.cod_executora, 1, 4)
		ORDER BY cod_secretaria;
	`

	rows, err := database.DB.Query(query, ano, quadrimestre)
	if err != nil {
		fmt.Println("❌ Erro na consulta de compliance:", err)
		http.Error(w, "Erro ao processar dados de conformidade", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []models.ComplianceSecretaria = []models.ComplianceSecretaria{}

	for rows.Next() {
		var sec models.ComplianceSecretaria
		if err := rows.Scan(&sec.CodSecretaria, &sec.AcoesExigidas, &sec.AcoesEntregues); err == nil {
			lista = append(lista, sec)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

// 5. DetalhesSecretaria faz o drill-down mostrando o que falta entregar
func DetalhesSecretaria(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	codSecretaria := r.URL.Query().Get("secretaria")
	anoStr := r.URL.Query().Get("ano")
	quadStr := r.URL.Query().Get("quadrimestre")

	if codSecretaria == "" {
		http.Error(w, "Secretaria não informada", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			p.cod_executora,
			p.acao,
			CASE WHEN m.id IS NOT NULL THEN true ELSE false END as entregue,
			COALESCE(m.meta_prevista, 0) as meta_prevista,
			COALESCE(m.meta_realizada, 0) as meta_realizada
		FROM planejamento_loa p
		LEFT JOIN monitoramento_audesp m 
			ON p.cod_executora = m.cod_executora 
			AND p.acao = m.acao 
			AND m.ano_vigente = $2 
			AND m.quadrimestre = $3
		WHERE SUBSTR(p.cod_executora, 1, 4) = $1
		ORDER BY entregue ASC, p.cod_executora ASC;
	`

	rows, err := database.DB.Query(query, codSecretaria, anoStr, quadStr)
	if err != nil {
		fmt.Println("❌ Erro ao buscar detalhes:", err)
		http.Error(w, "Erro no servidor", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lista []models.DetalheAcao = []models.DetalheAcao{}

	for rows.Next() {
		var d models.DetalheAcao
		if err := rows.Scan(&d.CodExecutora, &d.Acao, &d.Entregue, &d.MetaPrevista, &d.MetaRealizada); err == nil {
			lista = append(lista, d)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}