package objects

import (
	"log"
	"math"
)

type TimingPoint struct {
	Time         int64
	BaseBpm, Bpm float64
	SampleSet    int
	SampleIndex  int
	SampleVolume float64
	Kiai	bool
}

func (t TimingPoint) GetRatio() float64 {
	return t.Bpm / t.BaseBpm
}

type Timings struct {
	Points           []TimingPoint
	queue            []TimingPoint
	SliderMult       float64
	Current          TimingPoint
	fullBPM, partBPM float64
	BaseSet          int
	LastSet          int
	TickRate         float64
}

func NewTimings() *Timings {
	return &Timings{BaseSet: 1, LastSet: 1}
}

func (tim *Timings) AddPoint(time int64, bpm float64, sampleset, sampleindex int, samplevolume float64, isKiai bool) {
	point := TimingPoint{Time: time, Bpm: bpm, SampleSet: sampleset, SampleIndex: sampleindex, SampleVolume: samplevolume}
	if point.Bpm > 0 {
		tim.fullBPM = point.Bpm
	} else {
		point.Bpm = tim.fullBPM / math.Max(0.1, -100.0/point.Bpm)
	}
	point.BaseBpm = tim.fullBPM
	point.Kiai = isKiai
	tim.Points = append(tim.Points, point)
	tim.queue = append(tim.queue, point)
}

func (tim *Timings) Update(time int64) {
	if len(tim.queue) > 0 {
		p := tim.queue[0]
		if p.Time <= time {
			tim.queue = tim.queue[1:]
			tim.partBPM = p.Bpm
			tim.Current = p
		}
	}
}

func clamp(a, min, max int) int {
	if a > max {
		return max
	}
	if a < min {
		return min
	}
	return a
}

func (tim *Timings) GetPoint(time int64) TimingPoint {
	for i, pt := range tim.Points {
		if time < pt.Time {
			return tim.Points[clamp(i-1, 0, len(tim.Points)-1)]
		}
	}
	return tim.Points[len(tim.Points)-1]
}

func (tim Timings) GetSliderTimeS(time int64, pixelLength float64) int64 {
	res := int64(tim.GetPoint(time).Bpm * pixelLength / (100.0 * tim.SliderMult))
	if res < 0 {
		log.Println("E?", tim.GetPoint(time).Bpm, pixelLength, tim.SliderMult)
	}
	return res
}

func (tim Timings) GetSliderTime(pixelLength float64) int64 {
	return int64(tim.partBPM * pixelLength / (100.0 * tim.SliderMult))
}

func (tim Timings) GetSliderTimeP(point TimingPoint, pixelLength float64) int64 {
	return int64(point.Bpm * pixelLength / (100.0 * tim.SliderMult))
}

func (tim *Timings) Reset() {
	tim.queue = make([]TimingPoint, len(tim.Points))
	copy(tim.queue, tim.Points)
	tim.Current = tim.queue[0]
}

func (tim *Timings) Log() {
	log.Println(len(tim.Points))
}
