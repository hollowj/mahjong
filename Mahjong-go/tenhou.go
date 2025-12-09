package mahjong

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
)

// TenhouYamaFromSeed reproduces the C++ tenhou::tenhou_yama_from_seed behavior.
// Input: MT seed as base64 string. Output: yama []int (length 136)
func TenhouYamaFromSeed(mtseedB64 string) ([]int, error) {
	// decode base64 seed to bytes
	seedBytes, err := base64.StdEncoding.DecodeString(mtseedB64)
	if err != nil {
		return nil, err
	}
	if len(seedBytes) < 4*mtN {
		return nil, errors.New("seed length too short")
	}

	// build RTseed uint32 array by taking 4 bytes each and then reversing endian (same as C++ convertEndian)
	rtseed := make([]uint32, mtN)
	for i := 0; i < mtN; i++ {
		// assemble same byte order as C++: MTseed[4*i]<<24 | MTseed[4*i+1]<<16 ...
		v := uint32(seedBytes[4*i])<<24 | uint32(seedBytes[4*i+1])<<16 | uint32(seedBytes[4*i+2])<<8 | uint32(seedBytes[4*i+3])
		// convertEndian: reverse bytes
		b0 := byte((v >> 24) & 0xff)
		b1 := byte((v >> 16) & 0xff)
		b2 := byte((v >> 8) & 0xff)
		b3 := byte(v & 0xff)
		rtseed[i] = uint32(b3)<<24 | uint32(b2)<<16 | uint32(b1)<<8 | uint32(b0)
	}

	// initialize MT with RTseed
	initByArray(rtseed)

	// prepare src and rnd arrays (mirroring C++ sizes)
	rndDwords := (sha512.Size / 4) * 9 // 144
	srcLen := rndDwords * 2            // 288
	src := make([]uint32, srcLen)
	for i := 0; i < srcLen; i++ {
		src[i] = genrandInt32()
	}

	// compute rnd by hashing src in 9 blocks of 128 bytes
	rnd := make([]uint32, rndDwords)
	blocks := (rndDwords * 4) / sha512.Size // should be 9
	for bi := 0; bi < blocks; bi++ {
		// each block hashes 128 bytes => 32 uint32s
		chunk := make([]byte, sha512.Size*2)
		// fill chunk as little-endian uint32, matching C++ memory layout on little-endian machines
		for j := 0; j < sha512.Size*2/4; j++ {
			idx := bi*(sha512.Size*2/4) + j
			binary.LittleEndian.PutUint32(chunk[j*4:(j+1)*4], src[idx])
		}
		sum := sha512.Sum512(chunk)
		// place digest into rnd as little-endian uint32s
		for k := 0; k < sha512.Size/4; k++ {
			rnd[bi*(sha512.Size/4)+k] = binary.LittleEndian.Uint32(sum[k*4 : (k+1)*4])
		}
	}

	// shuffle yama
	yama := make([]int, 136)
	for i := 0; i < 136; i++ {
		yama[i] = i
	}
	for i := 0; i < 136-1; i++ {
		tmpIndex := i + int(rnd[i]%uint32(136-i))
		yama[i], yama[tmpIndex] = yama[tmpIndex], yama[i]
	}
	return yama, nil
}
