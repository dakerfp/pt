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
	AlbedoChannel
)

type FloatDistribution struct {
	M, V float64
	N    int
}

func (fd *FloatDistribution) AddSample(v float64) {
	fd.N++
	if fd.N == 1 {
		fd.M = v
		fd.V = 0
		return
	}
	m := fd.M
	fd.M = v + (v-m)/float64(fd.N)
	fd.V = fd.V + (v-m)*(v-fd.M)
}

type ColorDistribution struct {
	M, V Color
	N    int
}

func (p *ColorDistribution) AddSample(v Color) {
	p.N++
	if p.N == 1 {
		p.M = v
		p.V = Black
		return
	}
	m := p.M
	p.M = p.M.Add(v.Sub(p.M).DivScalar(float64(p.N)))
	p.V = p.V.Add(v.Sub(m).Mul(v.Sub(p.M)))
}

func (p *ColorDistribution) Variance() Color {
	if p.N < 2 {
		return Black
	}
	return p.V.DivScalar(float64(p.N - 1))
}

type VectorDistribution struct {
	M, V Vector
	N    int
}

func (p *VectorDistribution) AddSample(v Vector) {
	p.N++
	if p.N == 1 {
		p.M = v
		p.V = Vector{}
		return
	}
	m := p.M
	p.M = p.M.Add(v.Sub(p.M).DivScalar(float64(p.N)))
	p.V = p.V.Add(v.Sub(m).Mul(v.Sub(p.M)))
}

type Pixel struct {
	color  ColorDistribution
	albedo ColorDistribution
	normal VectorDistribution
	dist   FloatDistribution
}

func (p *Pixel) AddSample(sample Color) {
	p.color.AddSample(sample)
}

func (p *Pixel) AddSampleFeature(sample Features) {
	p.color.AddSample(sample.Color)
	p.albedo.AddSample(sample.Albedo)
	p.normal.AddSample(sample.Normal)
	p.dist.AddSample(sample.Distance)
}

func (p *Pixel) Color() Color {
	return p.color.M
}

func (p *Pixel) Variance() Color {
	return p.color.Variance()
}

func (p *Pixel) Distance() float64 {
	return p.dist.M
}

func (p *Pixel) Albedo() Color {
	return p.albedo.M
}

func (p *Pixel) StandardDeviation() Color {
	return p.Variance().Pow(0.5)
}

func (p *Pixel) Normal() Color {
	n := p.normal.M
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
	return b.Pixels[y*b.W+x].color.N
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
			maxSamples = math.Max(maxSamples, float64(pixel.color.N))
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
				p := float64(b.Pixels[y*b.W+x].color.N) / maxSamples
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
