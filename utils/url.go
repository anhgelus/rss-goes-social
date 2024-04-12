package utils

func GetFullUrl(base string, uri string) string {
	if base[len(base)-1] == '/' {
		if uri[0] == '/' {
			return base + uri[1:]
		}
		return base + uri
	}
	if uri[0] == '/' {
		return base + uri
	}
	return base + "/" + uri
}
