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
			err = children.AddMember(test.childName, member)
			if err != nil {
				t.Error(err)
			}
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

func TestNoChild(t *testing.T) {
	var err error
	children := NewChildren()
	err = children.AddMember("fakeChild", "host")
	if err == nil {
		t.Error("No error returned when adding to a non-existent child group")
	} else if err != errChildNotFound {
		t.Errorf("Expected an errChildNotFound but got %v", err)
	}
	_, err = children.MemberSlice("fakeChild")
	if err == nil {
		t.Error("No error returned when attempting to retrieve a non-existent child group")
	} else if err != errChildNotFound {
		t.Errorf("Expected an errChildNotFound but got %v", err)
	}
}
