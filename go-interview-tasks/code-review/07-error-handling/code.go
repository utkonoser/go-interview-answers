//go:build !solution

// Задача на code review: сохранить данные в файл.
package errorhandling

import "os"

// Save записывает data в path.
func Save(path string, data []byte) error {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return nil
}
