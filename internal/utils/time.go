package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)


type TimeUtils struct{}


func NewTimeUtils() *TimeUtils {
	return &TimeUtils{}
}

func (tu *TimeUtils) ConvertStringToPgTime(timeStr string) (pgtype.Time, error) {
if timeStr == "" {
	return pgtype.Time{Valid: false}, nil
}

timeFormats := []string{
	"15:04",      // 24-hour format without seconds
	"15:04:05",   // 24-hour format with seconds
	"3:04 PM",    // 12-hour format without seconds
	"3:04:05 PM", // 12-hour format with seconds
	"15:04:05.000", // 24-hour format with milliseconds
}


var parsedTime time.Time
var err error

for _, format := range timeFormats {
	parsedTime, err = time.Parse(format, timeStr)
	if err == nil {
		break
	}
}

if err != nil {
	return pgtype.Time{}, fmt.Errorf("invalid time format '%s': supported formats are HH:MM, HH:MM:SS, H:MM PM,  H:MM:SS PM", timeStr)
} 

microseconds := int64(parsedTime.Hour())*3600000000 +
int64(parsedTime.Minute())*60000000 +
int64(parsedTime.Second())*1000000  +
int64(parsedTime.Nanosecond())/1000

return pgtype.Time{
	Microseconds: microseconds,
	Valid:   true,
}, nil
}


func (tu *TimeUtils) ConvertStringToPgTimePtr(timeStr *string) (*pgtype.Time, error) {
  if timeStr == nil || *timeStr == "" {	
	return nil, nil
}

	pgTime, err := tu.ConvertStringToPgTime(*timeStr)
	if err != nil {
	return nil, err
}

	return &pgTime, nil
}


func (tu *TimeUtils) ConvertPgTimeToString(pgTime pgtype.Time) string {
	if !pgTime.Valid {
		return ""
	}

	totalSeconds := pgTime.Microseconds / 1000000
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600)  /  60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func (tu *TimeUtils) ValidateTimeFormat(timeStr string)  error {
	if timeStr == "" {
		return nil
	}
	_, err := tu.ConvertStringToPgTime(timeStr)
	return err
}


func (tu *TimeUtils) GetSupportedTimeFormats() []string{
	return []string{
		"15:04 (24-hour format)",
		"15:04:05 (24-hour format with seconds)",
		"3:04 PM (12-hour format)",
		"3:04:05 PM (12-hour format with seconds)",
		"15:04:05.000 (24-hour format with milliseconds)",
	}
}

func PgTimeToTime(pgTime pgtype.Time) time.Time {
	if !pgTime.Valid {
		return time.Time{}
	}

	seconds := pgTime.Microseconds / 1_000_000
	nanoseconds := (pgTime.Microseconds % 1_000_000)* 1000
	return time.Date(0, 1, 1, int(seconds/3600), int((seconds%3600)/60), int(seconds%60), int(nanoseconds), time.UTC)
}
