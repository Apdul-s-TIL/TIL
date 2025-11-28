package main

import "fmt"

type Doctor struct {
	next Department
}

func (d *Doctor) execute(p *Patient) {
	if p.dockerCheckUpDone {
		fmt.Println("Doctor check-up already done")
		d.next.execute(p)
		return
	}

	fmt.Println("Doctor: checking up patient")
	p.dockerCheckUpDone = true
	d.next.execute(p)
}

func (d *Doctor) setNext(next Department) {
	d.next = next
}
