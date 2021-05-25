package controllers

type DockerConfigSecret struct {
	Auths map[string]Credentials `json:"auths"`
}

type Credentials struct {
	Auth string `json:"auth"`
}
