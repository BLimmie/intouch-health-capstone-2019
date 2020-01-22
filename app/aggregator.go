package app

import (
	"errors"
	"sort"
	"time"
)

func AggregatorFromSession(session *Session) (*Aggregator, error) {
	defaultAU := make(map[string]float32)
	for _, label := range auLabels {
		defaultAU[label] = 0
	}
	agg := &Aggregator{
		session:     session,
		conclusions: make(map[string]func(session *Session) (interface{}, error)),
	}
	session.TextMetrics = append([]TextMetrics{
		{
			Time:      session.CreatedTime,
			Text:      "",
			Sentiment: 0,
		}}, session.TextMetrics...)
	session.ImageMetrics = append([]FrameMetrics{
		{
			Time:          session.CreatedTime,
			ImageFilename: "",
			Emotion: map[string]string{
				"joy":      "NOT_LIKELY",
				"sorrow":   "NOT_LIKELY",
				"anger":    "NOT_LIKELY",
				"surprise": "NOT_LIKELY",
			},
			AU: defaultAU,
		},
	}, session.ImageMetrics...)
	sort.Slice(session.TextMetrics, func(i, j int) bool {
		t1 := timeFromInt(session.TextMetrics[i].Time)
		t2 := timeFromInt(session.TextMetrics[j].Time)
		return t1.Before(t2)
	})
	sort.Slice(session.ImageMetrics, func(i, j int) bool {
		t1 := timeFromInt(session.ImageMetrics[i].Time)
		t2 := timeFromInt(session.ImageMetrics[j].Time)
		return t1.Before(t2)
	})
	agg.init()
	return agg, nil
}

func (agg *Aggregator) init() {
	agg.conclusions["Average Text Sentiment"] = isTextPositive
	agg.conclusions["Time in Facial Emotions"] = timeSpentEmotion
	agg.conclusions["Percent in Facial Emotion over last 10 seconds"] = emotionOverTime_RunningAggregate
}

func (agg *Aggregator) Run() (interface{}, error) {
	fullRes := make(map[string]interface{})
	for key, f := range agg.conclusions {
		res, err := f(agg.session)
		if err != nil {
			return nil, err
		}
		fullRes[key] = res
	}
	return fullRes, nil
}

func isTextPositive(session *Session) (interface{}, error) {
	var sum float32
	for _, m := range session.TextMetrics {
		sum += m.Sentiment
	}
	return map[string]float32{
		"AvgTextSentiment": sum / float32(len(session.TextMetrics)),
	}, nil
}

func timeSpentEmotion(session *Session) (interface{}, error) {
	emotions := map[string]time.Duration{
		"joy":      time.Duration(0),
		"sorrow":   time.Duration(0),
		"anger":    time.Duration(0),
		"surprise": time.Duration(0),
	}
	lastTime := time.Unix(session.CreatedTime, 0)
	startTime := lastTime

	for _, m := range session.ImageMetrics {
		t := timeFromInt(m.Time)
		dur := t.Sub(lastTime)
		for emotion, level := range m.Emotion {
			if level == "VERY_LIKELY" {
				emotions[emotion] += dur
			} else if level == "POSSIBLE" {
				emotions[emotion] += dur / 2
			}
		}
		lastTime = t
	}
	var emotions_s = make(map[string]string)
	totalTime := lastTime.Sub(startTime)
	var emotions_p = make(map[string]float32)
	for e, t := range emotions {
		emotions_s[e] = t.String()
		emotions_p[e] = float32(t) / float32(totalTime)
	}
	return map[string]interface{}{
		"Total Time": emotions_s,
		"Percentage": emotions_p,
	}, nil
}

func emotionOverTime_RunningAggregate(session *Session) (interface{}, error) {
	initialTime := timeFromInt(session.CreatedTime)
	fullList := make(map[int]interface{})
	lastTime := timeFromInt(session.ImageMetrics[len(session.ImageMetrics)-1].Time)
	seconds := 10
	for currentTime := initialTime.Add(time.Second * 10); currentTime.Before(lastTime); currentTime = currentTime.Add(time.Second) {
		tenSecondsAgo := currentTime.Add(time.Duration(-10)*time.Second)
		firstSegmentIdx := 0
		for {
			t := timeFromInt(session.ImageMetrics[firstSegmentIdx].Time)
			if t.Before(tenSecondsAgo) || t.Equal(tenSecondsAgo){
				//This should guarantee that firstSegmentIdx>0, bug if not
				firstSegmentIdx += 1
			} else {
				break
			}
			if firstSegmentIdx == 0 {
				return nil, errors.New("frame sentiment not initialized properly")
			}
			if firstSegmentIdx >= len(session.ImageMetrics){
				return nil, errors.New("out of bounds in image metrics")
			}
		}
		//time of first segment that's after initial time
		m := session.ImageMetrics[firstSegmentIdx-1]
		n := session.ImageMetrics[firstSegmentIdx]
		t := timeFromInt(n.Time)
		dur := t.Sub(tenSecondsAgo)
		emotions := map[string]time.Duration{
			"joy":      time.Duration(0),
			"sorrow":   time.Duration(0),
			"anger":    time.Duration(0),
			"surprise": time.Duration(0),
		}
		for emotion, level := range m.Emotion {
			if level == "VERY_LIKELY" {
				emotions[emotion] += dur
			} else if level == "POSSIBLE" {
				emotions[emotion] += dur / 2
			}
		}
		for i := firstSegmentIdx; timeFromInt(session.ImageMetrics[i].Time).Before(currentTime); i++ {
			m := session.ImageMetrics[i]
			nextTime := timeFromInt(session.ImageMetrics[i+1].Time)
			if nextTime.After(currentTime) {
				nextTime = currentTime
			}
			dur := nextTime.Sub(timeFromInt(m.Time))
			for emotion, level := range m.Emotion {
				if level == "VERY_LIKELY" {
					emotions[emotion] += dur
				} else if level == "POSSIBLE" {
					emotions[emotion] += dur / 2
				}
			}
		}
		var emotions_s = make(map[string]string)
		totalTime := time.Second*10
		var emotions_p = make(map[string]float32)
		for e, t := range emotions {
			emotions_s[e] = t.String()
			emotions_p[e] = float32(t) / float32(totalTime)
		}
		fullList[seconds] = map[string]interface{}{
			"Total Time": emotions_s,
			"Percentage": emotions_p,
		}
		seconds++
	}
	return fullList, nil
}
