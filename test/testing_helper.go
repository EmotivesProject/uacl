package test

func CreateStringAtLength(length int) string {
	var str string
	for i := 0; i < length; i++ {
		str += "x"
	}

	return str
}
