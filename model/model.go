package model

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type AIRequest struct {
	Inputs Inputs `json:"inputs"`
}

type TapasResponse struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

type ChatResponse struct {
	GeneratedText string `json:"generated_text"`
}
