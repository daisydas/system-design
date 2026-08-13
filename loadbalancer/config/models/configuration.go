package models

type Configuration struct {
	CurrentAPIServers []string `json:"api_servers"`
	Algorithm         string   `json:"algorithm"`
	UnhealthyServers  []string `json:"unhealthy_servers"`
}
