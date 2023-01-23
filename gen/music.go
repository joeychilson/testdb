package gen

type Artist struct {
	Name  string `fake:"{name}"`
	Image string `fake:"{imageurl}"`
}

type Album struct {
	Name  string `fake:"{phrase}"`
	Cover string `fake:"{imageurl}"`
}

type Song struct {
	Title  string  `fake:"{phrase}"`
	Length float64 `fake:"{float64}"`
	Path   string  `fake:"{url}"`
	Mtime  int32   `fake:"{uint32}"`
}
