package domain

type Menu struct {
	ID         int    `json:"id,omitempty"`
	Icon       string `json:"icon,omitempty"`
	RouterLink string `json:"router_link,omitempty"`
	Text       string `json:"text,omitempty"`
}

type Children struct {
	ID         int    `json:"id,omitempty"`
	Text       string `json:"text,omitempty"`
	Icon       string `json:"icon,omitempty"`
	RouterLink string `json:"router_link,omitempty"`
}
