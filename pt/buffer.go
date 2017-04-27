package pt

import (
	"image"
	"math"
)

type Channel int

const (
	ColorChannel = Channel(iota)
	VarianceChannel
	StandardDeviationChannel
	SamplesChannel
	DistanceChannel
	NormalChannel
)

type Pixel struct {
	Samples int
	M, V    Color

	DistSum      float64
	MNorm, VNorm Vector
}

func (p *Pixel) AddSample(sample Color) {
	p.Samples++
	if p.Samples == 1 {
		p.M = sample
		return
	}
	m := p.M
	p.M = p.M.Add(sample.Sub(p.M).DivScalar(float64(p.Samples)))
	p.V = p.V.Add(sample.Sub(m).Mul(sample.Sub(p.M)))
}

func (p *Pixel) AddSampleFeature(sample Features) {
	p.AddSample(sample.Color)
	if p.Samples == 1 {
		p.MNorm = sample.Normal
		return
	}
	m := p.MNorm
	p.MNorm = p.MNorm.Add(sample.Normal.Sub(p.MNorm).DivScalar(float64(p.Samples)))
	p.VNorm = p.VNorm.Add(sample.Normal.Sub(m).Mul(sample.Normal.Sub(p.MNorm)))
	p.DistSum += sample.Distance // XXX
}

func (p *Pixel) Color() Color {
	return p.M
}

func (p *Pixel) Variance() Color {
	if p.Samples < 2 {
		return Black
	}
	return p.V.DivScalar(float64(p.Samples - 1))
}

func (p *Pixel) Distance() float64 {
	return p.DistSum / float64(p.Samples)
}

func (p *Pixel) StandardDeviation() Color {
	return p.Variance().Pow(0.5)
}

func (p *Pixel) Normal() Color {
	n := p.MNorm
	n = n.AddScalar(1.0).DivScalar(2.0)
	return Color{n.X, n.Y, n.Z}
}

type Buffer struct {
	W, H   int
	Pixels []Pixel

	// XXX
	MaxDist float64
	MinDist float64
}

func NewBuffer(w, h int) *Buffer {
	pixels := make([]Pixel, w*h)
	return &Buffer{w, h, pixels, -math.MaxFloat64, math.MaxFloat64}
}

func (b *Buffer) Copy() *Buffer {
	pixels := make([]Pixel, b.W*b.H)
	copy(pixels, b.Pixels)
	return &Buffer{b.W, b.H, pixels, b.MaxDist, b.MinDist}
}

func (b *Buffer) AddSample(x, y int, sample Color) {
	b.Pixels[y*b.W+x].AddSample(sample)
}

func (b *Buffer) AddSampleFeature(x, y int, sample Features) {
	b.Pixels[y*b.W+x].AddSampleFeature(sample)
	b.MaxDist = math.Max(b.MaxDist, sample.Distance)
	b.MinDist = math.Min(b.MinDist, sample.Distance)
}

func (b *Buffer) Samples(x, y int) int {
	return b.Pixels[y*b.W+x].Samples
}

func (b *Buffer) Color(x, y int) Color {
	return b.Pixels[y*b.W+x].Color()
}

func (b *Buffer) Variance(x, y int) Color {
	return b.Pixels[y*b.W+x].Variance()
}

func (b *Buffer) StandardDeviation(x, y int) Color {
	return b.Pixels[y*b.W+x].StandardDeviation()
}

func (b *Buffer) Distance(x, y int) Color {
	d := b.Pixels[y*b.W+x].Distance()
	d = (d - b.MinDist) / (b.MaxDist - b.MinDist)
	return Color{d, d, d}
}

func (b *Buffer) Normal(x, y int) Color {
	return b.Pixels[y*b.W+x].Normal()
}

func (b *Buffer) Image(channel Channel) image.Image {
	result := image.NewRGBA64(image.Rect(0, 0, b.W, b.H))
	var maxSamples float64
	if channel == SamplesChannel {
		for _, pixel := range b.Pixels {
			maxSamples = math.Max(maxSamples, float64(pixel.Samples))
		}
	}
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			var c Color
			switch channel {
			case ColorChannel:
				c = b.Pixels[y*b.W+x].Color().Pow(1 / 2.2)
			case VarianceChannel:
				c = b.Pixels[y*b.W+x].Variance()
			case StandardDeviationChannel:
				c = b.Pixels[y*b.W+x].StandardDeviation()
			case SamplesChannel:
				p := float64(b.Pixels[y*b.W+x].Samples) / maxSamples
				c = Color{p, p, p}
			case DistanceChannel:
				// c = Color{float64(x) / 100, float64(x) / 100, float64(x) / 100}
				c = b.Distance(x, y)
			case NormalChannel:
				// c = Color{float64(x) / 100, float64(x) / 100, float64(x) / 100}
				c = b.Normal(x, y)
			}
			result.SetRGBA64(x, y, c.RGBA64())
		}
	}
	return result
}
