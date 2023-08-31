package models

// Used only for swagger

type SegmentAssignRequest struct {
	SegmentsToDeleteNames []string `json:"to_delete"`
	SegmentsToAddNames    []string `json:"to_add"`
}

type SegmentAssignWithTTLRequest struct {
	SegmentName string `json:"segment_name" validate:"required"`
	TTL         int32  `json:"ttl" validate:"required,gt=0"`
}
