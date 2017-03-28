package share

import (
	"os"
)

func InitConfig(dir, list string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(list, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}
