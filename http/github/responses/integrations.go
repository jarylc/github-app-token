package responses

type Integrations []struct {
	Id      int `json:"id"`
	Account struct {
		Login string `json:"login"`
	}
}
