package store

import "time"

// BusinessHoursRule business hours rule model, used to save one rule of business hour,
// save multiple instance in Store model to apply multiple rules, the later
// ones overwrite prior ones. The format of FromTime and ToTime is hh:mm
type BusinessHoursRule struct {
	FromDay  time.Weekday `json:"from_day" bson:"fd" validate:"required,min=1,max=7"`
	ToDay    time.Weekday `json:"to_day" bson:"td" validate:"required,min=1,max=7"`
	FromTime string       `json:"from_time" bson:"ft" validate:"required,hh-mm"`
	ToTime   string       `json:"to_time" bson:"tt" validate:"required,hh-mm"`
}
