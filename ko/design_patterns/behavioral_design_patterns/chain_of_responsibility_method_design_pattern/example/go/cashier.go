package main

import "fmt"

type Cashier struct {
	next Department
}

func (c *Cashier) execute(p *Patient) {
	if p.paymentDone {
		fmt.Println("Payment already processed")
		return
	}

	fmt.Println("Cashier: Processing payment for patient")
	p.paymentDone = true
}

func (c *Cashier) setNext(next Department) {
	c.next = next
}
