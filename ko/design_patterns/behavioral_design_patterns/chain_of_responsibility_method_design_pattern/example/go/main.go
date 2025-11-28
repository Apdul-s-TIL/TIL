package main

func main() {
	patient := &Patient{name: "Apdul"}

	cashier := &Cashier{}
	medical := &Medical{}
	doctor := &Doctor{}
	reception := &Reception{}

	reception.setNext(doctor)
	doctor.setNext(medical)
	medical.setNext(cashier)

	reception.execute(patient)
}
