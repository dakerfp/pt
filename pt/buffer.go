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
	HitsChannel
	DiffuseColorChannel
)

type Pixel struct {
	color        ColorDistribution
	normal       VectorDistribution
	dist         FloatDistribution
	hits         FloatDistribution
	albedo       ColorDistribution
	bump         FloatDistribution
	emmitance    FloatDistribution
	index        FloatDistribution
	gloss        FloatDistribution
	tint         FloatDistribution
	reflectivity FloatDistribution
	transparent  FloatDistribution
}

func (p *Pixel) AddSample(sample Color) {
	p.color.AddSample(sample)
}

func (p *Pixel) AddSampleFeature(sample Features) {
	p.color.AddSample(sample.Color)
	p.normal.AddSample(sample.Normal)
	p.dist.AddSample(sample.Distance)
	hits := 0.0
	if sample.Hits {
		hits = 1.0
	}
	p.hits.AddSample(hits)
	mat := sample.Material
	p.albedo.AddSample(mat.Color)
	p.bump.AddSample(mat.BumpMultiplier)
	p.emmitance.AddSample(mat.Emittance)
	p.index.AddSample(mat.Index)
	p.gloss.AddSample(mat.Gloss)
	p.tint.AddSample(mat.Tint)
	p.reflectivity.AddSample(mat.Reflectivity)
	transp := 0.0
	if mat.Transparent {
		transp = 1.0
	}
	p.transparent.AddSample(transp)
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

const FeatureRawSize = 38

func (p *Pixel) Raw() []float64 {
	col := p.color.Avg()
	norm := p.normal.Avg()
	albedo := p.albedo.Avg()
	colVar := p.color.Variance()
	normVar := p.normal.Variance()
	albedoVar := p.albedo.Variance()
	distVar := p.dist.Variance()
	gamma := col.Pow(1 / 2.2)
	return []float64{
		// Primary features
		col.R,
		col.G,
		col.B,
		norm.X,
		norm.Y,
		norm.Z,
		p.dist.Avg(),
		albedo.R,
		albedo.G,
		albedo.B,
		gamma.R,
		gamma.G,
		gamma.B,
		p.bump.Avg(),
		p.emmitance.Avg(),
		p.index.Avg(),
		p.gloss.Avg(),
		p.tint.Avg(),
		p.reflectivity.Avg(),
		p.transparent.Avg(),
		// Secondary features
		colVar.R,
		colVar.G,
		colVar.B,
		normVar.X,
		normVar.Y,
		normVar.Z,
		distVar,
		albedoVar.R,
		albedoVar.G,
		albedoVar.B,
		p.bump.Variance(),
		p.emmitance.Variance(),
		p.index.Variance(),
		p.gloss.Variance(),
		p.tint.Variance(),
		p.reflectivity.Variance(),
		p.transparent.Variance(),

		// Other
		p.hits.Avg(),
	}
}

func (p *Pixel) Hits() float64 {
	return p.hits.M
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

func (b *Buffer) Pixel(x, y int) Pixel {
	return b.Pixels[y*b.W+x]
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
	return Gray((d - b.MinDist) / (b.MaxDist - b.MinDist))
}

func (b *Buffer) Hits(x, y int) Color {
	return Gray(b.Pixels[y*b.W+x].Hits())
}

func (b *Buffer) Normal(x, y int) Color {
	return b.Pixels[y*b.W+x].Normal()
}

func Gray(v float64) Color {
	return Color{v, v, v}
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
				c = Gray(float64(b.Pixels[y*b.W+x].color.N) / maxSamples)
			case DistanceChannel:
				c = b.Distance(x, y)
			case NormalChannel:
				c = b.Normal(x, y)
			case HitsChannel:
				c = b.Hits(x, y)
			case DiffuseColorChannel:
				c = b.Hits(x, y)
			}
			result.SetRGBA64(x, y, c.RGBA64())
		}
	}
	return result
}

func (b *Buffer) Raw() ([]int, []float64) {
	shape := []int{b.H, b.W, FeatureRawSize}
	result := make([]float64, FeatureRawSize*b.W*b.H)
	i := 0
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			p := b.Pixel(x, y)
			copy(result[i:i+FeatureRawSize], p.Raw())
			i += FeatureRawSize
		}
	}
	return shape, result
}
