package api

import (
	"encoding/xml"
	"fmt"
	"math/big"
	"net/http"

	"utile.space/api/utils"
)

/*
The following code is somehow a direct translation of the Python code provided on the wikipedia page below:
https://en.wikipedia.org/wiki/Chudnovsky_algorithm#Python_code
*/
const floatPrecision uint = 100000 // 100K

func newBigFloat(v float64) *big.Float {
	return big.NewFloat(v).SetPrec(floatPrecision)
}

func binarySplit(a int, b int) (*big.Float, *big.Float, *big.Float) {
	var Pab, Qab, Rab *big.Float

	o := newBigFloat(1)

	A := newBigFloat(float64(a))

	if b == a+1 {
		i := newBigFloat(10939058860032000)
		j := newBigFloat(545140134)
		k := newBigFloat(13591409)
		e := newBigFloat(0).Set(o.Add(o.Mul(newBigFloat(6), A), newBigFloat(-5)))
		f := newBigFloat(0).Set(o.Add(o.Mul(newBigFloat(2), A), newBigFloat(-1)))
		g := newBigFloat(0).Set(o.Add(o.Mul(newBigFloat(6), A), newBigFloat(-1)))

		Pab = newBigFloat(1)
		Pab = Pab.Mul(Pab, e)
		Pab = Pab.Mul(Pab, f)
		Pab = Pab.Mul(Pab, g)
		Pab = Pab.Mul(Pab, newBigFloat(-1))

		Qab = newBigFloat(1)
		Qab = Qab.Mul(Qab, i)
		Qab = Qab.Mul(Qab, cube(A))

		Rab = newBigFloat(1)
		Rab = Rab.Mul(Rab, Pab)

		j.Mul(j, A)
		j.Add(j, k)
		Rab = Rab.Mul(Rab, j)
	} else {
		m := (a + b) / 2
		Pam, Qam, Ram := binarySplit(a, m)
		Pmb, Qmb, Rmb := binarySplit(m, b)

		o1 := newBigFloat(1)
		o2 := newBigFloat(1)
		o3 := newBigFloat(1)
		o4 := newBigFloat(1)

		Pab = o1.Mul(Pam, Pmb)
		Qab = o2.Mul(Qam, Qmb)

		o3.Mul(Qmb, Ram)
		o4.Mul(Pam, Rmb)
		Rab = o3.Add(o3, o4)
	}
	return Pab, Qab, Rab
}

func cube(v *big.Float) *big.Float {
	result := newBigFloat(1)
	result = result.Mul(result, v)
	result = result.Mul(result, v)
	result = result.Mul(result, v)
	return result
}

// chudnovsky computes π using the Chudnovsky algorithm
func chudnovsky(n int) *big.Float {
	_, Q1n, R1n := binarySplit(1, n)
	k := newBigFloat(426880.0)
	l := newBigFloat(1).Sqrt(newBigFloat(10005.0))
	m := newBigFloat(13591409.0)

	deno := newBigFloat(1).Mul(newBigFloat(1).Mul(k, l), Q1n)
	divi := newBigFloat(1).Add(newBigFloat(1).Mul(m, Q1n), R1n)

	return newBigFloat(1).Quo(deno, divi)
}

func chudnovskyTau(n int) *big.Float {
	pi := chudnovsky(n)
	pi.Mul(pi, newBigFloat(2.0))

	return pi
}

// @Summary		Pi Value
// @Description	Calculate Pi value up to 10K decimals
// @Tags			pi
// @Produce		json,xml,application/yaml,plain
// @Success		200	{object}	BigNumberResult
// @Router			/pi [get]
func CalculatePi(w http.ResponseWriter, r *http.Request) {
	pi := chudnovsky(10000)

	var answer BigNumberResult
	answer.Name = "Pi"
	answer.Value = fmt.Sprintf("%.10000f", pi)

	utils.Output(w, r.Header["Accept"], answer, answer.Value)
}

// @Summary		Tau Value
// @Description	Calculate Tau value up to 10K decimals
// @Tags			pi
// @Produce		json,xml,application/yaml,plain
// @Success		200	{object}	BigNumberResult
// @Router			/tau [get]
func CalculateTau(w http.ResponseWriter, r *http.Request) {
	tau := chudnovskyTau(10000)

	var answer BigNumberResult
	answer.Name = "Tau"
	answer.Value = fmt.Sprintf("%.10000f", tau)

	utils.Output(w, r.Header["Accept"], answer, answer.Value)
}

type BigNumberResult struct {
	XMLName xml.Name `json:"-" xml:"bignumber" yaml:"-"`
	Name    string   `json:"name" xml:"name" yaml:"name"`
	Value   string   `json:"value" xml:"value" yaml:"value"`
}
