package aoc_test


func TestLCM(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.LCM([]int{}...), 0)
	is.Equal(aoc.LCM(5), 5)
	is.Equal(aoc.LCM(5, 3), 15)
	is.Equal(aoc.LCM(5, 3, 2), 30)
}

func TestPower2(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.Power2(0), 1)
	is.Equal(aoc.Power2(1), 2)
	is.Equal(aoc.Power2(2), 4)
}

func TestABS(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.ABS(1), 1)
	is.Equal(aoc.ABS(0), 0)
	is.Equal(aoc.ABS(-1), 1)
}