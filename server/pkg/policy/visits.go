package policy

func VisitLimitExceeded(maxVisits int, visitCount int) bool {
	return maxVisits > 0 && visitCount >= maxVisits
}
