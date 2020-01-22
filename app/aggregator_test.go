package app

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestAggregator(t *testing.T) {
	currentTime, _ := removeMonotonicTimeFromTime(time.Now())
	s := Session{
		ID:           primitive.NewObjectID(),
		Link:         "/",
		CreatedTime:  currentTime.Unix(),
		Patient:      UserRef{},
		Provider:     UserRef{},
		TextMetrics:  []TextMetrics{
			{
				Time:      currentTime.Add(10*time.Second).Unix(),
				Text:      "",
				Sentiment: .1,
			},
			{
				Time:      currentTime.Add(20*time.Second).Unix(),
				Text:      "",
				Sentiment: .8,
			},
			{
				Time:      currentTime.Add(30*time.Second).Unix(),
				Text:      "",
				Sentiment: 0,
			},
			{
				Time:      currentTime.Add(40*time.Second).Unix(),
				Text:      "",
				Sentiment: -0.9,
			},
		},
		ImageMetrics: []FrameMetrics{
			{
				Time:          currentTime.Add(10*time.Second).Unix(),
				ImageFilename: "",
				Emotion:       map[string]string{
					"joy": "VERY_LIKELY",
					"sorrow": "VERY_UNLIKELY",
					"anger": "VERY_UNLIKELY",
					"surprise": "POSSIBLE",
				},
				AU:            nil,
			},
			{
				Time:          currentTime.Add(20*time.Second).Unix(),
				ImageFilename: "",
				Emotion:       map[string]string{
					"joy": "VERY_LIKELY",
					"sorrow": "VERY_UNLIKELY",
					"anger": "VERY_UNLIKELY",
					"surprise": "POSSIBLE",
				},
				AU:            nil,
			},
			{
				Time:          currentTime.Add(30*time.Second).Unix(),
				ImageFilename: "",
				Emotion:       map[string]string{
					"joy": "VERY_UNLIKELY",
					"sorrow": "VERY_LIKELY",
					"anger": "VERY_UNLIKELY",
					"surprise": "VERY_UNLIKELY",
				},
				AU:            nil,
			},
			{
				Time:          currentTime.Add(40*time.Second).Unix(),
				ImageFilename: "",
				Emotion:       map[string]string{
					"joy": "VERY_UNLIKELY",
					"sorrow": "VERY_LIKELY",
					"anger": "VERY_LIKELY",
					"surprise": "POSSIBLE",
				},
				AU:            nil,
			},
		},
	}
	agg, err := AggregatorFromSession(&s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	res, err := agg.Run()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(res)
}