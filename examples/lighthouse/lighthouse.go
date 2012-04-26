package main

// "LIGHTHOUSE" NESTED SAMPLING APPLICATION
// (GNU General Public License software, (C) Sivia and Skilling 2006)
// Translated to Go by Issac Trotts (2012)
//
//              u=0                                 u=1
//               -------------------------------------
//          y=2 |:::::::::::::::::::::::::::::::::::::| v=1
//              |::::::::::::::::::::::LIGHT::::::::::|
//         north|::::::::::::::::::::::HOUSE::::::::::|
//              |:::::::::::::::::::::::::::::::::::::|
//              |:::::::::::::::::::::::::::::::::::::|
//          y=0 |:::::::::::::::::::::::::::::::::::::| v=0
// --*--------------*----*--------*-**--**--*-*-------------*--------
//             x=-2          coastline -.east      x=2
// Problem:
//  Lighthouse at (x,y) emitted n flashes observed at D[.] on coast.
// Inputs:
//  Prior(u)    is uniform (=1) over (0,1), mapped to x = 4*u - 2; and
//  Prior(v)    is uniform (=1) over (0,1), mapped to y = 2*v; so that
//  Position    is 2-dimensional -2 < x < 2, 0 < y < 2 with flat prior
//  Likelihood  is L(x,y) = PRODUCT[k] (y/pi) / ((D[k] - x)^2 + y^2)
// Outputs:
//  Evidence    is Z = INTEGRAL L(x,y) Prior(x,y) dxdy
//  Posterior   is P(x,y) = L(x,y) / Z estimating lighthouse position
//  Information is H = INTEGRAL P(x,y) math.Log(P(x,y)/Prior(x,y)) dxdy

import (
	"log"
	"math"
	"math/rand"
	"go_mininest"
)

type Object struct {
	u     float64 // uniform-prior controlling parameter for x
	v     float64 // uniform-prior controlling parameter for y
	x     float64 // Geographical easterly position of lighthouse
	y     float64 // Geographical northerly position of lighthouse
	logL  float64 // logLikelihood = ln Prob(data | position)
}

func (o *Object) Copy() go_mininest.Object {
	ret := new(Object)
	*ret = *o
	return ret
}

func (o *Object) LogL() float64 {
	return o.logL
}

// logLikelihood function
// Easterly position
// Northerly position
func logLhood(x float64, y float64) float64 {
	N := 64 // # arrival positions
	D := []float64{4.73, 0.45, -1.73, 1.09, 2.19, 0.12,
		1.31, 1.00, 1.32, 1.07, 0.86, -0.49, -2.59, 1.73, 2.11,
		1.61, 4.98, 1.71, 2.23, -57.20, 0.96, 1.25, -1.56, 2.45,
		1.19, 2.17, -10.66, 1.91, -4.16, 1.92, 0.10, 1.98, -2.51,
		5.55, -0.47, 1.91, 0.95, -0.78, -0.84, 1.72, -0.01, 1.48,
		2.70, 1.21, 4.41, -4.79, 1.33, 0.81, 0.20, 1.58, 1.29,
		16.19, 2.75, -2.38, -1.79, 6.50, -18.53, 0.72, 0.94, 3.64,
		1.94, -0.11, 1.57, 0.57} // up to N=64 data
	var k int   // data index
	logL := 0.0 // logLikelihood accumulator
	for k = 0; k < N; k++ {
		num := (y / 3.1416)
		denom := ((D[k]-x)*(D[k]-x) + y*y)
		logL += math.Log(num / denom)
	}
	return logL
}

// Set Object according to prior
func SampleFromPrior() *Object {
	Obj := new(Object)
	Obj.u = rand.Float64()       // uniform in (0,1)
	Obj.v = rand.Float64()       // uniform in (0,1)
	Obj.x = 4.0*Obj.u - 2.0 // map to x
	Obj.y = 2.0 * Obj.v     // map to y
	Obj.logL = logLhood(Obj.x, Obj.y)
	return Obj
}

// Evolve object within likelihood constraint
// logLstar: Likelihood constraint L > Lstar
func (Obj *Object) Explore(logLstar float64) {
	step := 0.1    // Initial guess suitable step-size in (0,1)
	m := 20        // MCMC counter (pre-judged # steps)
	accept := 0    // # MCMC acceptances
	reject := 0    // # MCMC rejections
	var Try Object // Trial object
	for ; m > 0; m-- {
		// Trial object
		Try.u = Obj.u + step*(2.*rand.Float64()-1.) // |move| < step
		Try.v = Obj.v + step*(2.*rand.Float64()-1.) // |move| < step
		Try.u -= math.Floor(Try.u)             // wraparound to stay within (0,1)
		Try.v -= math.Floor(Try.v)             // wraparound to stay within (0,1)
		Try.x = 4.0*Try.u - 2.0                // map to x
		Try.y = 2.0 * Try.v                    // map to y
		Try.logL = logLhood(Try.x, Try.y)      // trial likelihood value
		// Accept if and only if within hard likelihood constraint
		if Try.logL > logLstar {
			*Obj = Try
			accept++
		} else {
			reject++
		}
		// Refine step-size to let acceptance ratio converge around 50%
		if accept > reject {
			step *= math.Exp(1.0 / float64(accept))
		}
		if accept < reject {
			step /= math.Exp(1.0 / float64(reject))
		}
	}
}

// Posterior properties, here mean and stddev of x,y
// Objects defining posterior
// Evidence (= total weight = SUM[Samples] Weight)
func Results(Samples []go_mininest.Object, logZ float64) {
}

func plus(x float64, y float64) float64 {
	if x > y {
		return x + math.Log(1+math.Exp(y-x))
	}
	return y + math.Log(1+math.Exp(x-y))
}

func main() {
	// Sample objects from the prior
	var objects [100]go_mininest.Object
	for i := range objects {
		objects[i] = SampleFromPrior()
	}

	results := go_mininest.GetArray(go_mininest.GoSampling(objects[:]), 1000)

	// Summarize the results.
	logZ := results[len(results)-1].LogZ  // final estimate of log(Z)
	x := 0.0
	xx := 0.0 // 1st and 2nd moments of x
	y := 0.0
	yy := 0.0 // 1st and 2nd moments of y
	for _, result := range(results) {
		switch s := result.Sample.(type) {
		case *Object:
			w := math.Exp(result.LogWt - logZ)
			x += w * s.x
			xx += w * s.x * s.x
			y += w * s.y
			yy += w * s.y * s.y
		default:
			log.Fatal("Unexpected type %T\n", s)
		}
	}
	log.Printf("mean(x) = %.5f, stddev(x) = %.5f\n", x, math.Sqrt(xx-x*x))
	log.Printf("mean(y) = %.5f, stddev(y) = %.5f\n", y, math.Sqrt(yy-y*y))
}

