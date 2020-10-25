package animation

import (
	"github.com/wieku/danser-go/bmath"
	"math"
)

type TransformationType int64
type TransformationStatus int64

const (
	Fade = TransformationType(1 << iota)
	Rotate
	Scale
	ScaleVector
	Move
	MoveX
	MoveY
	Color3
	Color4
	HorizontalFlip
	VerticalFlip
	Additive
)

const (
	NotStarted = TransformationStatus(1 << iota)
	Going
	Ended
)

func timeClamp(start, end, time float64) float64 {
	return math.Max(0, math.Min(1.0, (time-start)/(end-start)))
}

type Transformation struct {
	transformationType TransformationType
	startValues        [4]float64
	endValues          [4]float64
	easing             func(float64) float64
	startTime, endTime float64
}

func NewBooleanTransform(transformationType TransformationType, startTime, endTime float64) *Transformation {
	if transformationType&(HorizontalFlip|VerticalFlip|Additive) == 0 {
		panic("Wrong TransformationType used!")
	}

	return &Transformation{transformationType: transformationType, startTime: startTime, endTime: endTime}
}

func NewSingleTransform(transformationType TransformationType, easing func(float64) float64, startTime, endTime, startValue, endValue float64) *Transformation {
	if transformationType&(Fade|Rotate|Scale|MoveX|MoveY) == 0 {
		panic("Wrong TransformationType used!")
	}

	transformation := &Transformation{transformationType: transformationType, startTime: startTime, endTime: endTime, easing: easing}
	transformation.startValues[0] = startValue
	transformation.endValues[0] = endValue
	return transformation
}

func NewVectorTransform(transformationType TransformationType, easing func(float64) float64, startTime, endTime, startValueX, startValueY, endValueX, endValueY float64) *Transformation {
	if transformationType&(ScaleVector|Move) == 0 {
		panic("Wrong TransformationType used!")
	}

	transformation := &Transformation{transformationType: transformationType, startTime: startTime, endTime: endTime, easing: easing}
	transformation.startValues[0] = startValueX
	transformation.startValues[1] = startValueY
	transformation.endValues[0] = endValueX
	transformation.endValues[1] = endValueY
	return transformation
}

func NewVectorTransformV(transformationType TransformationType, easing func(float64) float64, startTime, endTime float64, start, end bmath.Vector2d) *Transformation {
	if transformationType&(ScaleVector|Move) == 0 {
		panic("Wrong TransformationType used!")
	}

	transformation := &Transformation{transformationType: transformationType, startTime: startTime, endTime: endTime, easing: easing}
	transformation.startValues[0] = start.X
	transformation.startValues[1] = start.Y
	transformation.endValues[0] = end.X
	transformation.endValues[1] = end.Y
	return transformation
}

func NewColorTransform(transformationType TransformationType, easing func(float64) float64, startTime, endTime float64, start, end bmath.Color) *Transformation {
	if transformationType&(ScaleVector|Move) == 0 {
		panic("Wrong TransformationType used!")
	}

	transformation := &Transformation{transformationType: transformationType, startTime: startTime, endTime: endTime, easing: easing}
	transformation.startValues[0] = start.R
	transformation.startValues[1] = start.G
	transformation.startValues[2] = start.B
	transformation.startValues[3] = start.A

	transformation.endValues[0] = end.R
	transformation.endValues[1] = end.G
	transformation.endValues[2] = end.B
	transformation.endValues[3] = end.A
	return transformation
}

//Missing color

func (t *Transformation) GetStatus(time float64) TransformationStatus {
	if time < t.startTime {
		return NotStarted
	} else if time >= t.endTime {
		return Ended
	}
	return Going
}

func (t *Transformation) getProgress(time float64) float64 {
	return t.easing(timeClamp(t.startTime, t.endTime, time))
}

func (t *Transformation) GetSingle(time float64) float64 {
	return t.startValues[0] + t.getProgress(time)*(t.endValues[0]-t.startValues[0])
}

func (t *Transformation) GetDouble(time float64) (float64, float64) {
	progress := t.getProgress(time)
	return t.startValues[0] + progress*(t.endValues[0]-t.startValues[0]), t.startValues[1] + progress*(t.endValues[1]-t.startValues[1])
}

func (t *Transformation) GetVector(time float64) bmath.Vector2d {
	return bmath.NewVec2d(t.GetDouble(time))
}

func (t *Transformation) GetBoolean(time float64) bool {
	return time >= t.startTime && time < t.endTime
}

func (t *Transformation) GetColor(time float64) bmath.Color {
	progress := t.getProgress(time)
	return bmath.Color{
		R: t.startValues[0] + progress*(t.endValues[0]-t.startValues[0]),
		G: t.startValues[1] + progress*(t.endValues[1]-t.startValues[1]),
		B: t.startValues[2] + progress*(t.endValues[2]-t.startValues[2]),
		A: t.startValues[3] + progress*(t.endValues[3]-t.startValues[3]),
	}
}

func (t *Transformation) GetStartTime() float64 {
	return t.startTime
}

func (t *Transformation) GetEndTime() float64 {
	return t.endTime
}

func (t *Transformation) GetType() TransformationType {
	return t.transformationType
}
