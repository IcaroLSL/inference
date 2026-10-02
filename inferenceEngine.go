package main

// Definição do nível de orçamento e tipos de quarto
const (
	OrcamentoBaixo = 1
	OrcamentoMedio = 2
	OrcamentoAlto  = 3
)

const (
	QuartoStandard = "Quarto Standard"
	QuartoLuxo     = "Quarto Luxo"
	SuitePremium   = "Suíte Premium"
	QuartoFamilia  = "Quarto Família"
	Indefinido     = "indefinido"
)

// Estruturas de dados para a sessão de inferência
type facts struct {
	people        int
	budget        int
	focus         string 
	child         bool
	accommodation string
}


// Estrutura para definir regras de inferência
type Rule struct {
	ID         string
	Condicao   func(f facts) bool
	Result     string
}

// regras de inferência para determinar o tipo de quarto com base nos fatos
var rules = []Rule{
	{"R01", func(f facts) bool { return f.people == 1 && f.budget == OrcamentoBaixo }, QuartoStandard},
	{"R02", func(f facts) bool {
		return f.people == 2 && f.budget == OrcamentoBaixo && f.focus == "economia"
	}, QuartoStandard},
	{"R03", func(f facts) bool {
		return f.people == 2 && f.budget == OrcamentoMedio && f.focus == "conforto"
	}, QuartoLuxo},
	{"R04", func(f facts) bool {
		return f.people == 2 && f.budget == OrcamentoAlto && f.focus == "privacidade"
	}, SuitePremium},
	{"R05", func(f facts) bool { return f.people >= 3 && f.child }, QuartoFamilia},
	{"R06", func(f facts) bool { return f.people >= 4 && f.focus == "espaco" }, QuartoFamilia},
	{"R07", func(f facts) bool {
		return f.people == 3 && f.budget == OrcamentoAlto && f.focus == "conforto"
	}, SuitePremium},
	{"R08", func(f facts) bool { return f.people == 2 && f.accommodation == "lua de mel" }, SuitePremium},
	{"R09", func(f facts) bool {
		return f.people >= 3 && f.budget == OrcamentoMedio && f.focus == "economia"
	}, QuartoFamilia},
	{"R10", func(f facts) bool {
		return f.people == 1 && f.budget == OrcamentoAlto && f.focus == "conforto"
	}, QuartoLuxo},
	{"R11", func(f facts) bool {
		return f.people == 3 && f.budget == OrcamentoMedio && f.focus == "conforto" && !f.child && f.accommodation == "normal"
	}, QuartoLuxo},
}

// Executor de inferência
func runInference(raw InferenceFacts) (string, string) {
	f := facts{
		people:        toInt(raw["people"]),
		budget:        toInt(raw["budget"]),
		focus:         normalize(toString(raw["focus"])),
		child:         toBool(raw["child"]),
		accommodation: normalize(toString(raw["accommodation"])),
	}

	for _, r := range rules {
		if r.Condicao(f) {
			return r.Result, r.ID
		}
	}
	return Indefinido, ""
}
