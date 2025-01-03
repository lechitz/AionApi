package domain

type LoginDomain struct {
	ID       uint64 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
