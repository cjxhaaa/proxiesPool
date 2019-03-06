package main

import "fmt"

type passSS struct {
	X int
}

type passPP struct {
	X int
}


func changeSS(pass passSS) {
	pass.X = 2
}

func changePP(pass *passPP) {
	pass.X = 2
}

func main() {
	ss := passSS{X:1}
	pp := passPP{X:3}
	changeSS(ss)
	changePP(&pp)
	fmt.Println(ss.X)
	fmt.Println(pp.X)
}
