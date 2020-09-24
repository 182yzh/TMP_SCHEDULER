package flowscheduler

type SelectorType = uint8 

const (
	IN_SET = 0;
	NOT_IN_SET = 1;
	EXISTS_KEY = 2;
	NOT_EXISTS_KEY = 3;
)


type LabelSelector struct{	
	lstype SelectorType
	key string
	values map[string]bool
}



func NewLabelSelector()*LabelSelector{
	ans := new(LabelSelector)
	ans.values = make(map[string]bool)
	return ans
}

func (ls *LabelSelector)SetLabelSelector(key string ,values []string,lstype SelectorType){
	ls.key = key
	ls.lstype = lstype
	for _,v := range values {
		ls.values[v]=true
	}
}
//check the labels can satisfy the selector
func (ls *LabelSelector)SatisfiesLabels(labels []Label) bool {
	labelsMap := make(map[string]string)
	for _,v := range labels {
		labelsMap[v.key]=v.value
	}
	return ls.SatisfiesLabelsUseMap(labelsMap)
}

//labelsMap is the map[Labelkey]LabelsValues
// labels 是某个节点上的labels，而selector是验证是否满足
func (ls *LabelSelector)SatisfiesLabelsUseMap(labelsMap map[string]string) bool{
	switch  ls.lstype {
	case IN_SET:
		return ls.SatisfiesLabelsInSet(labelsMap)
	case NOT_IN_SET:
		return ls.SatisfiesLabelsNotInSet(labelsMap)
	case EXISTS_KEY:
		return ls.SatisfiesLabelsExistsKey(labelsMap)
	case NOT_EXISTS_KEY:
		return ls.SatisfiesLabelsNotExistsKey(labelsMap)
	}
	return false
}

func (ls *LabelSelector)SatisfiesLabelsExistsKey(labelsMap map[string]string) bool{
	_,v := labelsMap[ls.key]
	return v
}

func (ls *LabelSelector)SatisfiesLabelsNotExistsKey(labelsMap map[string]string) bool{
	_,v := labelsMap[ls.key]
	return !v
}

func (ls *LabelSelector)SatisfiesLabelsInSet(labelsMap map[string]string) bool {
	v,ok := labelsMap[ls.key]
	if !ok {
		return false
	}
	_,ok = ls.values[v]
	return ok
}

func (ls *LabelSelector)SatisfiesLabelsNotInSet(labelsMap map[string]string)bool {
	v,ok := labelsMap[ls.key]
	if !ok {
		return true
	}
	_,ok = ls.values[v]
	return !ok
}
