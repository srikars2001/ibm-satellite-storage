package mustachejsonvalidation

type ParamKeys struct {
	prefix string
	value  string
	suffix string
}

type Stack []ParamKeys

// IsEmpty checks whether stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// push ParamKey struct into stack
func (s *Stack) Push(pk ParamKeys) {
	*s = append(*s, pk)
}

// Peek top element of stack if stack is empty returns empty struct,false
// else returns tje paramkeys struct,true
func (s *Stack) Peek() (ParamKeys, bool) {
	if s.IsEmpty() {
		return ParamKeys{}, false
	}

	index := len(*s) - 1
	element := (*s)[index]
	return element, true
}

// Remove top element of stack and return true, else return nil,false
func (s *Stack) Pop() (ParamKeys, bool) {
	if s.IsEmpty() {
		return ParamKeys{}, false
	}

	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element, true
}
