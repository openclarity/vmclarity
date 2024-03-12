package types

func (dst *Annotations) MergeFrom(src *Annotations) {
	*dst = append(*dst, *src...)
}

func (dst *ScanSummary) MergeFrom(src *ScanSummary) {
	// TODO(ramizpolic): Implement later
}
