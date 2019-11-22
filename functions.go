package main

func main() {
	// for example
	// array:=[]string{`a`,`b`,`c`,`d`,`a`,`d`,`e`,`f`,`g`,`c`,`e`}
}

/* 删除数组中某一个元素 */
func del(arr *[]string, s string) {
	ar := *arr
	for i, v := range ar {
		if v != s {
			continue
		}
		ar = append(ar[:i], ar[i+1:]...)
		break
	}
	*arr = ar
}

// 删除数组中相同的元素
func delSame(arr *[]string) {
	ar := *arr
a:
	for {
		for i, n := range ar {
			if i != len(ar)-1 {
				for j, m := range ar[i+1:] {
					if n != m {
						continue
					} else {
						ar = append(ar[:j+i+1], ar[j+i+2:]...)
						continue a
					}
				}
			}
		}
		break
	}
	*arr = ar
}

// 插入新元素，若数组中已存在，则忽略
func update(arr *[]string, s string) {
	if len(*arr) != 0 {
		for i, v := range *arr {
			if i != len(*arr)-1 {
				if v != s {
					continue
				}
				break
			}
			if v != s {
				*arr = append(*arr, s)
			}
		}
		return
	}
	*arr = append(*arr, s)
}

// 删除相邻相同的元素
func delNearby(arr *[]string) {
	ar := *arr
a:
	for {
		for i, v := range ar {
			if i != len(ar)-1 {
				if v == ar[i+1] {
					ar = append(ar[:i], ar[i+1:]...)
					continue a
				}
			}
		}
		break
	}
	*arr = ar
}
