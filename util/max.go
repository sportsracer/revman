package util

// Return the item that occurs most frequently in iter
func MaxOccur(iter Iterable) (val interface{}, count int) {

	max := 0
	var maxItem interface{}
	counter := make(map[interface{}]int)
	for item := range iter.Iter() {
		count, ok := counter[item]
		if ok {
			count += 1
		} else {
			count = 1
		}
		if count > max {
			max, maxItem = count, item
		}
		counter[item] = count
	}

	return maxItem, max
}
