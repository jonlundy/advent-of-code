package aoc

type Vector struct {
	Offset Point
	Scale  int
}

func (v Vector) Point() Point {
	return v.Offset.Scale(v.Scale)
}

type Point [2]int

func (p Point) Add(a Point) Point {
	return Point{p[0] + a[0], p[1] + a[1]}
}
func (p Point) Scale(m int) Point {
	return Point{p[0] * m, p[1] * m}
}
func (p Point) Less(b Point) bool {
	if p[0] != b[0] {
		return p[0] < b[0]
	}
	return p[1] < b[1]
}

func Transpose[T any](matrix [][]T) [][]T {
	rows, cols := len(matrix), len(matrix[0])

	m := make([][]T, cols)
	for i := range m {
		m[i] = make([]T, rows)
	}

	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			m[i][j] = matrix[j][i]
		}
	}
	return m
}

// NumPoints the number of the points inside an outline plus the number of points in the outline
func NumPoints(outline []Point, borderLength int) int {
	// shoelace - find the float area in a shape
	sum := 0
	for _, p := range Pairwise(outline) {
		row1, col1 := p[0][0], p[0][1]
		row2, col2 := p[1][0], p[1][1]

		sum += row1*col2 - row2*col1
	}
	area := sum / 2

	// pick's theorem - find the number of points in a shape given its area
	return (ABS(area) - borderLength/2 + 1) + borderLength
}

type Map[T any] [][]T

func (m *Map[T]) Get(p Point) (Point, T, bool) {
	var zero T
	if !m.Valid(p) {
		return [2]int{0, 0}, zero, false
	}

	return p, (*m)[p[0]][p[1]], true
}
func (m *Map[T]) Size() (int, int) {
	if m == nil || len(*m) == 0 {
		return 0, 0
	}
	return len(*m), len((*m)[0])
}
func (m *Map[T]) Valid(p Point) bool {
	rows, cols := m.Size()
	return p[0] >= 0 && p[0] < rows && p[1] >= 0 && p[1] < cols
}
