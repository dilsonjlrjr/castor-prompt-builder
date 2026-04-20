package parser

// Role representa um papel carregado de roles/**/*.md
type Role struct {
	ID          string   `yaml:"id"`
	Nome        string   `yaml:"nome"`
	Tom         string   `yaml:"tom"`
	Habilidades []string `yaml:"habilidades"`
	GapsComuns  []string `yaml:"gaps_comuns"`
	Descricao   string   // corpo do .md abaixo do frontmatter
	Categoria   string   // derivado do nome do subdiretório
}

// FieldType define o tipo de campo de um modelo
type FieldType string

const (
	FieldText        FieldType = "text"
	FieldTextarea    FieldType = "textarea"
	FieldSelect      FieldType = "select"
	FieldMultiselect FieldType = "multiselect"
	FieldList        FieldType = "list"
	FieldSteps       FieldType = "steps"
)

// StepField define um subcampo dentro de um campo tipo steps
type StepField struct {
	ID    string    `yaml:"id"`
	Label string    `yaml:"label"`
	Tipo  FieldType `yaml:"tipo"`
}

// Field representa um campo de um modelo de prompt
type Field struct {
	ID          string      `yaml:"id"`
	Label       string      `yaml:"label"`
	Tipo        FieldType   `yaml:"tipo"`
	Obrigatorio bool        `yaml:"obrigatorio"`
	Opcoes      []string    `yaml:"opcoes"`
	StepCampos  []StepField `yaml:"step_campos"`
}

// Model representa um modelo de prompt carregado de models/*.md
type Model struct {
	ID        string  `yaml:"id"`
	Nome      string  `yaml:"nome"`
	Descricao string  `yaml:"descricao"`
	Campos    []Field `yaml:"campos"`
	Template  string  // corpo do .md abaixo do frontmatter
}

// Step representa uma fase de execução preenchida pelo usuário
type Step struct {
	Titulo    string
	Descricao string
}
