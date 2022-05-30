package ch9

var deposits = make(chan int)
var balances = make(chan int)
var withdraw = make(chan withdrawal)

type withdrawal struct {
	amount  int
	success chan bool
}

func Deposit(amount int) { deposits <- amount }

func Balance() int { return <-balances }

func Withdraw(amount int) bool {
	success := make(chan bool)
	withdraw <- withdrawal{amount, success}
	return <-success
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case draw := <-withdraw:
			if draw.amount >= balance {
				balance -= draw.amount
				draw.success <- true
			} else {
				draw.success <- false
			}
		}
	}
}

func init() {
	go teller() // start monitor goroutine
}
