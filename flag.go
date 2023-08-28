package cdb

type Flag func() []string

func Username(name string) Flag {
	return func() []string {
		return []string{"--username", name}
	}
}

func Password(pass string) Flag {
	return func() []string {
		return []string{"--password", pass}
	}
}

func AsOpf() []string {
	return []string{"--as-opf"}
}
