package bitvec

import "testing"

func BenchmarkConstruct(b *testing.B) {
	b.Run("BitVec", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = NewBitVec(100, 2)
		}
	})

	b.Run("DiBit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewDiBit(100)
		}
	})
}

func BenchmarkSet(b *testing.B) {

	b.Run("BitVec", func(b *testing.B) {
		vec, _ := NewBitVec(100, 2)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = vec.Set(79, 3)
		}
	})

	b.Run("DiBit", func(b *testing.B) {
		vec := NewDiBit(100)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = vec.Set(79, 3)
		}
	})

}
func BenchmarkUnset(b *testing.B) {

	b.Run("BitVec", func(b *testing.B) {
		vec, _ := NewBitVec(100, 2)
		_ = vec.Set(79, 3)
		_ = vec.Set(81, 1)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = vec.Unset(81)
		}
	})

	b.Run("DiBit", func(b *testing.B) {
		vec := NewDiBit(100)
		_ = vec.Set(79, 3)
		_ = vec.Set(81, 1)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = vec.Unset(81)
		}
	})

}

func BenchmarkHas(b *testing.B) {

	b.Run("BitVec", func(b *testing.B) {
		vec, _ := NewBitVec(100, 2)
		_ = vec.Set(79, 3)
		_ = vec.Set(81, 1)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = vec.Has(79, 3)
		}
	})

	b.Run("DiBit", func(b *testing.B) {
		vec := NewDiBit(100)
		_ = vec.Set(79, 3)
		_ = vec.Set(81, 1)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = vec.Has(79, 3)
		}
	})

}

func BenchmarkState(b *testing.B) {

	b.Run("BitVec", func(b *testing.B) {
		vec, _ := NewBitVec(100, 2)
		_ = vec.Set(79, 3)
		_ = vec.Set(81, 1)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = vec.State(81)
		}
	})

	b.Run("DiBit", func(b *testing.B) {
		vec := NewDiBit(100)
		_ = vec.Set(79, 3)
		_ = vec.Set(81, 1)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = vec.State(81)
		}
	})

}
