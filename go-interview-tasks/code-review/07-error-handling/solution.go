//go:build solution

package errorhandling

import "os"

// fix: defer f.Close() до проверки err → паника на nil *os.File при ошибке Create.
// fix: ошибка Write игнорируется — возвращаем nil при частичной записи.
func Save(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}
