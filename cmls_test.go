package count_min_log_sketch

import (
	"crypto/md5"
	"testing"

	"strconv"

	"github.com/seiflotfy/count-min-log"
	"github.com/stretchr/testify/assert"
)

func generateKey(numCookies int, numStrategies int) chan []byte {
	result := make(chan []byte, 1)

	go func() {
		for c := 0; c < numCookies; c++ {
			cookieHash := md5.Sum([]byte(strconv.Itoa(c)))
			for s := 0; s < numStrategies; s++ {
				strategyId := []byte(strconv.Itoa(s))

				key := make([]byte, 24)
				copy(key[0:16], cookieHash[:])
				copy(key[16:len(key)], strategyId)

				result <- key
			}
		}
		close(result)
	}()

	return result
}

func TestGenerateKey(t *testing.T) {
	results := [][]byte{
		{0xcf, 0xcd, 0x20, 0x84, 0x95, 0xd5, 0x65, 0xef, 0x66, 0xe7, 0xdf, 0xf9, 0xf9, 0x87, 0x64, 0xda, 0x30, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, // cookie: 0, strategy: 0
		{0xcf, 0xcd, 0x20, 0x84, 0x95, 0xd5, 0x65, 0xef, 0x66, 0xe7, 0xdf, 0xf9, 0xf9, 0x87, 0x64, 0xda, 0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, // cookie: 0, strategy: 1
		{0xc4, 0xca, 0x42, 0x38, 0xa0, 0xb9, 0x23, 0x82, 0xd, 0xcc, 0x50, 0x9a, 0x6f, 0x75, 0x84, 0x9b, 0x30, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},  // cookie: 1, strategy: 0
		{0xc4, 0xca, 0x42, 0x38, 0xa0, 0xb9, 0x23, 0x82, 0xd, 0xcc, 0x50, 0x9a, 0x6f, 0x75, 0x84, 0x9b, 0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},  // cookie: 1, strategy: 1
	}

	i := 0
	for got := range generateKey(2, 2) {
		expected := results[i]
		assert.Equal(t, expected, got)
		i += 1
	}

	assert.Equal(t, len(results), 4)
}

func TestCmlSketch8(t *testing.T) {
	key := []byte{0xcf, 0xcd, 0x20, 0x84, 0x95, 0xd5, 0x65, 0xef, 0x66, 0xe7, 0xdf, 0xf9, 0xf9, 0x87, 0x64, 0xda, 0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	sketch := cmlSketch8(t, 100, 100)
	assert.Equal(t, 5, int(sketch.Frequency(key)))
}

func cmlSketch8(tb testing.TB, numCookies int, numStrategies int) *cml.Sketch8 {
	// Epsilon = how much error is added to our count with each item we add to the sketch
	// Delta = with what probability do we allow the count estimate to be outside of our epsilon error rate
	//	sketch, _ := cml.NewSketch8ForEpsilonDelta(0.01, 0.05)
	sketch, _ := cml.NewDefaultSketch8()
	for key := range generateKey(numCookies, numStrategies) {
		sketch.IncreaseCount(key)
		sketch.IncreaseCount(key)
		sketch.IncreaseCount(key)
		sketch.IncreaseCount(key)
		sketch.IncreaseCount(key)
	}

	return sketch
}

func cmlSketch16(tb testing.TB, numCookies int, numStrategies int) *cml.Sketch16 {
	// Epsilon = how much error is added to our count with each item we add to the sketch
	// Delta = with what probability do we allow the count estimate to be outside of our epsilon error rate
	//	sketch, _ := cml.NewSketch16ForEpsilonDelta(0.01, 0.05)
	sketch, _ := cml.NewDefaultSketch16()
	for key := range generateKey(numCookies, numStrategies) {
		sketch.IncreaseCount(key)
		sketch.IncreaseCount(key)
		sketch.IncreaseCount(key)
		sketch.IncreaseCount(key)
		sketch.IncreaseCount(key)
	}

	return sketch
}

func BenchmarkBaseCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numCookies := 500000000
		numStrategies := 1
		counter := make(map[string]uint)
		for key := range generateKey(numCookies, numStrategies) {
			counter[string(key)] = 1
		}
	}
}

func BenchmarkCMLSketch8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numCookies := 10000
		numStrategies := 100
		cmlSketch8(b, numCookies, numStrategies)
	}
}

func BenchmarkCMLSketch16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numCookies := 10000
		numStrategies := 100
		cmlSketch16(b, numCookies, numStrategies)
	}
}
