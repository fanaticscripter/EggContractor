package api

import "sort"

type By func(m1, m2 *CoopStatus_Member) bool

var ByEggsLaid = By(func(m1, m2 *CoopStatus_Member) bool {
	return m1.EggsLaid > m2.EggsLaid
})

var ByLayingRate = By(func(m1, m2 *CoopStatus_Member) bool {
	return m1.EggsPerSecond > m2.EggsPerSecond
})

var ByEarningBonus = By(func(m1, m2 *CoopStatus_Member) bool {
	return m1.EarningBonusOom > m2.EarningBonusOom
})

func (by By) Sort(members []*CoopStatus_Member) {
	sort.Stable(&memberSorter{
		members: members,
		by:      by,
	})
}

type memberSorter struct {
	members []*CoopStatus_Member
	by      By
}

// Len is the number of elements in the collection.
func (s *memberSorter) Len() int {
	return len(s.members)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s *memberSorter) Less(i int, j int) bool {
	return s.by(s.members[i], s.members[j])
}

// Swap swaps the elements with indexes i and j.
func (s *memberSorter) Swap(i int, j int) {
	s.members[i], s.members[j] = s.members[j], s.members[i]
}
