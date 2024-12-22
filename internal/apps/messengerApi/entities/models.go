package entities

import "time"

type MessageInfo struct {
	Id             string     `json:"id,omitempty" bson:"_id,omitempty"`
	To             string     `json:"to,omitempty" bson:"to,omitempty"`
	Content        string     `json:"content,omitempty" bson:"content,omitempty"`
	Status         string     `json:"status,omitempty" bson:"status,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty" bson:"sent_at,omitempty"`
	FailedToSentAt *time.Time `json:"failed_to_sent_at,omitempty" bson:"failed_to_sent_at,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
