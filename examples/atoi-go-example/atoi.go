
// line 1 "atoi.rl"
// -*-go-*-
//
// Convert a string to an integer.
//
// To compile:
//
//   ragel -Z -T0 -o atoi.go atoi.rl
//   go build -o atoi atoi.go
//   ./atoi
//
// To show a diagram of your state machine:
//
//   ragel -V -Z -p -o atoi.dot atoi.rl
//   xdot atoi.dot
//

package main

import (
	"os"
	"fmt"
)


// line 28 "atoi.go"
var _atoi_actions []byte = []byte{
	0, 1, 0, 1, 1, 
}

var _atoi_key_offsets []byte = []byte{
	0, 0, 4, 6, 9, 
}

var _atoi_trans_keys []byte = []byte{
	43, 45, 48, 57, 48, 57, 10, 48, 
	57, 
}

var _atoi_single_lengths []byte = []byte{
	0, 2, 0, 1, 0, 
}

var _atoi_range_lengths []byte = []byte{
	0, 1, 1, 1, 0, 
}

var _atoi_index_offsets []byte = []byte{
	0, 0, 4, 6, 9, 
}

var _atoi_trans_targs []byte = []byte{
	2, 2, 3, 0, 3, 0, 4, 3, 
	0, 0, 
}

var _atoi_trans_actions []byte = []byte{
	0, 1, 3, 0, 3, 0, 0, 3, 
	0, 0, 
}

const atoi_start int = 1
const atoi_first_final int = 3
const atoi_error int = 0

const atoi_en_main int = 1


// line 27 "atoi.rl"


func atoi(data string) (val int) {
	cs, p, pe := 0, 0, len(data)
	neg := false

	
// line 79 "atoi.go"
	{
	cs = atoi_start
	}

// line 84 "atoi.go"
	{
	var _klen int
	var _trans int
	var _acts int
	var _nacts uint
	var _keys int
	if p == pe {
		goto _test_eof
	}
	if cs == 0 {
		goto _out
	}
_resume:
	_keys = int(_atoi_key_offsets[cs])
	_trans = int(_atoi_index_offsets[cs])

	_klen = int(_atoi_single_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case data[p] < _atoi_trans_keys[_mid]:
				_upper = _mid - 1
			case data[p] > _atoi_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_atoi_range_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case data[p] < _atoi_trans_keys[_mid]:
				_upper = _mid - 2
			case data[p] > _atoi_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
	cs = int(_atoi_trans_targs[_trans])

	if _atoi_trans_actions[_trans] == 0 {
		goto _again
	}

	_acts = int(_atoi_trans_actions[_trans])
	_nacts = uint(_atoi_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _atoi_actions[_acts-1] {
		case 0:
// line 34 "atoi.rl"

 neg = true 
		case 1:
// line 35 "atoi.rl"

 val = val * 10 + (int(data[p]) - '0') 
// line 170 "atoi.go"
		}
	}

_again:
	if cs == 0 {
		goto _out
	}
	p++
	if p != pe {
		goto _resume
	}
	_test_eof: {}
	_out: {}
	}

// line 44 "atoi.rl"


	if neg {
		val = -1 * val;
	}

	if cs < atoi_first_final {
		fmt.Println("atoi: there was an error:", cs, "<", atoi_first_final)
		fmt.Println(data)
		for i := 0; i < p; i++ {
			fmt.Print(" ")
		}
		fmt.Println("^")
	}

	return val
}

//////////////////////////////////////////////////////////////////////

type atoiTest struct {
	s string
	v int
}

var atoiTests = []atoiTest{
	atoiTest{"7", 7},
	atoiTest{"666", 666},
	atoiTest{"-666", -666},
	atoiTest{"+666", 666},
	atoiTest{"1234567890", 1234567890},
	atoiTest{"+1234567890\n", 1234567890},
	atoiTest{"+ 1234567890", 1234567890}, // i will fail
}

func main() {
	res := 0
	for _, test := range atoiTests {
		res := atoi(test.s)
		if res != test.v {
			fmt.Fprintf(os.Stderr, "FAIL atoi(%#v) != %#v\n", test.s, test.v)
			res = 1
		}
	}
	os.Exit(res)
}
