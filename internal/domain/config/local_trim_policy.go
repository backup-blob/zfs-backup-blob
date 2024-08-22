package config

type LocalTrimPolicy string

func (e *LocalTrimPolicy) GetIncrementalCount() int {
	return e.countChar('I')
}

func (e *LocalTrimPolicy) GetFullCount() int {
	return e.countChar('F')
}

func (e *LocalTrimPolicy) countChar(c rune) int {
	policy := []rune(*e)

	var counter int

	for i := 0; i < len(policy); i++ {
		char := policy[i]
		if char == c {
			counter++
		}
	}

	return counter
}
