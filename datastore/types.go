package datastore

import (
	"database/sql"
	"encoding/json"
	"time"
)

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

type JsonNullDate struct {
	sql.NullString
}

func (v JsonNullDate) MarshalJSON() ([]byte, error) {
	if v.Valid {
		t, err := time.Parse(time.RFC3339, v.String)
		if err != nil {
			return json.Marshal(nil)
		}
		return json.Marshal(t.Unix())
	} else {
		return json.Marshal(nil)
	}
}

type OrganizationView struct {
	Id          JsonNullInt64  `json:"id,omitempty"`
	Name        JsonNullString `json:"name,omitempty"`
	Entitlement JsonNullString `json:"entitlement,omitempty"`
	ExternalID  JsonNullString `json:"externalid,omitempty" db:"externalID"`
	CreatedAt   JsonNullInt64  `json:"createdat,omitempty"`
	UpdatedAt   JsonNullInt64  `json:"updatedat,omitempty"`
	RenewAt     JsonNullDate   `json:"renewaldate,omitempty" db:"renewaldate"`
	TimeZone    JsonNullString `json:"timezone,omitempty"`
	GroupId     JsonNullInt64  `json:"groupid,omitempty"`
	Tam         JsonNullString `json:"tam,omitempty"`
	Se          JsonNullString `json:"se,omitempty"`
}

type TicketView struct {
	Id                JsonNullInt64  `json:"id,omitempty"`
	ExternalID        JsonNullString `json:"externalid,omitempty"`
	Subject           JsonNullString `json:"subject,omitempty"`
	Status            JsonNullString `json:"status,omitempty"`
	RequesterId       JsonNullInt64  `json:"requesterid,omitempty"`
	SubmitterId       JsonNullInt64  `json:"submitterid,omitempty"`
	AssigneeId        JsonNullInt64  `json:"AssigneeId,omitempty"`
	Recipient         JsonNullString `json:"recipient,omitempty"`
	OrganizationId    JsonNullInt64  `json:"organizationid,omitempty"`
	GroupId           JsonNullInt64  `json:"groupid,omitempty"`
	CreatedAt         JsonNullInt64  `json:"createdat,omitempty"`
	UpdatedAt         JsonNullInt64  `json:"updatedat,omitempty"`
	Priority          JsonNullString `json:"priority,omitempty"`
	Component         JsonNullString `json:"component,omitempty"`
	TicketTime        JsonNullInt64  `json:"tickettime,omitempty" db:"ticketTime"`
	Cause             JsonNullString `json:"cause,omitempty"`
	Version           JsonNullString `json:"version,omitempty"`
	BundleUsed        JsonNullString `json:"bundleused"`
	TTFR              JsonNullInt64  `json:"ttfr,omitempty"`
	TTR               JsonNullInt64  `json:"ttr,omitempty"`
	SolvedAt          JsonNullInt64  `json:"solvedat,omitempty"`
	AgentWaitTime     JsonNullInt64  `json:"agentwaittime,omitempy"`
	RequesterWaitTime JsonNullInt64  `json:"requesterwaittime,omitempty"`
}
