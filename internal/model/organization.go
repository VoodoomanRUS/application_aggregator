package model

// Organization — доменная модель организации
type Organization struct {
	UUID  string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
