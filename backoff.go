package backoff

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

var (
	_ Backoff = &Fixed{}
	_ Backoff = &Random{}
	_ Backoff = &Linear{}
	_ Backoff = &Polynomial{}
	_ Backoff = &Exponential{}
	_ Backoff = &Fibonacci{}
)

// NoLimit is the alias for MaxInt64 in [time.Duration].
// It is the actual limit of durations.
const NoLimit = time.Duration(math.MaxInt64)

// Backoff provides backoff algorithm.
type Backoff interface {
	// Attempt returns the n-th backoff duration.
	// Returned duration can be zero or positive.
	// n<=0 always returns zero.
	Attempt(n int) time.Duration
}

// NewFixed returns a new instance of [Fixed].
//   - interval: fixed value of backoff duration. (>=0)
//
// Backoff duration for n-th attempt will be:
//
//	d = interval
func NewFixed(interval time.Duration) (*Fixed, error) {
	if interval < 0 {
		return nil, rangeError("expect interval>=0. got " + interval.String())
	}
	return &Fixed{
		interval: interval,
	}, nil
}

// Fixed provides fixed-interval backoff strategy.
type Fixed struct {
	interval time.Duration
}

func (b *Fixed) Attempt(n int) time.Duration {
	if n <= 0 {
		return 0
	}
	return b.interval
}

// NewRandom returns a new instance of [Random].
//   - offset: offset for the backoff duration. (>=0)
//   - limit: maximum backoff duration. (>=offset)
//
// Backoff duration for n-th attempt will be:
//
//	d = offset + random(limit-offset)
func NewRandom(offset, limit time.Duration) (*Random, error) {
	if err := checkCommonParams(offset, limit, dummy, dummy); err != nil {
		return nil, err
	}
	return &Random{
		offset:   offset,
		deltaMax: float64(limit - offset),
	}, nil
}

// Random provides random backoff strategy.
type Random struct {
	offset   time.Duration
	deltaMax float64
}

func (b *Random) Attempt(n int) time.Duration {
	if n <= 0 {
		return 0
	}
	return b.offset + time.Duration(rand.Float64()*b.deltaMax)
}

// NewLinear returns a new instance of [Linear].
//   - offset: offset for the backoff duration. (>=0)
//   - limit: maximum backoff duration. (>=offset)
//   - coeff: coefficient for linear function. (>=0)
//   - jitter: parameter for random fluctuation. 1.0 for full jitter. (>=0.0, <=1.0)
//
// Backoff duration for n-th attempt with no jitter will be:
//
//	d = offset + coeff * (n)
func NewLinear(offset, limit, coeff time.Duration, jitter float64) (*Linear, error) {
	if err := checkCommonParams(offset, limit, coeff, jitter); err != nil {
		return nil, err
	}
	return &Linear{
		offset:   float64(offset),
		deltaMax: float64(limit - offset),
		coeff:    float64(coeff),
		jitter:   jitter,
	}, nil
}

// Linear provides linear backoff strategy.
type Linear struct {
	offset   float64
	deltaMax float64
	coeff    float64
	jitter   float64
}

func (b *Linear) Attempt(n int) time.Duration {
	if n <= 0 {
		return 0
	}
	delta := b.coeff * float64(n)
	delta = min(delta, b.deltaMax)
	return withJitter(b.jitter, b.offset, delta)
}

// NewPolynomial returns a new instance of [Polynomial].
//   - offset: offset for the backoff duration. (>=0)
//   - limit: maximum backoff duration. (>=offset)
//   - coeff: coefficient for polynomial function. (>=0)
//   - exponent: exponent of n. (>=0)
//   - jitter: parameter for random fluctuation. 1.0 for full jitter. (>=0.0, <=1.0)
//
// Backoff duration for n-th attempt with no jitter will be:
//
//	d = offset + coeff * n^exponent
func NewPolynomial(offset, limit, coeff time.Duration, exponent, jitter float64) (*Polynomial, error) {
	if err := checkCommonParams(offset, limit, coeff, jitter); err != nil {
		return nil, err
	}
	if exponent < 0 {
		return nil, rangeError("expect exponent>=0. got " + fmt.Sprint(exponent))
	}
	return &Polynomial{
		offset:   float64(offset),
		deltaMax: float64(limit - offset),
		coeff:    float64(coeff),
		exponent: exponent,
		jitter:   jitter,
	}, nil
}

// Polynomial provides polynomial backoff strategy.
type Polynomial struct {
	offset   float64
	deltaMax float64
	coeff    float64
	exponent float64
	jitter   float64
}

func (b Polynomial) Attempt(n int) time.Duration {
	if n <= 0 {
		return 0
	}
	pow := math.Pow(float64(n), b.exponent)
	if math.IsInf(pow, 1) || math.IsNaN(pow) { // https://github.com/golang/go/issues/7394
		pow = math.MaxFloat64
	}
	delta := b.coeff * pow
	delta = min(delta, b.deltaMax)
	return withJitter(b.jitter, b.offset, delta)
}

// NewExponential returns a new instance of [Exponential].
//   - offset: offset for the backoff duration. (>=0)
//   - limit: maximum backoff duration. (>=offset)
//   - coeff: coefficient for exponential function. (>=0)
//   - base: base number for exponential function. (>=1)
//   - jitter: parameter for random fluctuation. 1.0 for full jitter. (>=0.0, <=1.0)
//
// Backoff duration for n-th attempt with no jitter will be:
//
//	d = offset + coeff * base^n
func NewExponential(offset, limit, coeff time.Duration, base, jitter float64) (*Exponential, error) {
	if err := checkCommonParams(offset, limit, coeff, jitter); err != nil {
		return nil, err
	}
	if base < 1 {
		return nil, rangeError("expect base>=1. got " + fmt.Sprint(base))
	}
	return &Exponential{
		offset:   float64(offset),
		deltaMax: float64(limit - offset),
		coeff:    float64(coeff),
		base:     base,
		jitter:   jitter,
	}, nil
}

// Exponential provides exponential backoff strategy.
type Exponential struct {
	offset   float64
	deltaMax float64
	coeff    float64
	base     float64
	jitter   float64
}

func (b *Exponential) Attempt(n int) time.Duration {
	if n <= 0 {
		return 0
	}
	pow := math.Pow(b.base, float64(n))
	if math.IsInf(pow, 1) || math.IsNaN(pow) { // https://github.com/golang/go/issues/7394
		pow = math.MaxFloat64
	}
	delta := b.coeff * pow
	delta = min(delta, b.deltaMax)
	return withJitter(b.jitter, b.offset, delta)
}

// NewFibonacci returns a new instance of [Fibonacci].
//   - offset: offset for the backoff duration. (>=0)
//   - limit: maximum backoff duration. (>=offset)
//   - coeff: coefficient for fibonacci number. (>=0)
//   - jitter: parameter for random fluctuation. 1.0 for full jitter. (>=0.0, <=1.0)
//
// Backoff duration for n-th attempt with no jitter will be:
//
//	d = offset + coeff * fibonacci(n)
func NewFibonacci(offset, limit, coeff time.Duration, jitter float64) (*Fibonacci, error) {
	if err := checkCommonParams(offset, limit, coeff, jitter); err != nil {
		return nil, err
	}
	return &Fibonacci{
		offset:   float64(offset),
		deltaMax: float64(limit - offset),
		coeff:    float64(coeff),
		jitter:   jitter,
	}, nil
}

// Fibonacci provides fibonacci backoff strategy.
type Fibonacci struct {
	offset   float64
	deltaMax float64
	coeff    float64
	jitter   float64
}

func (b *Fibonacci) Attempt(n int) time.Duration {
	if n <= 0 {
		return 0
	}
	delta := b.coeff * fibonacci[min(n, 92)]
	delta = min(delta, b.deltaMax)
	return withJitter(b.jitter, b.offset, delta)
}

// dummy is the dummy value for checkCommonParams.
// It can be used like checkCommonParams(1,100,dummy,dummy).
const dummy = 1

// checkCommonParams checks common parameters.
// It returns error when
//   - offset < 0
//   - limit < offset
//   - coeff < 0
//   - jitter < 0 || jitter > 1
func checkCommonParams(offset, limit, coeff time.Duration, jitter float64) error {
	if offset < 0 {
		return rangeError("expect offset>=0. got " + offset.String())
	}
	if limit < offset {
		return rangeError("expect limit>=offset")
	}
	if coeff < 0 {
		return rangeError("expect coeff>=0. got " + coeff.String())
	}
	if jitter < 0 || jitter > 1 {
		return rangeError("expect 0<=jitter<=1. got " + fmt.Sprint(jitter))
	}
	return nil
}

// fibonacci is the list of fibonacci numbers.
var fibonacci = [...]float64{
	0:  0,
	1:  1,
	2:  1,
	3:  2,
	4:  3,
	5:  5,
	6:  8,
	7:  13,
	8:  21,
	9:  34,
	10: 55,
	11: 89,
	12: 144,
	13: 233,
	14: 377,
	15: 610,
	16: 987,
	17: 1597,
	18: 2584,
	19: 4181,
	20: 6765,
	21: 10946,
	22: 17711,
	23: 28657,
	24: 46368,
	25: 75025,
	26: 121393,
	27: 196418,
	28: 317811,
	29: 514229,
	30: 832040,
	31: 1346269,
	32: 2178309,
	33: 3524578,
	34: 5702887,
	35: 9227465,
	36: 14930352,
	37: 24157817,
	38: 39088169,
	39: 63245986,
	40: 102334155,
	41: 165580141,
	42: 267914296,
	43: 433494437,
	44: 701408733,
	45: 1134903170,
	46: 1836311903, // int32 limit
	47: 2971215073, // int32 overflow
	48: 4807526976,
	49: 7778742049,
	50: 12586269025,
	51: 20365011074,
	52: 32951280099,
	53: 53316291173,
	54: 86267571272,
	55: 139583862445,
	56: 225851433717,
	57: 365435296162,
	58: 591286729879,
	59: 956722026041,
	60: 1548008755920,
	61: 2504730781961,
	62: 4052739537881,
	63: 6557470319842,
	64: 10610209857723,
	65: 17167680177565,
	66: 27777890035288,
	67: 44945570212853,
	68: 72723460248141,
	69: 117669030460994,
	70: 190392490709135,
	71: 308061521170129,
	72: 498454011879264,
	73: 806515533049393,
	74: 1304969544928657,
	75: 2111485077978050,
	76: 3416454622906707,
	77: 5527939700884757,
	78: 8944394323791464,
	79: 14472334024676221,
	80: 23416728348467685,
	81: 37889062373143906,
	82: 61305790721611591,
	83: 99194853094755497,
	84: 160500643816367088,
	85: 259695496911122585,
	86: 420196140727489673,
	87: 679891637638612258,
	88: 1100087778366101931,
	89: 1779979416004714189,
	90: 2880067194370816120,
	91: 4660046610375530309,
	92: 7540113804746346429,  // int64 limit
	93: 12200160415121876738, // int64 overflow
}
