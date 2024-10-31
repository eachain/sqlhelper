package sqlhelper

func Escape(s string) string {
	n, ok := checkStringValid(s)
	if ok {
		return s
	}
	return escapeStringBackslash(n, s)
}

func checkStringValid(v string) (n int, valid bool) {
	valid = true
	for i := 0; i < len(v); i++ {
		c := v[i]
		switch c {
		case '\x00', '\n', '\r', '\x1a', '\'', '"', '\\':
			valid = false
			n += 2
		default:
			n++
		}
	}
	return
}

func escapeStringBackslash(n int, v string) string {
	buf := make([]byte, n)
	pos := 0
	for i := 0; i < len(v); i++ {
		c := v[i]
		switch c {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		default:
			buf[pos] = c
			pos++
		}
	}

	return string(buf)
}
