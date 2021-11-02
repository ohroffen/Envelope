package algo

import (
	"math/rand"
	"time"
)

var TotalAmountOfMoney int64 = 100 * 100 // 红包雨总预算
var TotalAmountOfEnvelope int64 = 30     // 红包总个数

const (
	SnatchRatio    float64 = 0.8       // 多大概率抢到红包，[0, 1.0)之间
	MaxSnatchCount int     = 5         // 每个人能抢的红包总个数
	MaxAmount      int64   = 100 * 100 // 单个红包最大金额
	MinAmount      int64   = 50        // 单个红包最小金额
)

// GetRandomMoney
// 调用者负责维护剩余红包数量与剩余总金额
func GetRandomMoney() int64 {
	money := int64(0)
	if TotalAmountOfMoney == 0 || TotalAmountOfEnvelope == 0 {
		return money
	}

	if TotalAmountOfEnvelope == 1 {
		money = TotalAmountOfMoney
		TotalAmountOfMoney -= money
		TotalAmountOfEnvelope--
		return money
	}

	// 最大可调度金额
	max := TotalAmountOfMoney - MinAmount*TotalAmountOfEnvelope
	if max <= 0 {
		return 0
	}
	// fmt.Printf("剩余总金额：%d，剩余总个数：%d\n", TotalAmountOfMoney, TotalAmountOfEnvelope)

	// fmt.Printf("剩余可调度金额：%d\n", max)

	// 每个红包平均调度金额
	avgMax := max / TotalAmountOfEnvelope
	// fmt.Printf("每个红包平均调度金额：%d\n", avgMax)

	// 根据平均调度金额来生成每个红包金额
	Init()
	randNum := rand.Float64() - 0.5
	// fmt.Printf("浮动比率：%f\n", randNum)

	avgMax += int64(randNum * float64(avgMax))
	// fmt.Printf("最终浮动金额：%d\n", avgMax)

	money = MinAmount + avgMax
	// fmt.Printf("最终金额：%d\n", money)
	if money < MinAmount {
		money = MinAmount
	}

	if money > MaxAmount {
		money = MaxAmount
	}

	TotalAmountOfMoney -= money
	TotalAmountOfEnvelope--

	return money
	//if remainSize == 1 {
	//	//remainSize--
	//	return float64(math.Round(remainMoney * 100)) / 100
	//}
	//
	////r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//rand.Seed(time.Now().Unix())
	//min := 0.01
	//max := remainMoney / float64(remainSize) * 2
	//
	//money := rand.Float64() * max
	//if money < min {
	//	money = 0.01
	//}
	////money += min
	//money = math.Floor(money * 100) / 100
	////remainSize--
	////remainMoney -= money
	//
	//return money
}

func Init() {
	// 初始化随机数的资源库，如果不执行这行，不管运行多少次都返回同样的值
	rand.Seed(time.Now().UnixNano())
}
