package ga

import (
	"math/rand"
	"time"
)

// 个体
type Individual interface {
	Fitness() float64
	Cross(ind Individual) Individual
	Mutate() Individual
}

// 种群
type Population struct {
	size    int          // 种群大小
	elite   Individual   // 精英
	fitness float64      // 精英适应度
	fitsum  float64      // 适应度和
	fs      []float64    // 个体适应度
	inds    []Individual // 个体
	_fs     []float64
	_inds   []Individual
	_rand   *rand.Rand
}

// 创建种群
func New(inds []Individual) *Population {
	p := new(Population)
	p.size = len(inds)
	if p.size < 2 {
		panic("个体数量至少为2") // TODO
	}
	p.elite = inds[0]
	p.fs = make([]float64, p.size)
	p.inds = make([]Individual, p.size)
	p._fs = make([]float64, p.size)
	p._inds = make([]Individual, p.size)
	copy(p.inds, inds)

	for i := 0; i < p.size; i++ {
		p.fs[i] = p.inds[i].Fitness()
		p.fitsum += p.fs[i]
		if p.fs[i] > p.fitness {
			p.elite, p.fitness = p.inds[i], p.fs[i]
		}
	}

	p._rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return p
}

// 种群规模
func (p *Population) Size() int {
	return p.size
}

// 精英
func (p *Population) Elite() Individual {
	return p.elite
}

// 精英适应度
func (p *Population) Fitness() float64 {
	return p.fitness
}

// 种群
func (p *Population) Inds() []Individual {
	return p.inds
}

// 下一代
func (p *Population) NextGen() {
	s := 0.0
	for i := 0; i < p.size; i++ {
		x, y := p.pick()
		z := x.Cross(y).Mutate()
		p._inds[i], p._fs[i] = z, z.Fitness()
		s += p._fs[i]
		if p._fs[i] > p.fitness {
			p.elite, p.fitness = p._inds[i], p._fs[i]
		}
	}
	p.inds, p._inds = p._inds, p.inds
	p.fs, p._fs = p._fs, p.fs
	p.fitsum = s
}

// 进化
func (p *Population) Evolve(n int, m int) int {
	k, ef := 0, p.Fitness()
	for i := 0; i < m; i++ {
		p.NextGen()
		if ef >= p.Fitness() {
			k += 1
		} else {
			k, ef = 0, p.Fitness()
		}
		if k >= n {
			return i + 1
		}
	}
	return m
}

// 选择两个个体
func (p *Population) pick() (Individual, Individual) {
	i, j := 0, 0

	for r := p.fitsum * p._rand.Float64(); i < p.size; i++ {
		if r < p.fs[i] {
			break
		}
		r -= p.fs[i]
	}
	if i >= p.size {
		i = p.size - 1
	}

	for r := (p.fitsum - p.fs[i]) * p._rand.Float64(); j < p.size; j++ {
		if j == i {
			continue
		}
		if r < p.fs[j] {
			break
		}
		r -= p.fs[j]
	}
	if j >= p.size {
		j = p.size - 1
	}

	return p.inds[i], p.inds[j]
}
