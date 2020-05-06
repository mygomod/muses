package utils

// InSlice checks given string in string slice or not.
func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// a [1,2] b [2,3]       return [1]
func ArrayDiffInt(a, b []int) []int {
	bMap := make(map[int]struct{}, 0)
	for _, id := range b {
		bMap[id] = struct{}{}
	}

	output := make([]int, 0)
	for _, id := range a {
		if _, flag := bMap[id]; !flag {
			output = append(output, id)
		}
	}
	return output
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
