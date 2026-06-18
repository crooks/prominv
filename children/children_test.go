package children

import (
	"testing"
)

func isMember(stringSlice []string, key string) bool {
	for _, s := range stringSlice {
		if s == key {
			return true
		}
	}
	return false
}

func TestNewChild(t *testing.T) {
	var testdata = []struct {
		childName       string
		hostNames       []string
		addChildSuccess bool
	}{
		{"exampleChildGroup", []string{"exampleHost1", "exampleHost2"}, true},
		{"exampleChildGroup", []string{"exampleHost1", "exampleHost1"}, false},
	}
	children := NewChildren()
	for _, test := range testdata {
		err := children.AddChild(test.childName)
		if err == nil && !test.addChildSuccess {
			t.Error("duplicate child failed to raise an error")
		}
		if err != nil && test.addChildSuccess {
			t.Errorf("Error adding %s: %v", test.childName, err)
		}
		for _, member := range test.hostNames {
			children.AddMember(test.childName, member)
		}
		outMembers, err := children.MemberSlice(test.childName)
		if err != nil {
			t.Error(err)
		}
		if len(outMembers) != len(test.hostNames) {
			t.Errorf("testSlice contains %d records instead of %d", len(outMembers), len(test.hostNames))
		}
		for _, item := range test.hostNames {
			if !isMember(outMembers, item) {
				t.Errorf("Hostname %s not found in results slice", item)
			}
		}
	} // Loop of tests
}

func TestGetAllChildren(t *testing.T) {
	inChildren := []string{"abc", "def", "ghi"}
	children := NewChildren()
	for _, newChild := range inChildren {
		children.AddChild(newChild)
	}
	outChildren := children.GetAllChildren(true)
	if len(inChildren) != len(outChildren) {
		t.Errorf("Unexpected number of Child records. Wanted: %d, Got: %d", len(inChildren), len(outChildren))
	}
	for _, ic := range inChildren {
		matched := false
		for _, oc := range outChildren {
			if oc == ic {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("Test element %s not found in AllChildren output", ic)
			t.Errorf("%v vs. %v", inChildren, outChildren)
		}
	}
}

func TestDelChild(t *testing.T) {
	inChildren := []string{"abc", "def", "ghi"}
	children := NewChildren()
	for _, newChild := range inChildren {
		children.AddChild(newChild)
	}
	children.DelChild("ghi")
	if len(children) != 2 {
		t.Errorf("Unexpected number of children.  Expected=2, Got=%d", len(children))
	}
	// Try and delete the same key again. Nothing should change.
	children.DelChild("ghi")
	if len(children) != 2 {
		t.Errorf("Unexpected number of children.  Expected=2, Got=%d", len(children))
	}
}

func TestNewMember(t *testing.T) {
	inChild := "tc"
	inMember := "tm"
	children := NewChildren()
	children.AddMember(inChild, inMember)

	// This tests the GetChild function can return the created Child.
	outChild, err := children.GetChild(inChild)
	if err != nil {
		t.Errorf("Child retrieval failed with: %v", err)
	}
	if outChild.Name != inChild {
		t.Errorf("Unexpected child name.  Wanted: %s, Got: %s", inChild, outChild.Name)
	}

	// Pretty much the same test as above but this time recovering the Members of the Child.
	members, err := children.MemberSlice(inChild)
	if len(members) != 1 {
		t.Errorf("Unexpected number of members in %s.  Should be 1 but have %d", inChild, len(members))
	}
	if members[0] != inMember {
		t.Errorf("Unexpected member name in %s. Wanted: %s, Got: %s", inChild, inMember, members[0])
	}
}
