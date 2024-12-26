package render

func ConvertSlice[From any, To any](to []To, from []From, conv func(From) To) {
	for i, v := range from {
		if i >= len(to) {
			return
		}

		to[i] = conv(v)
	}
}

func ConvertSliceWithAlloc[From any, To any](from []From, conv func(From) To) []To {
	to := make([]To, len(from))
	for i, v := range from {
		to[i] = conv(v)
	}

	return to
}
