package utils

func Mapify[S []I, I any, C comparable](s S, c func(item I) C) map[C]I {
	m := make(map[C]I, len(s))
	for _, v := range s {
		m[c(v)] = v
	}
	return m
}
