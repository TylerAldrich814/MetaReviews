package model

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defined a record type. Together with RecordID
// identifies unique records across all types.
type RecordType string

// Existing Record Types
const (
  RecordTypeMove = RecordType("movie")
)

// UserID defines a user id.
type UserID string 

// RatingValue defines a value of a rating record.
type RatingValue int32

// Rating defines an individual rating created by a user for some record
type Rating struct {
  RecordID   RecordID    `json:"record_id"`
  RecordType RecordType  `json:"record_type"`
  UserID     UserID      `json:"user_id"`
  Value      RatingValue `json:"value"`
}

type RatingEventType  string
const (
  RatingEventTypePut    = "put"
  RatingEventTypeDelete = "delete"
)

// RatingEvent defines an event containing rating
type RatingEvent struct {
  UserID     UserID          `json:"user_id"`
  RecordID   RecordID        `json:"record_id"`
  RecordType RecordType      `json:"record_type"`
  Value      RatingValue     `json:"value"`
  EventType  RatingEventType `json:"event_type"`
}
