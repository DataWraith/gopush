package gopush

// RandomCode implements the standard Push random code generation algorithm described
// at http://faculty.hampshire.edu/lspector/push3-description.html#RandomCode
func (i *Interpreter) RandomCode(maxPoints int64) Code {
	return i.randomCodeWithSize(1 + i.Rand.Int63n(maxPoints))
}

func (i *Interpreter) decompose(number int64, maxParts int64) []int64 {
	if number == 1 || maxParts == 1 {
		return []int64{number}
	}

	thisPart := 1 + i.Rand.Int63n(number-1)

	return append([]int64{thisPart}, i.decompose(number-thisPart, maxParts-1)...)
}

func (i *Interpreter) randomCodeWithSize(numPoints int64) Code {
	if numPoints == 1 {
		return i.randomInstruction()
	}

	sizesThisLevel := i.decompose(numPoints-1, numPoints-1)
	codeFragments := make([]Code, 0, len(sizesThisLevel))

	for _, v := range sizesThisLevel {
		codeFragments = append(codeFragments, i.randomCodeWithSize(v))
	}

	// Shuffle the codeFragments slice
	for j := range codeFragments {
		k := i.Rand.Intn(j + 1)
		codeFragments[j], codeFragments[k] = codeFragments[k], codeFragments[j]
	}

	c := Code{
		Length: 0,
		List:   codeFragments,
	}

	for _, fragment := range codeFragments {
		c.Length += fragment.Length
	}

	return c
}
