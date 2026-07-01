//go:build solution

package deferinloop

import "os"

// fix: defer в цикле — все Close() выполнятся только при выходе из функции.
// Держим открытыми все файлы до конца → исчерпание fd, неверное чтение.
func FirstBytes(paths []string) ([]byte, error) {
	out := make([]byte, len(paths))

	for i, path := range paths {
		b, err := readFirstByte(path)
		if err != nil {
			return nil, err
		}
		out[i] = b
	}

	return out, nil
}

func readFirstByte(path string) (byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var b [1]byte
	if _, err := f.Read(b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}
