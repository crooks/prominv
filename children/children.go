package children

import (
	"errors"
	"sort"
)

var (
	errChildNotFound error = errors.New("child Not Found")
	errChildExists   error = errors.New("child already exists")
)

type Children map[string]Child

type Child struct {
	Name    string
	Members map[string]int
}

func NewChildren() Children {
	var c Children = make(map[string]Child)
	return c
}

func (children Children) AddChild(childName string) error {
	_, ok := children[childName]
	if ok {
		return errChildExists
	}
	child := new(Child)
	child.Name = childName
	child.Members = make(map[string]int)
	children[childName] = *child
	return nil
}

func (children Children) GetChild(childName string) (Child, error) {
	child, ok := children[childName]
	if !ok {
		return child, errChildNotFound
	}
	return child, nil
}

// getAllChildren returns a slice of all the keys in the Children map
func (children Children) GetAllChildren() (c []string) {
	for k := range children {
		c = append(c, k)
	}
	sort.Strings(c)
	return
}

func (children Children) MemberSlice(childName string) (members []string, err error) {
	child, ok := children[childName]
	if !ok {
		err = errChildNotFound
		return
	}
	for k := range child.Members {
		members = append(members, k)
	}
	sort.Strings(members)
	return
}

// AddMember adds a host to a specific Child map.
// There is no test for duplication in this method as it does no harm to attempt to add a host more than once.
func (children Children) AddMember(childName, hostname string) (err error) {
	child, ok := children[childName]
	if !ok {
		err = errChildNotFound
		return
	}
	child.Members[hostname] = 1
	return
}
