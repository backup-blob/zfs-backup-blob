package config

type RemoteTrimPolicy string

func (e *RemoteTrimPolicy) GetIncrementalCount() int {
	return e.countChar('I')
}

func (e *RemoteTrimPolicy) GetFullCount() int {
	return e.countChar('F')
}

func (e *RemoteTrimPolicy) countChar(c rune) int {
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
