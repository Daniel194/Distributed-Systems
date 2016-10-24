package main

import "fmt"


//Fibonacci Numbers.
func fib() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a + b
		return a
	}
}

func main() {
	var k int

	fmt.Print("Give a number : ")
	n, err := fmt.Scanf("%d", &k)

	if err != nil || n != 1 {
		// handle invalid input
		fmt.Println("Invalid input : ", n, err)
	}

	f := fib()

	for i := 0; i < k; i++ {
		fmt.Println(f())
	}

}
