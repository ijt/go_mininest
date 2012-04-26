package mininest

import (
	"log"
	"math"
	"math/rand"
)

type Object interface {
	Explore(logLstar float64)
	LogL() float64 // Log-likelihood
	Copy() Object
}

type Result struct {
	H float64
	LogZ float64
	LogWt float64
	Sample Object
}

func GoSampling(Obj []Object) chan Result {
	cOut := make(chan Result)
	go Sampling(Obj, cOut)
	return cOut
}

// Generates an infinite stream of samples from the posterior, given a slice of
// Objects sampled from the prior.
func Sampling(Obj []Object, cOut chan Result) {
	H := 0.0                            // Information, initially 0
	logZ := -math.MaxFloat64            // ln(Evidence Z, initially 0)
	n := len(Obj)

	// Outermost interval of prior mass
	logwidth := math.Log(1.0 - math.Exp(-1.0/float64(n))) // ln(width in prior mass)
	for {
		// Find worst object in collection, with Weight = width * Likelihood
		worst := 0
		for i, o := range(Obj) {
			if o.LogL() < Obj[worst].LogL() {
				worst = i
			}
		}
		logWt := logwidth + Obj[worst].LogL()
		// Update Evidence Z and Information H
		logZnew := plus(logZ, logWt) // ln(Likelihood constraint)
		H = (math.Exp(logWt-logZnew)*Obj[worst].LogL() +
			math.Exp(logZ-logZnew)*(H+logZ) - logZnew)
		logZ = logZnew
		// Posterior sample
		sample := Obj[worst].Copy()
		// Kill worst object in favour of copy of different survivor
		// don't kill if n is only 1
		logLstar := Obj[worst].LogL() // new likelihood constraint
		if n > 1 {
			for {
				// aCopy is the index of the object to copy.
				aCopy := int(float64(n)*rand.Float64()) % n // force 0 <= aCopy < n
				if aCopy != worst {
					Obj[worst] = Obj[aCopy] // overwrite worst object
					break
				}
			}
		}
		// Evolve copied object within constraint
		Obj[worst].Explore(logLstar)
		// Shrink interval
		logwidth -= 1.0 / float64(n)

		cOut <- Result { H, logZ, logWt, sample }
	}
}

func plus(x float64, y float64) float64 {
	if x > y {
		return x + math.Log(1+math.Exp(y-x))
	}
	return y + math.Log(1+math.Exp(x-y))
}

// Gathers the results from a channel into an array.
func GetArray(c chan Result, iterations int) []Result {
	results := make([]Result, 0, iterations)
	i := 0
	for r := range(c) {
		if i += 1; i > iterations { break }
		results = append(results, r)
	}
	// Exit with evidence Z, information H, and optional posterior Samples
	r := results[len(results)-1]
	log.Printf("# iterates = %d\n", len(results))
	log.Printf("Evidence: ln(Z) = %.3f +- %.5f\n", r.LogZ,
		math.Sqrt(float64(r.H)/float64(len(results))))
	log.Printf("Information: H = %.5f nats = %.5f bits\n", r.H, r.H/math.Log(2.))
	return results
}

