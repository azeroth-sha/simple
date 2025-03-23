package guid

type GUID []byte

func (g GUID) String() string {
	if g == nil || len(g) == 0 {
		return ""
	}
	return string(g)
}

func (g GUID) Bytes() []byte {
	return g
}

// TODO
