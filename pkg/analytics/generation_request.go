package analytics

type GenerationRequest struct {
	Duration        int      `json:"duration"`
	GameType        string   `json:"gameType"`
	GenerationCount int      `json:"generationCount"`
	Host            string   `json:"host"`
	ID              int      `json:"id"`
	MapType         string   `json:"mapType"`
	Parameters      []string `json:"parameters"`
	RequestID       string   `json:"requestId"`
	UserAgent       string   `json:"userAgent"`
}

//     Timestamp       []int    `json:"timestamp"`
