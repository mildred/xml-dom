package xmldom

type characterData struct {
	characterData string
}

func (c *characterData) Data() string {
	return c.characterData
}

func (c *characterData) SetData(s string) {
	c.characterData = s
}

func (c *characterData) Length() uint {
	return uint(len(c.characterData))
}
