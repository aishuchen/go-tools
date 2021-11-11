// Application which greets you.
package rabbitmq

import "fmt"

func main() {
	fmt.Println(greet())
}

func greet() string {
	return "Hi!"
}
