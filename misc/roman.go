package misc

import (
	"errors"
	"fmt"
)

/*

  File:    roman.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Roman Numerals.
*/

// Paragraph generates a string useful for paragraph rows.
//goland:noinspection ALL
func Paragraph(num int) (string, error) {
	//goland:noinspection SpellCheckingInspection
	var tensCode = [...]string{"x", "xx", "xxx", "xl", "l", "lx", "lxx", "lxxx", "xc"}
	var onesCode = [...]string{"i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix"}
	if num > 99 {
		return fmt.Sprintf("*%5d.", num), errors.New("Invalid Roman")
	}
	s := ""
	tens := num / 10
	ones := num % 10
	if tens > 0 {
		s += tensCode[tens-1]
	}
	if ones > 0 {
		s += onesCode[ones-1]
	}
	return fmt.Sprintf("%6s.", s), nil

}

// Roman generates a string of roman numerals - up to 3999.
// 4000 and greater would require unicode characters 216x thru 218x.
//goland:noinspection ALL
func Roman(num int) (string, error) {
	var thousands = [...]string{"m", "mm", "mmm"}
	var hundreds = [...]string{"c", "cc", "ccc", "cd", "d", "dc", "dcc", "dccc", "cm"}
	var tens = [...]string{"x", "xx", "xxx", "xl", "l", "lx", "lxx", "lxxx", "xc"}
	var ones = [...]string{"i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix"}
	if num > 3999 {
		return fmt.Sprintf("*%d*", num), errors.New("Invalid Roman")
	}
	s := ""
	k := num / 1000
	h := num % 1000 / 100
	t := num % 100 / 10
	o := num % 10
	//	log.Printf("     %d %d %d %d\n", k, h, t, o)
	if k > 0 {
		s += thousands[k-1]
	}
	if h > 0 {
		s += hundreds[h-1]
	}
	if t > 0 {
		s += tens[t-1]
	}
	if o > 0 {
		s += ones[o-1]
	}
	return s, nil
}

/*
Table of Roman numerals in Unicode:
Value 	1	2	3	4	5	6	7	8	9	10	11	12	50	100	500	1,000
U+216x	Ⅰ	Ⅱ	Ⅲ	Ⅳ	Ⅴ	Ⅵ	Ⅶ	Ⅷ	Ⅸ	Ⅹ	Ⅺ	Ⅻ	Ⅼ	 Ⅽ	 Ⅾ	  Ⅿ
U+217x	ⅰ	ⅱ	ⅲ	ⅳ	ⅴ	ⅵ	ⅶ	ⅷ	ⅸ	ⅹ	ⅺ	ⅻ	ⅼ	 ⅽ	 ⅾ	  ⅿ

2180 ↀ ROMAN NUMERAL ONE THOUSAND
2181 ↁ ROMAN NUMERAL FIVE THOUSAND
2182 ↂ ROMAN NUMERAL TEN THOUSAND
2183 Ↄ ROMAN NUMERAL REVERSED ONE HUNDRED

2187 ↇ ROMAN NUMERAL FIFTY THOUSAND
2188 ↈ ROMAN NUMERAL ONE HUNDRED THOUSAND

1100    MC
9,999   IXCMXCIX      Bar over IX
9,997	IXCMXCVII
9,998	IXCMXCVIII
10,000	X             Bar over X
10,001	XI
*/
