package algo

import (
	"fmt"
	"testing"
)

func TestGetRandomMoney(t *testing.T) {

	moneyArr := make([]int64, 0)
	for ;TotalAmountOfEnvelope > 0; {
		x := GetRandomMoney()
		moneyArr = append(moneyArr, x)
	}

	fmt.Println("分配的红包金额:")
	fmt.Println(moneyArr)

	fmt.Println("分配的红包金额总和:")
	var sum int64 = 0
	for _, num := range moneyArr {
		sum += num
	}
	fmt.Println(sum)
}
