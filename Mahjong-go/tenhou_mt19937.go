package mahjong

// Mersenne Twister MT19937 port (minimal subset used by Tenhou shuffle)

const (
	mtN = 624
	mtM = 397
)

var (
	mt    [mtN]uint32
	mtIdx = mtN + 1
)

func initGenrand(s uint32) {
	mt[0] = s
	for i := 1; i < mtN; i++ {
		mt[i] = (1812433253*(mt[i-1]^(mt[i-1]>>30)) + uint32(i)) & 0xffffffff
	}
	mtIdx = mtN
}

func initByArray(initKey []uint32) {
	initGenrand(19650218)
	i := 1
	j := 0
	k := mtN
	if len(initKey) > mtN {
		k = len(initKey)
	}
	for ; k > 0; k-- {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 30)) * 1664525)) + initKey[j] + uint32(j)
		mt[i] &= 0xffffffff
		i++
		j++
		if i >= mtN {
			mt[0] = mt[mtN-1]
			i = 1
		}
		if j >= len(initKey) {
			j = 0
		}
	}
	for k = mtN - 1; k > 0; k-- {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 30)) * 1566083941)) - uint32(i)
		mt[i] &= 0xffffffff
		i++
		if i >= mtN {
			mt[0] = mt[mtN-1]
			i = 1
		}
	}
	mt[0] = 0x80000000
}

func genrandInt32() uint32 {
	var mag01 = [2]uint32{0x0, 0x9908b0df}
	if mtIdx >= mtN {
		var kk int
		if mtIdx == mtN+1 {
			initGenrand(5489)
		}
		for kk = 0; kk < mtN-mtM; kk++ {
			y := (mt[kk] & 0x80000000) | (mt[kk+1] & 0x7fffffff)
			mt[kk] = mt[kk+mtM] ^ (y >> 1) ^ mag01[y&0x1]
		}
		for ; kk < mtN-1; kk++ {
			y := (mt[kk] & 0x80000000) | (mt[kk+1] & 0x7fffffff)
			mt[kk] = mt[kk+(mtM-mtN)] ^ (y >> 1) ^ mag01[y&0x1]
		}
		y := (mt[mtN-1] & 0x80000000) | (mt[0] & 0x7fffffff)
		mt[mtN-1] = mt[mtM-1] ^ (y >> 1) ^ mag01[y&0x1]
		mtIdx = 0
	}
	y := mt[mtIdx]
	mtIdx++
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)
	return y
}
