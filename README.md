go_mininest
===========

This is a port from C to Go of example code for
John Skilling's Nested Sampling algorithm.

The original code is at http://www.inference.phy.cam.ac.uk/bayesys/.

# Installing
## From local source
    $ go install

## From github
    $ go get github.com/ijt/go_mininest

# Running Examples
    $ go run examples/lighthouse/lighthouse.go
    2012/04/26 11:10:24 # iterates = 1000
    2012/04/26 11:10:24 Evidence: ln(Z) = -159.288 +- 0.04767
    2012/04/26 11:10:24 Information: H = 2.27284 nats = 3.27902 bits
    2012/04/26 11:10:24 mean(x) = 1.21452, stddev(x) = 0.13085
    2012/04/26 11:10:24 mean(y) = 0.95857, stddev(y) = 0.13772

