package main

// Storage structure that keeps multiple copies of identical objects with a limit
type Storage struct {
	objects    map[*Product]int
	maxStorage int
	crtStorage int
}

// NewStorage creates a new *Storage
func NewStorage(maxStorage int) *Storage {
	s := new(Storage)
	s.maxStorage = maxStorage
	s.objects = make(map[*Product]int)

	return s
}

// Add add c Products to the storage, as long as it does not go above the maximum storage
func (s *Storage) Add(p *Product, c int) int {
	var toAdd int

	potentialMax := s.crtStorage + c
	if potentialMax > s.maxStorage {
		toAdd = s.maxStorage - s.crtStorage
	} else {
		toAdd = c
	}

	previousCount, present := s.objects[p]
	if !present {
		s.objects[p] = toAdd
	} else {
		s.objects[p] = previousCount + toAdd
	}

	s.crtStorage += toAdd
	return toAdd
}

// Remove extract c elements from Storage or the currently available amount if above
func (s *Storage) Remove(p *Product, c int) int {
	priorCount, present := s.objects[p]
	if !present {
		return 0
	}

	var removed int
	if priorCount > c {
		s.objects[p] = priorCount - c
		removed = c
	} else {
		delete(s.objects, p)
		removed = priorCount
	}

	s.crtStorage -= removed
	return removed
}

// Capacity returns the maximum number of Products that can be stored
func (s *Storage) Capacity() int {
	return s.maxStorage
}

// Size returns the current number of Products that are stored
func (s *Storage) Size() int {
	return s.crtStorage
}

// UniqueObjects returns the current number of unique objects
func (s *Storage) UniqueObjects() int {
	return len(s.objects)
}
