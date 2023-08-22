package docker

func ShortId(id string) string {
	return id[:shortIdLength]
}
