package ga_test

import (
	"algorithm/ga"
	"math"
	"math/rand"
	"testing"
)

type demo struct {
	x, y float64
}

func (d *demo) Fitness() float64 {
	if f := 1 / (d.x*d.x + d.y*d.y + 1); f > 0 {
		return f
	} else {
		return 0
	}
}

func (d *demo) Cross(ind ga.Individual) ga.Individual {
	a := ind.(*demo)
	return &demo{
		x: 0.5 * (d.x + a.x),
		y: 0.5 * (d.y + a.y),
	}
}

func (d *demo) Mutate() ga.Individual {
	if rand.Float64() < 0.1 {
		d.x = 10*rand.Float64() - 5
		d.y = 10*rand.Float64() - 5
	}
	return d
}

func TestGA(t *testing.T) {
	size := 100
	inds := make([]ga.Individual, size)
	for i := 0; i < size; i++ {
		inds[i] = &demo{
			x: 10*rand.Float64() - 5,
			y: 10*rand.Float64() - 5,
		}
	}
	p := ga.New(inds)
	p.Evolve(50, 1000)
	if math.Abs(p.Fitness()-1) > 1e-3 {
		t.Fail()
	}
}
