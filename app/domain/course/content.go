package course

type Content struct {
	NeedLogged  bool     `json:"need_logged"`
	Order       int      `json:"order"`
	FileUrl     []string `json:"file_url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
}

func (c *Content) IsPremium() bool {
	return c.NeedLogged
}
