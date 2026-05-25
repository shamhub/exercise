package envconfigreader

type IEnvconfigReader interface {
	Get(key string) (value string)
}
