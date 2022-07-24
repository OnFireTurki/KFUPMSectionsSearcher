package main

import (
	LFC "KFUPMSS/src"
)

func main() {

	LFC.InitClient(0)
	w := &LFC.Worker{Term: "202210", CourseSub: "MATH", CourseNume: "101", OptionSe: LFC.Options{Name: "", WaitListMode: false}}
	d := w.SingleSearch()
	for i := 0; i < len(d); i++ {
		d[i].PrintIT()
	}
}
