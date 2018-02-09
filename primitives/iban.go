package primitives

import (
	"log"
	"math/big"
)

const IBANSize = 34

type IBAN [IBANSize]byte

func NewIBAN(b []byte) *IBAN {
	iban := new(IBAN)
	copy(iban[IBANSize-4:], b[:4])
	copy(iban[:IBANSize-4], b[4:])

	log.Println(string(iban[:]))

	for i := IBANSize - 1; i >= 0; i-- {
		if (iban[i] >= 65) && (iban[i] <= 90) {
			iban[i] = iban[i] - 65 + 10
		}
	}

	log.Println(iban[:])

	return iban
}

func (iban *IBAN) CalculateChecksum(n math.Big) {

}

func (iban *IBAN) ConvertToInteger() *math.Big {

}

func (iban *IBAN) String() string {
	return string(iban[:])
}
