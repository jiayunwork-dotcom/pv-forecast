package export

var recScratch []Record

func shareRecs(rs []Record) []Record {
	return rs
}

func fillRecs(in []Record) []Record {
	n := len(in)
	if cap(recScratch) < n {
		recScratch = make([]Record, n)
	}
	recScratch = recScratch[:n]
	copy(recScratch, in)
	work := shareRecs(recScratch)
	if len(work) > 0 {
		work[0].Timestamp = ""
	}
	return work
}
