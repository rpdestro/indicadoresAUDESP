package models

// Monitoramento representa a estrutura exata do nosso relatório da AUDESP
type Monitoramento struct {
	CodExecutora  string  `json:"und_executora"`
	Programa      string  `json:"programa"`
	Acao          string  `json:"acao"`
	UnidadeMedida string  `json:"und_medida"`
	AnoVigente    int     `json:"ano_vigente"`
	Quadrimestre  int     `json:"quadrimestre"`
	MetaPrevista  float64 `json:"meta_prevista"`
	MetaRealizada float64 `json:"meta_realizada"`
	LocalExecucao string  `json:"local_execucao"`
	DataInicio    string  `json:"data_inicio"`
	DataFim       string  `json:"data_fim"`
	Justificativa string  `json:"justificativa"`
}

// PlanejamentoLOA representa os dados de devolução automática para a tela
type PlanejamentoLOA struct {
	Programa      string `json:"programa"`
	Acao          string `json:"acao"`
	UnidadeMedida string `json:"und_medida"`
}

// DashboardItem representa os dados resumidos para o painel de consulta visual
type DashboardItem struct {
	ID            string     `json:"id"`
	CodExecutora  string  `json:"cod_executora"`
	Acao          string  `json:"acao"`
	MetaPrevista  float64 `json:"meta_prevista"`
	MetaRealizada float64 `json:"meta_realizada"`
	LocalExecucao string  `json:"local_execucao"`
	AnoVigente    int     `json:"ano_vigente"`
	Quadrimestre  int     `json:"quadrimestre"`
}

// ComplianceSecretaria representa o status de entrega de cada secretaria
type ComplianceSecretaria struct {
	CodSecretaria  string `json:"cod_secretaria"`
	AcoesExigidas  int    `json:"acoes_exigidas"`
	AcoesEntregues int    `json:"acoes_entregues"`
}

// DetalheAcao representa o status individual de uma ação na LOA para o Drill-Down
type DetalheAcao struct {
	CodExecutora  string  `json:"cod_executora"`
	Acao          string  `json:"acao"`
	Entregue      bool    `json:"entregue"`
	MetaPrevista  float64 `json:"meta_prevista"`
	MetaRealizada float64 `json:"meta_realizada"`
}
