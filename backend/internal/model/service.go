package model

import "github.com/google/uuid"

type Service struct {
    ID          uuid.UUID `db:"id"           json:"id"`
    Name        string    `db:"name"         json:"name"`
    Category    string    `db:"category"     json:"category"`
    Description string    `db:"description"  json:"description"`
    DurationMin int       `db:"duration_min" json:"duration_min"`
    Price       float64   `db:"price"        json:"price"`
    IsActive    bool      `db:"is_active"    json:"is_active"`
    SortOrder   int       `db:"sort_order"   json:"sort_order"`
    NameKz string `json:"name_kz"`
    NameEn string `json:"name_en"`
    DescKz string `json:"desc_kz"`
    DescEn string `json:"desc_en"`
}

type MasterService struct {
    Service
    CustomPrice *float64 `db:"custom_price" json:"custom_price"`
}

func (ms MasterService) EffectivePrice() float64 {
    if ms.CustomPrice != nil && *ms.CustomPrice > 0 {
        return *ms.CustomPrice
    }
    return ms.Price
}

type Schedule struct {
    MasterID  string `db:"master_id"`
    DayOfWeek int    `db:"day_of_week"`
    StartTime string `db:"start_time"`
    EndTime   string `db:"end_time"`
}