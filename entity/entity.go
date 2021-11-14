package entity

type Envelope struct {
	EnvelopeID int64 `json:"envelope_id"`
	UserID     int64 `json:"uid"`
	Opened     bool  `json:"opened"`
	Value      int64 `json:"value"`
	SnatchTime int64 `json:"snatch_time"`
}
