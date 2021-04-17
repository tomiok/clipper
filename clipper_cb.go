package clipper

type Status int

type Clipper struct {
	Name       string
	MaxFailure int
	State      Status
}

type Settings struct {

}

func NewClipper(s Settings) *Clipper {
	return &Clipper{
		Name:       "",
		MaxFailure: 0,
		State:      0,
	}
}

func Do(name string, fn func() error) {

}
