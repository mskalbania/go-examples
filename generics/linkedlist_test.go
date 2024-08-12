package generics

import (
	"slices"
	"testing"
)

func TestShouldCorrectlyAddElementsToList(t *testing.T) {
	t.Log("Given empty linked list")
	{
		list := NewLinkedList[int]()
		t.Log("When some elements added")
		{
			list.Add(1)
			list.Add(2)
			t.Log("Then Elements added")
			{
				if list.start == nil || list.start.val != 1 {
					t.Errorf("Expected 1, got %v", list.start.val)
				}
				if list.start.next == nil || list.start.next.val != 2 {
					t.Errorf("Expected 2, got %v", list.start.next.val)
				}
			}
		}
	}
}

func TestShouldCorrectlyBuildStringRepresentation(t *testing.T) {
	t.Log("Given list with elements")
	{
		list := NewLinkedList[int]()
		list.Add(1)
		list.Add(2)
		t.Log("When string representation requested")
		{
			str := list.String()
			t.Log("Then correct string representation returned")
			{
				if str != "[1, 2]" {
					t.Errorf("Expected '[1, 2]', got %v", str)
				}
			}
		}
	}
}

func TestShouldCorrectlyRemoveElements(t *testing.T) {
	testData := []struct {
		description               string
		givenList                 *LinkedList[int]
		removeElement             int
		expectedReturnResult      bool
		expectedElementsRemaining []int
	}{
		{
			description:               "remove from empty list",
			givenList:                 buildLinkedList([]int{}),
			removeElement:             1,
			expectedReturnResult:      false,
			expectedElementsRemaining: []int{},
		},
		{
			description:               "remove from middle",
			givenList:                 buildLinkedList([]int{1, 2, 3}),
			removeElement:             2,
			expectedReturnResult:      true,
			expectedElementsRemaining: []int{1, 3},
		},
		{
			description:               "remove from end",
			givenList:                 buildLinkedList([]int{1, 2, 3}),
			removeElement:             3,
			expectedReturnResult:      true,
			expectedElementsRemaining: []int{1, 2},
		},
	}
	for _, testCase := range testData {
		t.Logf("Scenario: %s", testCase.description)
		t.Logf(" Given list: %s", testCase.givenList)
		{
			t.Log(" When element removed")
			{
				result := testCase.givenList.Delete(testCase.removeElement)
				t.Log(" Then correct status returned")
				{
					if result != testCase.expectedReturnResult {
						t.Errorf("  Expected result %v, got %v", testCase.expectedReturnResult, result)
					}
					asSlice := toSlice(testCase.givenList)
					if !slices.Equal(asSlice, testCase.expectedElementsRemaining) {
						t.Errorf("  Element not removed, expected %v, got %v", testCase.expectedElementsRemaining, asSlice)
					}
				}
			}
		}
	}
}

func TestShouldCorrectlySearch(t *testing.T) {
	testData := []struct {
		description   string
		givenList     *LinkedList[int]
		elementToFind int
		expectedIndex int
	}{
		{
			"search in empty",
			buildLinkedList([]int{}),
			1,
			-1,
		},
		{
			"search existing in the middle",
			buildLinkedList([]int{1, 2, 3}),
			2,
			1,
		},
		{
			"search existing in the end",
			buildLinkedList([]int{1, 2, 3}),
			3,
			2,
		},
		{
			"search non-existing",
			buildLinkedList([]int{1, 2, 3}),
			4,
			-1,
		},
	}
	for _, testCase := range testData {
		t.Logf("Scenario: %s", testCase.description)
		t.Logf(" Given list: %s", testCase.givenList)
		{
			t.Log(" When element searched")
			{
				result := testCase.givenList.Search(testCase.elementToFind)
				t.Log(" Then correct index returned")
				{
					if result != testCase.expectedIndex {
						t.Errorf("  Expected index %v, got %v", testCase.expectedIndex, result)
					}
				}
			}
		}
	}
}

func buildLinkedList(elements []int) *LinkedList[int] {
	list := NewLinkedList[int]()
	for _, e := range elements {
		list.Add(e)
	}
	return list
}

func toSlice(list *LinkedList[int]) []int {
	curr := list.start
	out := make([]int, 0)
	for {
		if curr == nil {
			break
		}
		out = append(out, curr.val)
		curr = curr.next
	}
	return out

}
