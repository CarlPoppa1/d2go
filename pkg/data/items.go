package data

import (
	"strings"

	"github.com/hectorgimenez/d2go/pkg/data/item"
	"github.com/hectorgimenez/d2go/pkg/data/stat"
)

type Inventory struct {
	Belt        Belt
	AllItems    []Item
	Gold        int
	StashedGold [4]int
}

func (i Inventory) Find(name item.Name, locations ...item.LocationType) (Item, bool) {
	for _, it := range i.AllItems {
		if strings.EqualFold(string(it.Name), string(name)) {
			// If no locations are specified, return the first item found
			if len(locations) == 0 {
				return it, true
			}

			for _, l := range locations {
				if it.Location.LocationType == l {
					return it, true
				}
			}
		}
	}

	return Item{}, false
}

func (i Inventory) FindByID(unitID UnitID) (Item, bool) {
	for _, it := range i.AllItems {
		if it.UnitID == unitID {
			return it, true
		}
	}

	return Item{}, false
}

func (i Inventory) ByLocation(locations ...item.LocationType) []Item {
	var items []Item

	for _, it := range i.AllItems {
		for _, l := range locations {
			if it.Location.LocationType == l {
				items = append(items, it)
			}
		}
	}

	return items
}

func (i Inventory) Matrix() [4][10]bool {
	invMatrix := [4][10]bool{} // false = empty, true = occupied
	for _, itm := range i.ByLocation(item.LocationInventory) {
		for k := range itm.Desc().InventoryWidth {
			for j := range itm.Desc().InventoryHeight {
				invMatrix[itm.Position.Y+j][itm.Position.X+k] = true
			}
		}
	}

	return invMatrix
}

type UnitID int

type Item struct {
	ID int
	UnitID
	Name                 item.Name
	Quality              item.Quality
	Position             Position
	Location             item.Location
	Ethereal             bool
	IsHovered            bool
	BaseStats            stat.Stats
	Stats                stat.Stats
	Identified           bool
	IsRuneword           bool
	IsNamed              bool
	IsStartItem          bool
	IsEar                bool
	IsBroken             bool
	IsEquipped           bool
	HasSockets           bool
	InTradeOrStoreScreen bool
	IsInSocket           bool
}

type Drop struct {
	Item         Item
	Rule         string
	RuleFile     string
	DropLocation string
}

func (i Item) Desc() item.Description {
	return item.Desc[i.ID]
}

func (i Item) Type() item.Type {
	return i.Desc().GetType()
}

func (i Item) IsPotion() bool {
	return i.IsHealingPotion() || i.IsManaPotion() || i.IsRejuvPotion()
}

func (i Item) IsHealingPotion() bool {
	return i.Type().IsType(item.TypeHealingPotion)
}

func (i Item) IsManaPotion() bool {
	return i.Type().IsType(item.TypeManaPotion)
}

func (i Item) IsRejuvPotion() bool {
	return i.Type().IsType(item.TypeRejuvPotion)
}

func (i Item) IsFromQuest() bool {
	return i.Type().IsType(item.TypeQuest)
}

func (i Item) FindStat(id stat.ID, layer int) (stat.Data, bool) {
	st, found := i.Stats.FindStat(id, layer)
	if found {
		return st, true
	}

	return i.BaseStats.FindStat(id, layer)
}
