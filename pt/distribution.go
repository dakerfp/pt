package pt

type FloatDistribution struct {
	Acc, M, V float64
	N         int
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

func (fd *FloatDistribution) Avg() float64 {
	return fd.M
}

func (fd *FloatDistribution) Variance() float64 {
	if fd.N < 2 {
		return 0
	}
	return fd.V / float64(fd.N-1)
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

func (cd *ColorDistribution) Avg() Color {
	return cd.M
}

func (cd *ColorDistribution) Variance() Color {
	if cd.N < 2 {
		return Black
	}
	return cd.V.DivScalar(float64(cd.N - 1))
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

func (vd *VectorDistribution) Avg() Vector {
	return vd.M
}

func (p *VectorDistribution) Variance() Vector {
	if p.N < 2 {
		return Vector{}
	}
	return p.V.DivScalar(float64(p.N - 1))
}
