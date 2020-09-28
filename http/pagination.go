package http

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func CreatePageToken(arrayData interface{}, dataLimit int, fieldNameId string, fieldNameDate string) (nextToken string) {
	voData := reflect.ValueOf(arrayData)
	if voData.Kind() != reflect.Slice {
		return
	}

	if !(voData.IsValid() && voData.Len() >= dataLimit) {
		return
	}

	lastRow := voData.Index(voData.Len() - 1)
	nextTimestamp := lastRow.FieldByName(fieldNameDate).Interface()
	timestamp, ok := nextTimestamp.(time.Time)
	if !ok {
		return
	}

	uuid := lastRow.FieldByName(fieldNameId).String()

	nextToken = fmt.Sprintf("%d_%s", timestamp.Unix(), uuid)
	return
}

func ParsePageToken(pageToken string) (token map[string]string) {
	pToken := strings.Split(pageToken, "_")
	token = make(map[string]string)
	var date string
	var id string

	if len(pToken) != 2 {
		return
	}

	date = pToken[0]
	id = pToken[1]

	dateInt, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		return
	}

	token["date"] = time.Unix(dateInt, 0).Format(time.RFC3339)
	token["id"] = id

	return
}
