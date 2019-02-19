package util

import (
	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 10

//生成哈希密码
func GeneratePasswordHash(input []byte, costs ...int) ([]byte, error) {
	//默认哈希成本（算法使用的cost:循环次数以 2 为底的对数）
	cost := DefaultCost
	if len(costs) > 0 {
		cost = costs[0]
	}
	//生成
	return bcrypt.GenerateFromPassword(input, cost)
}

//验证密码
func ValidatePassword(input, hash []byte) error {
	//比对，相同时返回nil
	return bcrypt.CompareHashAndPassword(hash, input)
}
