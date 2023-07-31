package go_train

import (
	"fmt"
	"testing"
	"time"
)

const YYYY_MM_DD = "2006-01-02"

func timeProcessingV1(sdaytime, edaytime time.Time) []interface{} {
	var fieldList []interface{}
	if sdaytime.Equal(edaytime) {
		fieldList = append(fieldList, sdaytime.Format(YYYY_MM_DD))
		return fieldList
	}
	// 逐一增加天数，直到达到结束日期
	for currentTime := sdaytime; !currentTime.After(edaytime); currentTime = currentTime.AddDate(0, 0, 1) {
		fieldList = append(fieldList, currentTime.Format(YYYY_MM_DD))
	}
	return fieldList
}

func TestTimeProcessingV1(t *testing.T) {
	t1, _ := time.Parse(YYYY_MM_DD, "2023-07-19")
	t2, _ := time.Parse(YYYY_MM_DD, "2023-07-18")
	res := timeProcessing(t1, t2)
	fmt.Println(res)
}

func timeProcessing(sdaytime, edaytime time.Time) []interface{} {
	var fieldList []interface{}

	if sdaytime.Equal(edaytime) {
		dateStr := sdaytime.Format(YYYY_MM_DD)
		fieldList = append(fieldList, dateStr)
		return fieldList
	}
	// 输出日期格式固定
	edaytime2Str := edaytime.Format(YYYY_MM_DD)
	fieldList = append(fieldList, sdaytime.Format(YYYY_MM_DD))
	for {
		sdaytime = sdaytime.AddDate(0, 0, 1)
		dateStr := sdaytime.Format(YYYY_MM_DD)
		fieldList = append(fieldList, dateStr)
		if dateStr == edaytime2Str {
			break
		}
	}
	return fieldList
}
