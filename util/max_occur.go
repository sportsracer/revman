package util

// Return the item that occurs most frequently in iter
func MaxOccur[T comparable](iter Iterable[T]) (val *T, count uint) {
	var max uint
	var maxItem *T
	counter := make(map[T]uint)
	for item := range iter {
		count, ok := counter[item]
		if ok {
			count += 1
		} else {
			count = 1
		}
		if count > max {
			// copy the item so we don't give out a pointer to the loop variable
			itemCopy := item
			max, maxItem = count, &itemCopy
		}
		counter[item] = count
	}
	return maxItem, max
}
