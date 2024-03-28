package types

func (a *ScanSummary) Add(b *ScanSummary) {
	if a == nil || b == nil {
		return
	}

	a.JobsDone += b.JobsDone
	a.JobsFailed += b.JobsFailed
	a.JobsRemaining += b.JobsRemaining

	a.JobsTotal = new(int)
	*a.JobsTotal += a.JobsDone
	*a.JobsTotal += a.JobsFailed
	*a.JobsTotal += a.JobsRemaining
}
