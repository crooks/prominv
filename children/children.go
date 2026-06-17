package children

import (
	"errors"
	"sort"
)

var (
	errChildNotFound error = errors.New("child not found")
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

// DelChild deletes a Child group from the children map.  If the child does not exist, an error is returned.
func (children Children) DelChild(childName string) error {
	_, ok := children[childName]
	if !ok {
		return errChildNotFound
	}
	delete(children, childName)
	return nil
}

func (children Children) GetChild(childName string) (Child, error) {
	child, ok := children[childName]
	if !ok {
		return child, errChildNotFound
	}
	return child, nil
}

// getAllChildren returns a slice of all the keys in the Children map.
// The "all" group can be excluded by setting includeAll=false.
func (children Children) GetAllChildren(includeAll bool) (c []string) {
	for k := range children {
		// "all" is special and should not be a child of itself
		if !includeAll && k == "all" {
			continue
		}
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

// AddMember adds a member server to a specified Child group.  If the group doesn't exists, it will be created.
func (children Children) AddMember(childName, hostName string) {
	var child Child
	if _, ok := children[childName]; ok {
		child = children[childName]
	} else {
		child.Name = childName
		child.Members = make(map[string]int)
		children[childName] = child
	}
	child.Members[hostName] = 1
}
