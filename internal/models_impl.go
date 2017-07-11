package gaefire

type WebApplicationInfoModel struct {
	ProjectId    string `json:"project_id"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type WebApplicationModel struct {
	Web WebApplicationInfoModel `json:web`
}
