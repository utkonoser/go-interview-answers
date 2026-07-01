//go:build !solution

// Задача на code review: прочитать первый байт из каждого файла.
package deferinloop

import "os"

// FirstBytes возвращает первый байт каждого файла (или 0 при пустом файле).
func FirstBytes(paths []string) ([]byte, error) {
	out := make([]byte, len(paths))

	for i, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		var b [1]byte
		if _, err := f.Read(b[:]); err != nil {
			return nil, err
		}
		out[i] = b[0]
	}

	return out, nil
}
