package mahjong

import (
	"testing"
)

func makeTile(bt BaseTile, id int) *Tile {
	return &Tile{Tile: bt, RedDora: false, ID: id}
}

func TestCheckKokushi_Valid(t *testing.T) {
	p := NewPlayer(East, true)
	s := &ScoreCounter{Player: p}
	// 13 terminal/honor each at least one + one duplicate
	tiles := []BaseTile{_1m, _9m, _1s, _9s, _1p, _9p, _1z, _2z, _3z, _4z, _5z, _6z, _7z, _1m}
	s.Tiles = tiles
	if !s.CheckKokushi() {
		t.Fatalf("expected Kokushi to be true")
	}
}

func TestCheckKokushi_Invalid(t *testing.T) {
	p := NewPlayer(East, true)
	s := &ScoreCounter{Player: p}
	// contains a non-terminal tile
	tiles := []BaseTile{_1m, _9m, _2m, _9s, _1p, _9p, _1z, _2z, _3z, _4z, _5z, _6z, _7z, _1m}
	s.Tiles = tiles
	if s.CheckKokushi() {
		t.Fatalf("expected Kokushi to be false")
	}
}

func TestCheckChuren_Valid(t *testing.T) {
	p := NewPlayer(East, true)
	s := &ScoreCounter{Player: p}
	// Construct a valid Nine Gates (1 and 9 duplicated enough to satisfy implementation)
	// Use man tiles: 1m x3, 2-8 x1, 9m x4 -> total 14
	tiles := make([]BaseTile, 0, 14)
	for i := 0; i < 3; i++ {
		tiles = append(tiles, _1m)
	}
	for i := _2m; i <= _8m; i++ {
		tiles = append(tiles, i)
	}
	for i := 0; i < 4; i++ {
		tiles = append(tiles, _9m)
	}
	s.Tiles = tiles
	if !s.CheckChuren() {
		t.Fatalf("expected Churen to be true")
	}
}

func TestCheckChuren_Invalid(t *testing.T) {
	p := NewPlayer(East, true)
	s := &ScoreCounter{Player: p}
	// Not all same suit
	tiles := []BaseTile{_1m, _1m, _2m, _3m, _4m, _5m, _6m, _7m, _8m, _9m, _9m, _9m, _1p, _2p}
	s.Tiles = tiles
	if s.CheckChuren() {
		t.Fatalf("expected Churen to be false")
	}
}

func TestCheckTenhou_And_Chihou(t *testing.T) {
	pOya := NewPlayer(East, true)
	pOya.FirstRound = true
	s := &ScoreCounter{Player: pOya}
	if !s.CheckTenhou() {
		t.Fatalf("expected Tenhou true for oya first round")
	}

	pChild := NewPlayer(East, false)
	pChild.FirstRound = true
	s2 := &ScoreCounter{Player: pChild}
	if !s2.CheckChihou() {
		t.Fatalf("expected Chihou true for non-oya first round")
	}
}

func TestCheckRinshan_And_Chankan(t *testing.T) {
	// Rinshan: LastAction is Kan/AnKan/KaKan and player has win tile in hand
	table := NewTable()
	p := NewPlayer(East, true)
	// make win tile
	win := _5m
	p.Hand = []*Tile{makeTile(win, 1)}
	s := &ScoreCounter{Player: p, Table: table, Tiles: []BaseTile{win}, WinTile: win}
	// simulate Kan
	table.LastAction = Kan
	if !s.CheckRinshan() {
		t.Fatalf("expected Rinshan true when last action is Kan and win in hand")
	}

	// Chankan: LastAction == KaKan and win is NOT in hand
	table.LastAction = KaKan
	p2 := NewPlayer(East, false)
	p2.Hand = []*Tile{makeTile(_1m, 2)}
	s2 := &ScoreCounter{Player: p2, Table: table, Tiles: []BaseTile{_5m}, WinTile: _5m}
	if !s2.CheckChankan() {
		t.Fatalf("expected Chankan true when last action is KaKan and win not in hand")
	}
}

func TestCheckHaitei_Hotei_Menzentsumo_Ippatsu_Daburu(t *testing.T) {
	table := NewTable()
	// make Yama length 14 so GetRemainTile == 0
	table.Yama = make([]*Tile, 14)
	p := NewPlayer(East, true)
	win := _3p
	p.Hand = []*Tile{makeTile(win, 10)}
	s := &ScoreCounter{Player: p, Table: table, Tiles: []BaseTile{win}, WinTile: win}
	if !s.CheckHaitei() {
		t.Fatalf("expected Haitei true when win in hand and no remain tiles")
	}

	// Hotei: win is not in hand (ron on last tile)
	p2 := NewPlayer(East, false)
	p2.Hand = []*Tile{makeTile(_1m, 11)}
	s2 := &ScoreCounter{Player: p2, Table: table, Tiles: []BaseTile{_3p}, WinTile: _3p}
	if !s2.CheckHotei() {
		t.Fatalf("expected Hotei true when ron on last tile")
	}

	// Menzentsumo
	p3 := NewPlayer(East, true)
	p3.Menzen = true
	s3 := &ScoreCounter{Player: p3}
	if !s3.CheckMenzentsumo() {
		t.Fatalf("expected Menzentsumo true when Menzen")
	}

	// Ippatsu & Double Riichi
	p4 := NewPlayer(East, false)
	p4.Ippatsu = true
	p4.DoubleRiichi = true
	s4 := &ScoreCounter{Player: p4}
	if !s4.CheckIppatsu() {
		t.Fatalf("expected Ippatsu true when Ippatsu flag set")
	}
	if !s4.CheckDabururiichi() {
		t.Fatalf("expected Dabururiichi true when DoubleRiichi flag set")
	}
}
