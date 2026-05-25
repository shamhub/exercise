package envconfigreader

import "os"

func (*envConfig) Get(key string) string {
	return os.Getenv(key)
}
