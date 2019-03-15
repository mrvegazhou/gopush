package helper

import (
	"fmt"
	"time"
)

type LocalTime time.Time

// MarshalJSON satify the json marshal interface
func (l LocalTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(l).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type LocalDate time.Time

// MarshalJSON satify the json marshal interface
func (l LocalDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(l).Format("2006-01-02"))
	return []byte(stamp), nil
}

// FormatTime 格式化时间
func FormatTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// UnFormatTime 将时间字符串转为Time
func UnFormatTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

// UnixTime 将时间转化为毫秒数
func UnixTime(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

func NowUnixTime() int64 {
	return time.Now().UnixNano() / 1000000
}

// UnunixTime 将毫秒数转为时间
func UnunixTime(unix int64) time.Time {
	return time.Unix(0, unix*1000000)
}