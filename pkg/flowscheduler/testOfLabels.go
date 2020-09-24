package flowscheduler

import "fmt"

func LabelsTest(){
	ls := NewLabelSelector()
	values := make([]string,0,1)
	values = append(values,"11")
	values = append(values,"33")
	values = append(values,"55")
	values = append(values,"77")
	values = append(values,"99")
	key:="num"
	ls.SetLabelSelector(key,values,IN_SET)
	fmt.Println(ls.key,ls.lstype)
	fmt.Println(ls.values)


	labels := make([]Label,0,1)
	l := Label{
		key:"num",
		value:"33",
	}
	labels = append(labels,l)
	l = Label{
		key:"alpha",
		value:"aa",
	}
	labels = append(labels,l)
	l = Label{
		key:"method",
		value:"add",
	}
	labels = append(labels,l)
	l = Label{
		key:"judge",
		value:"false",
	}
	labels = append(labels,l)
	fmt.Println(ls.SatisfiesLabels(labels))
	

	ls.lstype = NOT_IN_SET
	fmt.Println(ls.SatisfiesLabels(labels))

	ls.lstype = EXISTS_KEY
	fmt.Println(ls.SatisfiesLabels(labels))

	ls.lstype = NOT_EXISTS_KEY
	fmt.Println(ls.SatisfiesLabels(labels))
}
