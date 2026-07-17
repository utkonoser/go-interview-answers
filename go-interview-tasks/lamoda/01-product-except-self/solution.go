// Моё решение с собеса Lamoda.
package productexceptself

func MultOther(in []int) []int {
	allMult := 1
	countZero := 0
	for _, item := range in {
		if item == 0 {
			countZero++
			continue
		}

		allMult = allMult * item
	}

	for i := range in {
		if countZero > 1 {
			in[i] = 0
			continue
		}

		if countZero == 1 && in[i] != 0 {
			in[i] = 0
			continue
		}

		if in[i] == 0 {
			in[i] = allMult
			continue
		}

		in[i] = allMult / in[i]
	}

	return in
}
