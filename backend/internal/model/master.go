package model

import "github.com/google/uuid"

type Master struct {
    ID              uuid.UUID `db:"id"               json:"id"`
    FullName        string    `db:"full_name"         json:"full_name"`
    Bio             string    `db:"bio"               json:"bio"`
    ExperienceYears int       `db:"experience_years"  json:"experience_years"`
    AvatarURL       string    `db:"avatar_url"        json:"avatar_url"`
    Instagram       string    `db:"instagram"         json:"instagram"`
    IsActive        bool      `db:"is_active"         json:"is_active"`
    Services        []Service `db:"-"                 json:"services,omitempty"`
}