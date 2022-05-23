package helper

import (
	"math/rand"
	"strconv"
	"time"
)

var cpfs = make(map[string]struct{})

func GenerateCpf() string {
	cpf := ""

	for {
		rawCpf := generateCpfNumbers()
		rawCpf = calcVerifyingDigit(rawCpf, 10)
		rawCpf = calcVerifyingDigit(rawCpf, 11)

		cpf = cpfToString(rawCpf)
		if _, ok := cpfs[cpf]; ok {
			continue
		}
		cpfs[cpf] = struct{}{}
		break
	}

	return cpf
}

func generateCpfNumbers() []int {
	rand.Seed(time.Now().UnixNano())
	cpf := make([]int, 9)
	count := 0

	for count < 9 {
		cpf[count] = rand.Intn(10)
		count += 1
	}

	return cpf
}

func calcVerifyingDigit(cpf []int, initScore int) []int {
	sum := 0
	for _, number := range cpf {
		sum += (number * initScore)
		initScore -= 1
	}

	rest := sum % 11
	if rest < 2 {
		cpf = append(cpf, 0)
	} else {
		digit := 11 - rest
		cpf = append(cpf, digit)
	}

	return cpf
}

func cpfToString(cpf []int) string {
	resp := ""
	for _, number := range cpf {
		resp += strconv.Itoa(number)
	}
	return resp
}
