package course

type Content struct {
	NeedLogged  bool
	Order       int
	FileUrl     []string
	Title       string
	Description string
}

func (c *Content) IsPremium() bool {
	return c.NeedLogged
}
