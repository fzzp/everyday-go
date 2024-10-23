package main

import (
	"fmt"
	"testing"
)

/*
	在根目录下：
	go test .      # 执行测试
	go test ./...  # 执行包下所有测试
 	go test -v .   # 执行测试，输出明细
	go test -v . -run 测试函数名字， # 如：go test -v -run Test_isPrime2
	go test -cover . # 查看测试覆盖率
	go test -v -run Test_alpha # 执行Test_alpha开头的测试（测试分组）

	使用浏览器查看测试结果
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
*/

func Test_isPrime(t *testing.T) {
	reuslt, msg := isPrime(0)
	// 0 不是质数，reuslt = true 是错误的。
	// 写测试，一般思维是结果不等于正确的就是错误的，测试不通过；另一个思维就是结果符合的错误值，也是测试不通过。
	if reuslt {
		t.Errorf("with %d test parameter, got true, but expected false.", 0)
	}
	if msg != "0 is not prime, by definition!" {
		t.Error("wrong message returned: ", msg)
	}

	reuslt, msg = isPrime(7)
	if !reuslt {
		t.Errorf("with 7 as test parameter, got false, but expected true")
	}

	if msg != fmt.Sprintf("%d is a prime number!", 7) {
		t.Error("wrong message returned: ", msg)
	}
}

// 表格测试
func Test_isPrime2(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not prime", 9, false, fmt.Sprintf("%d is not a prime number it is divisible by %d", 9, 3)},
		{"zero prime", 0, false, fmt.Sprintf("%d is not prime, by definition!", 0)},
		{"负数", -10, false, fmt.Sprint("Negative numbers are not prime, by definition!")},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func Test_alpha_Demo(t *testing.T) {
	t.Log("测试组")
}

func Test_alpha_Example(t *testing.T) {
	t.Log("测试组")
}
