package extract

type Reader interface {
	ReadEntry() string
}
