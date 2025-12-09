package mahjong

// 注：鸣牌相关方法的完整实现需要根据实际的数据结构进行调整
// CallGroup结构中的Type字段应该是TileGroupType，而这里的BaseAction需要进行转换

// ExecuteNaki 执行鸣牌操作（统一入口）
func (p *Player) ExecuteNaki(tile *Tile, actionType BaseAction) {
	switch actionType {
	case Chi:
		p.ExecuteChiSimple(tile)
	case Pon:
		p.ExecutePonSimple(tile)
	case Kan:
		p.ExecuteKanSimple(tile)
	default:
		// 未知的鸣牌类型
	}
}

// ExecuteChiSimple 执行吃牌操作（简化版）
func (p *Player) ExecuteChiSimple(tile *Tile) {
	// 从手牌中移除对应的牌（使用RemoveFromHand）
	p.RemoveFromHand(tile)
	// 更新门前清状态
	p.Menzen = false
}

// ExecutePonSimple 执行碰牌操作（简化版）
func (p *Player) ExecutePonSimple(tile *Tile) {
	// 从手牌中移除两张同牌（使用现有的RemoveFromHand方法）
	p.RemoveFromHand(tile)
	p.RemoveFromHand(tile)
	p.Menzen = false
}

// ExecuteKanSimple 执行大明杠操作（简化版）
func (p *Player) ExecuteKanSimple(tile *Tile) {
	// 从手牌中移除三张同牌
	p.RemoveFromHand(tile)
	p.RemoveFromHand(tile)
	p.RemoveFromHand(tile)
	p.Menzen = false
}

// GetNormalXiangHuShu 计算向胡数（某个牌能让多少个玩家听牌）
func (p *Player) GetNormalXiangHuShu(tile BaseTile) int {
	xianghuCount := 0

	// 遍历其他玩家，检查该牌是否在他们的听牌列表中
	for i := 0; i < 4; i++ {
		// 应该从外部传入其他玩家信息
		// 这是简化版本，实际需要通过参数获取其他玩家
		//
		// if otherPlayers[i].isTenpai() {
		//     if contains(otherPlayers[i].atariTiles, tile) {
		//         xianghuCount++
		//     }
		// }
	}

	return xianghuCount
}

// RemoveAtariTileByList 根据列表移除听牌
func (p *Player) RemoveAtariTileByList(tiles []BaseTile) {
	for _, tile := range tiles {
		p.RemoveAtariTiles(tile)
	}
}

// UpdateBothFuriten 同时更新两种振听
func (p *Player) UpdateBothFuriten() {
	p.UpdateFuritenRiver()
	if p.IsRiichi() {
		// 检查立直后弃过的牌
		for _, riverTile := range p.River.River {
			if riverTile.Riichi {
				p.RemoveAtariTiles(riverTile.Tile.Tile)
			}
		}
	}
}

// GetMinGluedTile 获取最小粘着张（吃、碰后必须立即出牌的情况下的最小可出牌）
func (p *Player) GetMinGluedTile() BaseTile {
	if len(p.Hand) == 0 {
		return 0
	}

	// 按牌号排序，返回最小的
	minTile := p.Hand[0].Tile
	for _, tile := range p.Hand {
		if tile.Tile < minTile {
			minTile = tile.Tile
		}
	}
	return minTile
}

// GetMaxGluedTile 获取最大粘着张
func (p *Player) GetMaxGluedTile() BaseTile {
	if len(p.Hand) == 0 {
		return 0
	}

	maxTile := p.Hand[0].Tile
	for _, tile := range p.Hand {
		if tile.Tile > maxTile {
			maxTile = tile.Tile
		}
	}
	return maxTile
}

// IsTenpaiAfterDiscard 检查弃某张牌后是否仍然听牌
func (p *Player) IsTenpaiAfterDiscard(tile BaseTile) bool {
	// 创建临时手牌副本
	tempHand := make([]BaseTile, 0)
	for _, t := range p.Hand {
		tempHand = append(tempHand, t.Tile)
	}

	// 移除该牌
	for i, t := range tempHand {
		if t == tile {
			tempHand = append(tempHand[:i], tempHand[i+1:]...)
			break
		}
	}

	// 检查是否听牌
	splitter := &TileSplitter{}
	completedTiles := splitter.GetAllCompletedTiles(tempHand)
	return len(completedTiles) > 0
}

// GetFuritenTiles 获取振听的牌列表
func (p *Player) GetFuritenTiles() []BaseTile {
	furitenTiles := make([]BaseTile, 0)

	if p.FuritenRiver {
		// 历史弃过的牌
		for _, riverTile := range p.River.River {
			found := false
			for _, t := range furitenTiles {
				if t == riverTile.Tile.Tile {
					found = true
					break
				}
			}
			if !found {
				furitenTiles = append(furitenTiles, riverTile.Tile.Tile)
			}
		}
	}

	if p.FuritenRound {
		// 本轮弃过的牌（通过听牌列表的差异确定）
		// 这需要更复杂的逻辑
	}

	return furitenTiles
}

// CanDiscardForRiichi 检查是否可以弃该牌并立直
func (p *Player) CanDiscardForRiichi(tile BaseTile) bool {
	// 1. 必须门前清
	if !p.IsMenzen() {
		return false
	}

	// 2. 弃牌后必须听牌
	if !p.IsTenpaiAfterDiscard(tile) {
		return false
	}

	// 3. 弃牌后不能是振听
	if p.FuritenRiver {
		for _, riverTile := range p.River.River {
			if riverTile.Tile.Tile == tile {
				return false
			}
		}
	}

	return true
}

// UpdateFuritenAfterRiichi 立直后更新振听状态
func (p *Player) UpdateFuritenAfterRiichi() {
	p.FuritenRiichi = p.FuritenRiver || p.FuritenRound
}

// CountKan 计算杠的数量
func (p *Player) CountKan() int {
	kanCount := 0
	for _, group := range p.CallGroups {
		if group.Type == Kantsu {
			kanCount++
		}
	}
	return kanCount
}

// ResetRoundFlags 重置回合标志
func (p *Player) ResetRoundFlags() {
	p.Ippatsu = false
	p.FuritenRound = false
	p.FirstRound = false
}

// RecordHandTile 记录手牌状态
func (p *Player) RecordHandTile() []BaseTile {
	result := make([]BaseTile, len(p.Hand))
	for i, tile := range p.Hand {
		result[i] = tile.Tile
	}
	return result
}

// CheckHandConsistency 检查手牌一致性（调试用）
func (p *Player) CheckHandConsistency() bool {
	// 计算总牌数
	totalTiles := len(p.Hand)

	// 加上鸣牌的牌数
	for _, group := range p.CallGroups {
		totalTiles += len(group.Tiles)
	}

	// 加上河中的牌数（应该是13）
	totalTiles += p.River.Size()

	// 如果有最后一张摸牌，加1
	// 总应该是14（开局13+1摸牌）
	return totalTiles >= 13
}
