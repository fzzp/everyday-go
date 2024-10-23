package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// 质数：又称素数，指在大于1的自然数中，除了1和该数自身外，5无法被其他自然数整除的数。

// func main() {
// 	n := 9
// 	_, msg := isPrime(n)
// 	fmt.Println(msg)

// }

func main() {
	intro()
	doneChan := make(chan bool)
	go readUserInput(doneChan)

	<-doneChan

	close(doneChan)

	fmt.Println("Goodbye.")
}

func readUserInput(doneChan chan bool) {
	scaner := bufio.NewScanner(os.Stdin)
	for {
		res, done := checkNumbers(scaner)
		if done {
			doneChan <- done
			return
		}
		fmt.Println(res)
		prompt()
	}
}

func checkNumbers(scaner *bufio.Scanner) (string, bool) {
	scaner.Scan()
	if strings.EqualFold(scaner.Text(), "q") {
		return "", true
	}
	num, err := strconv.Atoi(scaner.Text())
	if err != nil {
		return "Please enter a whole number", false
	}
	_, msg := isPrime(num)
	return msg, false
}

func intro() {
	fmt.Sprintln("Is it Prime?")
	fmt.Println("--------------")
	fmt.Println("Enter a whole number, and we'll tell you if is a prime or not. Enter q to quit.")
	prompt()
}

func prompt() {
	fmt.Print("->")
}

// isPrime 判断一个数是否时质数
func isPrime(n int) (bool, string) {
	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not prime, by definition!", n)
	}
	if n < 0 {
		return false, "Negative numbers are not prime, by definition!"
	}

	// 任意一个合数最小因子数是2，从2开始
	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false, fmt.Sprintf("%d is not a prime number it is divisible by %d", n, i)
		}
	}

	return true, fmt.Sprintf("%d is a prime number!", n)
}
