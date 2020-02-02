package order

import (
	"testing"
)

func TestOrder(t *testing.T) {
	s := []string{"dotisff", "我是第123g123", "lkjasd", "spoi1eq", "第12讲但是", "第讲asd12", "萨拉丁就卡死", "qweq24s", "第20讲", "asdads25"}
	t.Log(Strings(s))

	// for _, v := range s {
	// 	fmt.Println(FdNum(v))
	// }

}
