package domains

type Health struct {
	Status     string
	Components HealthComponents
	Info       HealthInfo
}

type HealthComponents struct {
	Database string
}

type HealthInfo struct {
	Version     string
	Environment string
	Uptime      string
	InstanceID  string
	Hostname    string
	GoVersion   string
}
