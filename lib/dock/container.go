package dock

type BuildImagePayload struct {
	Image string
}

func (bip BuildImagePayload) CreatePayload() []byte {
	return []byte()
}

func CreateBuildImageJSONPayload(imgName string) []byte {
	return []byte(
		`{{}}`
	)
}