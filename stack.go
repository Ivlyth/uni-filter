package uni_filter

type Stack[T any] struct {
	entries []T
}

// NewStack creates a new stack object.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Clear removes all data from the stack
func (s *Stack[T]) Clear() {
	s.entries = []T{}
}

// IsEmpty returns true if the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.entries) == 0
}

// Size retrieves the number of entries stored upon the stack.
func (s *Stack[T]) Size() int {
	return len(s.entries)
}

// Push appends the specified value to the stack.
func (s *Stack[T]) Push(v T) {
	s.entries = append(s.entries, v)
}

// Pop removes a value from the stack.
func (s *Stack[T]) Pop() (v T) {

	// get the last entry.
	result := s.entries[len(s.entries)-1]

	// remove it
	s.entries = s.entries[:len(s.entries)-1]

	return result
}

// Top return the latest value from the stack, but not remove it
func (s *Stack[T]) Top() (v T) {
	return s.entries[len(s.entries)-1]
}
