package main

import "fmt"

// CMMDC
func main() {

	var n int
	var num int
	var result int

	fmt.Println("How many integers do you want to calculate?")

	fmt.Scanf("%d\n", &n)

	fmt.Print("Enter space separated integers:")

	numList := map[int]int{}

	for i := 1; i <= n; i++ {
		fmt.Scanf("%d", &num)
		numList[i] = num
	}

	result = gcd(numList[1], numList[2])

	for j := 3; j <= n; j++ {
		result = gcd(result, numList[j])
	}

	fmt.Println("The result of given integers is: ", result)
}

//Func to implement Euclid Algorithm
func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x % y
	}
	return x
}
