package orders

import "math/rand"

func GenerateOrderId() int32 {
	orderId := rand.Intn(999) + 1

	return int32(orderId)
}
