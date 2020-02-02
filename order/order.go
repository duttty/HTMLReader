package order

import (
	"strconv"
)

func Strings(s []string) []string {
	return classify(s)
}

func classify(s []string) []string {
	l := len(s)
	//可排序部分
	//值切片
	val := make([]string, l, l)
	//数切片
	num := make([]int, l, l)
	lNum := 0

	//不可排序部分
	no := make([]string, 0, l)
	for _, v := range s {
		o, _ := FdNum(v)
		if o == -1 {
			no = append(no, v)
		} else {
			i := lNum
			for ; i > 0; i-- {
				if o < num[i-1] {
					//右移
					val[i] = val[i-1]
					num[i] = num[i-1]
				} else {
					break

				}
			}
			//插入
			num[i] = o
			val[i] = v
			lNum++
		}
	}
	return append(val[:lNum], no...)
}

func FdNum(s string) (o, fIdx int) {
	str := ""
	fIdx = -1
	l := len(s) - 1
	for k, v := range s {
		//阿拉伯数字
		if v >= '0' && v <= '9' {
			if fIdx == -1 {
				fIdx = k
			}
			if k == l {
				str = s[fIdx:]
			}
		} else {
			if fIdx != -1 {
				str = s[fIdx:k]
				break
			}
		}
	}
	o, err := strconv.Atoi(str)
	if err != nil {
		return -1, fIdx
	}
	return
}
