package bot

import (
	"encoding/json"
	"time"
)

type MembersResponse []Member

type Member struct {
	AverageItemPower   json.Number        `json:"AverageItemPower"`
	Equipment          Equipment          `json:"Equipment"`
	Inventory          []interface{}      `json:"Inventory"`
	Name               string             `json:"Name"`
	ID                 string             `json:"Id"`
	GuildName          GuildName          `json:"GuildName"`
	GuildID            GuildID            `json:"GuildId"`
	AllianceName       *AllianceName      `json:"AllianceName"`
	AllianceID         AllianceIDEnum     `json:"AllianceId"`
	AllianceTag        AllianceTagEnum    `json:"AllianceTag"`
	Avatar             string             `json:"Avatar"`
	AvatarRing         AvatarRingEnum     `json:"AvatarRing"`
	DeathFame          int64              `json:"DeathFame"`
	KillFame           int64              `json:"KillFame"`
	FameRatio          float64            `json:"FameRatio"`
	LifetimeStatistics LifetimeStatistics `json:"LifetimeStatistics"`
}

type Equipment struct {
	MainHand interface{} `json:"MainHand"`
	OffHand  interface{} `json:"OffHand"`
	Head     interface{} `json:"Head"`
	Armor    interface{} `json:"Armor"`
	Shoes    interface{} `json:"Shoes"`
	Bag      interface{} `json:"Bag"`
	Cape     interface{} `json:"Cape"`
	Mount    interface{} `json:"Mount"`
	Potion   interface{} `json:"Potion"`
	Food     interface{} `json:"Food"`
}

type LifetimeStatistics struct {
	PVE           PVE       `json:"PvE"`
	Gathering     Gathering `json:"Gathering"`
	Crafting      Crafting  `json:"Crafting"`
	CrystalLeague int64     `json:"CrystalLeague"`
	FishingFame   int64     `json:"FishingFame"`
	FarmingFame   int64     `json:"FarmingFame"`
	Timestamp     time.Time `json:"Timestamp"`
}

type Crafting struct {
	Total    int64 `json:"Total"`
	Royal    int64 `json:"Royal"`
	Outlands int64 `json:"Outlands"`
	Avalon   int64 `json:"Avalon"`
}

type Gathering struct {
	Fiber Crafting `json:"Fiber"`
	Hide  Crafting `json:"Hide"`
	Ore   Crafting `json:"Ore"`
	Rock  Crafting `json:"Rock"`
	Wood  Crafting `json:"Wood"`
	All   Crafting `json:"All"`
}

type PVE struct {
	Total            int64 `json:"Total"`
	Royal            int64 `json:"Royal"`
	Outlands         int64 `json:"Outlands"`
	Avalon           int64 `json:"Avalon"`
	Hellgate         int64 `json:"Hellgate"`
	CorruptedDungeon int64 `json:"CorruptedDungeon"`
	Mists            int64 `json:"Mists"`
}

type AllianceIDEnum string

const (
	AllianceID             AllianceIDEnum = ""
	QRxb4JJ1RVygfbCEFUkDkQ AllianceIDEnum = "qRxb4jJ1RVygfbCEFUkDkQ"
)

type AllianceName string

const (
	AllianceNameKYSS AllianceName = "KYSS"
	KeepYourselfSfe  AllianceName = "Keep Yourself Sfe"
)

type AllianceTagEnum string

const (
	AllianceTag     AllianceTagEnum = ""
	AllianceTagKYSS AllianceTagEnum = "KYSS"
)

type AvatarRingEnum string

const (
	AvatarRing AvatarRingEnum = ""
	Ring1      AvatarRingEnum = "RING1"
)

type GuildID string

const (
	IMCZVRcgTm6IUhtph9Ddg GuildID = "ImCZVRcgTm6_iUhtph9ddg"
)

type GuildName string

const (
	SolExanime GuildName = "Sol Exanime"
)
