package rom

import (
	"fmt"
	"log"
)

// A Mutable is a memory data that can be changed by the randomizer.
type Mutable interface {
	Mutate([]byte) error // change ROM bytes
	Check([]byte) error  // verify that the mutable matches the ROM
}

// A MutableRange is a length of mutable bytes starting at a given address.
type MutableRange struct {
	Addr     Addr
	Old, New []byte
}

// MutableByte returns a special case of MutableRange with a range of a single
// byte.
func MutableByte(addr Addr, old, new byte) MutableRange {
	return MutableRange{Addr: addr, Old: []byte{old}, New: []byte{new}}
}

// MutableWord returns a special case of MutableRange with a range of a two
// bytes.
func MutableWord(addr Addr, old, new uint16) MutableRange {
	return MutableRange{
		Addr: addr,
		Old:  []byte{byte(old >> 8), byte(old)},
		New:  []byte{byte(new >> 8), byte(new)},
	}
}

// Mutate replaces bytes in its range.
func (mr MutableRange) Mutate(b []byte) error {
	addr := mr.Addr.FullOffset()
	for i, value := range mr.New {
		b[addr+i] = value
	}
	return nil
}

// Check verifies that the range matches the given ROM data.
func (mr MutableRange) Check(b []byte) error {
	addr := mr.Addr.FullOffset()
	for i, value := range mr.Old {
		if b[addr+i] != value {
			return fmt.Errorf("expected %x at %x; found %x",
				mr.Old[i], addr+i, b[addr+i])
		}
	}
	return nil
}

// A MutableSlot is an item slot (chest, gift, etc). It references room data
// and treasure data.
type MutableSlot struct {
	Treasure            *Treasure
	IDAddrs, SubIDAddrs []Addr
	SubIDOffset         byte // routine $16eb requires subID+1
	CollectMode         byte
}

// Mutate replaces the given IDs and subIDs in the given ROM data, and changes
// the associated treasure's collection mode as appropriate.
func (ms MutableSlot) Mutate(b []byte) error {
	for _, addr := range ms.IDAddrs {
		b[addr.FullOffset()] = ms.Treasure.id
	}
	for _, addr := range ms.SubIDAddrs {
		b[addr.FullOffset()] = ms.Treasure.subID + ms.SubIDOffset
	}
	ms.Treasure.mode = ms.CollectMode
	return ms.Treasure.Mutate(b)
}

// Check verifies that the slot's data matches the given ROM data.
func (ms MutableSlot) Check(b []byte) error {
	for _, addr := range ms.IDAddrs {
		if b[addr.FullOffset()] != ms.Treasure.id {
			return fmt.Errorf("expected %x at %x; found %x",
				ms.Treasure.id, addr.FullOffset(), b[addr.FullOffset()])
		}
	}
	for _, addr := range ms.SubIDAddrs {
		if b[addr.FullOffset()] != ms.Treasure.subID {
			return fmt.Errorf("expected %x at %x; found %x",
				ms.Treasure.subID, addr.FullOffset(), b[addr.FullOffset()])
		}
	}
	if ms.CollectMode != ms.Treasure.mode {
		return fmt.Errorf("slot/treasure collect mode mismatch: %x/%x",
			ms.CollectMode, ms.Treasure.mode)
	}

	return nil
}

var ItemSlots = map[string]*MutableSlot{
	"d0 sword chest": &MutableSlot{
		Treasure:    Treasures["sword L-1"],
		IDAddrs:     []Addr{{0x0a, 0x7b86}},
		SubIDAddrs:  []Addr{{0x0a, 0x7b88}},
		SubIDOffset: 1,
		CollectMode: CollectChest,
	},
	"maku key fall": &MutableSlot{
		Treasure:    Treasures["gnarled key"],
		IDAddrs:     []Addr{{0x15, 0x657d}, {0x09, 0x7dff}, {0x09, 0x7de6}},
		SubIDAddrs:  []Addr{{0x15, 0x6580}, {0x09, 0x7e02}},
		CollectMode: CollectFall,
	},
	"boomerang gift": &MutableSlot{
		Treasure:    Treasures["boomerang L-1"],
		IDAddrs:     []Addr{{0x0b, 0x6648}},
		SubIDAddrs:  []Addr{{0x0b, 0x6649}},
		CollectMode: CollectFind2,
	},
	"rod gift": &MutableSlot{
		Treasure:    Treasures["rod"],
		IDAddrs:     []Addr{{0x15, 0x7511}},
		SubIDAddrs:  []Addr{{0x15, 0x750f}},
		SubIDOffset: 1,
		CollectMode: CollectChest, // it's what the data says
	},
	"shovel gift": &MutableSlot{
		Treasure:    Treasures["shovel"],
		IDAddrs:     []Addr{{0x0b, 0x6a6e}},
		SubIDAddrs:  []Addr{{0x0b, 0x6a6f}},
		CollectMode: CollectFind2,
	},
	"d1 satchel": &MutableSlot{
		// addresses are backwards from a normal slot
		Treasure:    Treasures["satchel"],
		IDAddrs:     []Addr{{0x09, 0x669b}},
		SubIDAddrs:  []Addr{{0x09, 0x669a}},
		CollectMode: CollectFind2,
	},
	"d2 bracelet chest": &MutableSlot{
		Treasure:    Treasures["bracelet"],
		IDAddrs:     []Addr{{0x15, 0x5424}},
		SubIDAddrs:  []Addr{{0x15, 0x5425}},
		CollectMode: CollectChest,
	},
	"blaino gift": &MutableSlot{
		Treasure:    Treasures["ricky's gloves"],
		IDAddrs:     []Addr{{0x0b, 0x64ce}},
		SubIDAddrs:  []Addr{{0x0b, 0x64cf}},
		CollectMode: CollectFind1,
	},
	"floodgate key gift": &MutableSlot{
		Treasure:    Treasures["floodgate key"],
		IDAddrs:     []Addr{{0x09, 0x626b}},
		SubIDAddrs:  []Addr{{0x09, 0x626a}},
		CollectMode: CollectFind1,
	},
	"square jewel chest": &MutableSlot{
		Treasure:    Treasures["square jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7397}},
		SubIDAddrs:  []Addr{{0x0b, 0x739b}},
		CollectMode: CollectChest,
	},
	"x-shaped jewel chest": &MutableSlot{
		Treasure:    Treasures["x-shaped jewel"],
		IDAddrs:     []Addr{{0x15, 0x53cd}},
		SubIDAddrs:  []Addr{{0x15, 0x53ce}},
		CollectMode: CollectChest,
	},
	"star ore spot": &MutableSlot{
		Treasure:    Treasures["star ore"],
		IDAddrs:     []Addr{{0x08, 0x62f4}, {0x08, 0x62fe}},
		SubIDAddrs:  []Addr{}, // special case, not set at all
		CollectMode: CollectDig,
	},
	"d3 feather chest": &MutableSlot{
		Treasure:    Treasures["feather L-1"],
		IDAddrs:     []Addr{{0x15, 0x5458}},
		SubIDAddrs:  []Addr{{0x15, 0x5459}},
		CollectMode: CollectChest,
	},
	"master's plaque chest": &MutableSlot{
		Treasure:    Treasures["master's plaque"],
		IDAddrs:     []Addr{{0x15, 0x554d}},
		SubIDAddrs:  []Addr{{0x15, 0x554e}},
		CollectMode: CollectChest,
	},
	"flippers gift": &MutableSlot{
		Treasure:    Treasures["flippers"],
		IDAddrs:     []Addr{{0x0b, 0x7310}, {0x0b, 0x72f3}},
		SubIDAddrs:  []Addr{{0x0b, 0x7311}},
		CollectMode: CollectFind2,
	},
	"spring banana tree": &MutableSlot{
		Treasure:    Treasures["spring banana"],
		IDAddrs:     []Addr{{0x09, 0x66b0}},
		SubIDAddrs:  []Addr{{0x09, 0x66af}},
		CollectMode: CollectFind2,
	},
	"dragon key spot": &MutableSlot{
		Treasure:    Treasures["dragon key"],
		IDAddrs:     []Addr{{0x09, 0x628d}},
		SubIDAddrs:  []Addr{{0x09, 0x628c}},
		CollectMode: CollectFind1,
	},
	"pyramid jewel spot": &MutableSlot{
		Treasure:    Treasures["pyramid jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7350}},
		SubIDAddrs:  []Addr{{0x0b, 0x7351}},
		CollectMode: CollectUnderwater,
	},
	// don't use this slot; no one knows about it and it's not required for
	// anything in a normal playthrough
	/*
		"ring box L-2 gift": &MutableSlot{
			Treasure:    Treasures["ring box L-2"],
			IDAddrs:     []Addr{{0x0b, 0x5c1a}},
			SubIDAddrs:  []Addr{{0x0b, 0x5c1b}},
			CollectMode: CollectGoronGift,
		},
	*/
	"d4 slingshot chest": &MutableSlot{
		Treasure:    Treasures["slingshot L-1"],
		IDAddrs:     []Addr{{0x15, 0x5470}},
		SubIDAddrs:  []Addr{{0x15, 0x5471}},
		CollectMode: CollectChest,
	},
	"d5 magnet gloves chest": &MutableSlot{
		Treasure:    Treasures["magnet gloves"],
		IDAddrs:     []Addr{{0x15, 0x5480}},
		SubIDAddrs:  []Addr{{0x15, 0x5481}},
		CollectMode: CollectChest,
	},
	"round jewel gift": &MutableSlot{
		Treasure:    Treasures["round jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7334}},
		SubIDAddrs:  []Addr{{0x0b, 0x7335}},
		CollectMode: CollectFind2,
	},
	"noble sword spot": &MutableSlot{
		// two cases depending on which sword you enter with
		Treasure:    Treasures["sword L-2"],
		IDAddrs:     []Addr{{0x0b, 0x6417}, {0x0b, 0x641e}},
		SubIDAddrs:  []Addr{{0x0b, 0x6418}, {0x0b, 0x641f}},
		CollectMode: CollectFind1,
	},
	"d6 boomerang chest": &MutableSlot{
		Treasure:    Treasures["boomerang L-2"],
		IDAddrs:     []Addr{{0x15, 0x54c0}},
		SubIDAddrs:  []Addr{{0x15, 0x54c1}},
		CollectMode: CollectChest,
	},
	"rusty bell spot": &MutableSlot{
		Treasure:    Treasures["rusty bell"],
		IDAddrs:     []Addr{{0x09, 0x6476}},
		SubIDAddrs:  []Addr{{0x09, 0x6475}},
		CollectMode: CollectFind2,
	},
	"d7 cape chest": &MutableSlot{
		Treasure:    Treasures["feather L-2"],
		IDAddrs:     []Addr{{0x15, 0x54e1}},
		SubIDAddrs:  []Addr{{0x15, 0x54e2}},
		CollectMode: CollectChest,
	},
	"d8 HSS chest": &MutableSlot{
		Treasure:    Treasures["slingshot L-2"],
		IDAddrs:     []Addr{{0x15, 0x551d}},
		SubIDAddrs:  []Addr{{0x15, 0x551e}},
		CollectMode: CollectChest,
	},

	// these are "fake" item slots in that they don't slot real treasures
	"ember tree": &MutableSlot{
		Treasure: Treasures["ember tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x64cb}},
	},
	"mystery tree": &MutableSlot{
		Treasure: Treasures["mystery tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x67dd}},
	},
	"scent tree": &MutableSlot{
		Treasure: Treasures["scent tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x685c}},
	},
	"pegasus tree": &MutableSlot{
		Treasure: Treasures["pegasus tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x6870}},
	},
	"sunken gale tree": &MutableSlot{
		Treasure: Treasures["gale tree seeds 1"],
		IDAddrs:  []Addr{{0x11, 0x69b0}},
	},
	"tarm gale tree": &MutableSlot{
		Treasure: Treasures["gale tree seeds 2"],
		IDAddrs:  []Addr{{0x11, 0x6a46}},
	},
}

var codeMutables = map[string]Mutable{
	// have maku gate open from start
	"maku gate check": MutableByte(Addr{0x04, 0x61a3}, 0x7e, 0x66),

	// have horon village shop stock *and* sell items from the start, including
	// the flute. also don't disable the flute appearing until actually getting
	// ricky's flute; normally it disappears as soon as you enter the screen
	// northeast of d1 (or ricky's spot, whichever comes first).
	"horon shop stock check":   MutableByte(Addr{0x08, 0x4adb}, 0x05, 0x02),
	"horon shop sell check":    MutableByte(Addr{0x08, 0x48d0}, 0x05, 0x02),
	"horon shop flute check 1": MutableByte(Addr{0x08, 0x4b02}, 0xcb, 0xf6),
	"horon shop flute check 2": MutableByte(Addr{0x08, 0x4afc}, 0x6f, 0x7f),

	// subrosian dancing's flute prize is normally disabled by visiting the
	// same areas as the horon shop's flute.
	"dance hall flute check": MutableByte(Addr{0x09, 0x5e21}, 0x20, 0x80),

	// initiate all these events without requiring essences
	"ricky spawn check":         MutableByte(Addr{0x09, 0x4e68}, 0xcb, 0xf6),
	"rosa spawn check":          MutableByte(Addr{0x09, 0x678c}, 0x40, 0x02),
	"dimitri essence check":     MutableByte(Addr{0x09, 0x4e36}, 0xcb, 0xf6),
	"dimitri flipper check":     MutableByte(Addr{0x09, 0x4e4c}, 0x2e, 0x04),
	"master essence check 2":    MutableByte(Addr{0x0a, 0x4bea}, 0x40, 0x02),
	"master essence check 1":    MutableByte(Addr{0x0a, 0x4bf5}, 0x02, 0x00),
	"round jewel essence check": MutableByte(Addr{0x0a, 0x4f8b}, 0x05, 0x00),
	"pirate essence check":      MutableByte(Addr{0x08, 0x6c32}, 0x20, 0x00),
	"eruption check 1":          MutableByte(Addr{0x08, 0x7c41}, 0x07, 0x00),
	"eruption check 2":          MutableByte(Addr{0x08, 0x7cd3}, 0x07, 0x00),

	// count number of essences, not highest number essence
	"maku seed check 1": MutableByte(Addr{0x09, 0x7d8d}, 0xea, 0x76),
	"maku seed check 2": MutableByte(Addr{0x09, 0x7d8f}, 0x30, 0x18),

	// feather game: don't give fools ore, and don't return fools ore
	"get fools ore 1": MutableByte(Addr{0x14, 0x4111}, 0xe0, 0xf0),
	"get fools ore 2": MutableByte(Addr{0x14, 0x4112}, 0x2e, 0xf0),
	"get fools ore 3": MutableByte(Addr{0x14, 0x4113}, 0x5d, 0xf0),
	"lose fools ore":  MutableByte(Addr{0x3f, 0x454b}, 0x1e, 0x00),

	// stop the hero's cave event from giving you a second wooden sword that
	// you use to spin slash
	"wooden sword second item": MutableByte(Addr{0x0a, 0x7baf}, 0x05, 0x00),

	// change the noble sword's animation pointers to match regular items
	"noble sword anim 1": MutableWord(Addr{0x14, 0x4c67}, 0xe951, 0xa94f),
	"noble sword anim 2": MutableWord(Addr{0x14, 0x4e37}, 0x8364, 0xdf60),

	// getting the L-2 (or L-3) sword in the lost woods gives you two items;
	// one for the item itself and another that gives you the item and also
	// makes you do a spin slash animation. zero the second ID bytes so that
	// one slot doesn't give two items / the same item twice.
	"noble sword second item":  MutableByte(Addr{0x0b, 0x641a}, 0x05, 0x00),
	"master sword second item": MutableByte(Addr{0x0b, 0x6421}, 0x05, 0x00),
}

// like the item slots, these are unchanged by default until the randomizer
// touches them.
var dataMutables = map[string]Mutable{
	// these scenes use specific item sprites not tied to treasure data
	"wooden sword graphics": MutableRange{
		Addr: Addr{0x3f, 0x65f4},
		Old:  []byte{0x60, 0x00, 0x00},
		New:  []byte{0x60, 0x00, 0x00},
	},
	"rod graphics": MutableRange{
		Addr: Addr{0x3f, 0x6ba3},
		Old:  []byte{0x60, 0x10, 0x21},
		New:  []byte{0x60, 0x10, 0x21},
	},
	"noble sword graphics": MutableRange{
		Addr: Addr{0x3f, 0x6975},
		Old:  []byte{0x4e, 0x1a, 0x50},
		New:  []byte{0x4e, 0x1a, 0x50},
	},
	"master sword graphics": MutableRange{
		Addr: Addr{0x3f, 0x6978},
		Old:  []byte{0x4e, 0x1a, 0x40},
		New:  []byte{0x4e, 0x1a, 0x40},
	},
}

// get a collated map of all mutables
func getAllMutables() map[string]Mutable {
	slotMutables := make(map[string]Mutable)
	for k, v := range ItemSlots {
		slotMutables[k] = v
	}
	treasureMutables := make(map[string]Mutable)
	for k, v := range Treasures {
		treasureMutables[k] = v
	}

	mutableSets := []map[string]Mutable{
		codeMutables,
		treasureMutables,
		slotMutables,
		dataMutables,
	}

	// initialize master map w/ adequate capacity
	count := 0
	for _, set := range mutableSets {
		count += len(set)
	}
	allMutables := make(map[string]Mutable, count)

	// add mutables to master map
	for _, set := range mutableSets {
		for k, v := range set {
			if _, ok := allMutables[k]; ok {
				log.Fatalf("duplicate mutable key: %s", k)
			}
			allMutables[k] = v
		}
	}

	return allMutables
}
