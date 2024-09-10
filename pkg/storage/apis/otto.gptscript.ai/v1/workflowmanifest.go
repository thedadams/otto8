package v1

type WorkflowManifest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
	Steps       []Step            `json:"steps,omitempty"`
}

type Step struct {
	*AgentStep
	*ToolStep
	Name    string    `json:"name,omitempty"`
	Input   StepInput `json:"input,omitempty"`
	Tool    string    `json:"tool,omitempty"`
	If      *If       `json:"if,omitempty"`
	While   *While    `json:"while,omitempty"`
	ForEach *ForEach  `json:"forEach,omitempty"`
}

type AgentStep struct {
	Prompt string   `json:"prompt,omitempty"`
	Tools  []string `json:"tools,omitempty"`
}

type ToolStep struct {
	Tool     string `json:"tool,omitempty"`
	Metadata map[string]string
}

type StepInput struct {
	Content string            `json:"content,omitempty"`
	Args    map[string]string `json:"args,omitempty"`
}

type If struct {
	Condition string `json:"condition,omitempty"`
	Steps     []Step `json:"steps,omitempty"`
	Else      []Step `json:"else,omitempty"`
}

type While struct {
	Condition string `json:"condition,omitempty"`
	MaxLoops  int    `json:"maxLoops,omitempty"`
	Steps     []Step `json:"steps,omitempty"`
}

type ForEach struct {
	Items string `json:"items,omitempty"`
	Var   string `json:"var,omitempty"`
	Steps []Step `json:"steps,omitempty"`
}