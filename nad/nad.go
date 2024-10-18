package nad

import "embed"

//go:embed lib
var nadlib embed.FS

func ReadFiles(fn func(string, []byte) error) error {
	ls, err := nadlib.ReadDir("lib")

	if err != nil {
		return err
	}

	for _, e := range ls {
		if e.IsDir() {
			continue
		}

		b, bErr := nadlib.ReadFile("lib/" + e.Name())

		if bErr != nil {
			return bErr
		}

		fnErr := fn(e.Name(), b)

		if fnErr != nil {
			return fnErr
		}
	}

	return nil
}
