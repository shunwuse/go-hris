package dtos

type HealthResponse struct {
	Status     string                   `json:"status"`
	Components HealthComponentsResponse `json:"components"`
	Info       HealthInfoResponse       `json:"info"`
}

type HealthComponentsResponse struct {
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

type HealthInfoResponse struct {
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Uptime      string `json:"uptime"`
	InstanceID  string `json:"instance_id"`
	Hostname    string `json:"hostname"`
	GoVersion   string `json:"go_version"`
}
