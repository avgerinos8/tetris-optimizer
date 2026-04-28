package checks

import "errors"

func Checks(s string) Tetromino {
	if s != "" {
		s, err := verifyLines(s)
		if err != nil {
			// ERROR INVALID TETRONOMINO
			continue
		}
		data, err := to2DSlice(s)
		if err != nil {
			// INTERNAL ERROR could not convert to slice
			continue
		}
		return nil
	}
}

func verifyLines(s string) (string, error) {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '.' && s[i] != '#' && s[i] != '\n' {
			return "", errors.New("invalid tetromino")
		}
		if s[i] == '\n' {
			count++
			if count > 3 {
				return "", errors.New("invalid tetromino")
			}
		}
	}
	return s, nil
}

func to2DSlice(s string) ([][]int, error) {
	result := [][]int{}
	index := 0
	i := 0

	for {
		// 1. Check if we have reached the end of the string
		if index >= len(s) {
			break
		}

		// 2. Initialize a new row in the 2D slice
		result = append(result, []int{})

		for {
			// 3. Check for line breaks or end of string to terminate the current row
			if index >= len(s) || s[index] == '\n' {
				index++ // Skip the newline character
				break
			}

			// 4. Map characters to integers: '#' becomes 1, everything else becomes 0
			val := 0
			if s[index] == '#' {
				val = 1
			}

			// 5. Append the value to the current row and advance index
			result[i] = append(result[i], val)
			index++
		}

		// 6. Increment row counter
		i++
	}

	return result, nil
}
