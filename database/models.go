package database

import (
	"time"
)

type DevSetup struct {
	BotID   string `gorm:"primaryKey"`
	Version string
}

type WhiteLabels struct {
	BotID    string `gorm:"primaryKey;autoIncrement:false"`
	UserID   string `gorm:"primaryKey;autoIncrement:false"`
	BotToken string
}

type AuctionSetup struct {
	GuildID         string `gorm:"primaryKey"`
	Category        *string
	AlertRole       *string
	CurrencyID      uint
	Currency        *Currency
	LogChannel      *string
	HostRole        *string
	SnipeExtension  *time.Duration
	SnipeRange      *time.Duration
	IntegerOnly     *bool
	ChannelOverride *string
	ChannelLock     *bool
	ChannelPrefix   *string
}

// TODO Auction Setup should NEVER be used aside from initial auction creations
type Auction struct {
	ID             int `gorm:"primaryKey"`
	Event          Event
	EventID        uint
	Currency       *Currency
	CurrencyID     uint
	Bid            float64
	WinnerID       *string
	IncrementMin   *float64
	IncrementMax   *float64
	TargetPrice    *float64
	Buyout         *float64
	IntegerOnly    bool
	BidHistory     *string
	SnipeExtension *time.Duration
	SnipeRange     *time.Duration
}

// TODO Potentially add EventSettings to make it easier to create functions that work on multiple event types
type AuctionQueue struct {
	ID              int `gorm:"primaryKey"`
	ChannelID       string
	Bid             float64
	StartTime       time.Time
	EndTime         time.Time
	GuildID         string
	Item            string
	Host            string
	Currency        string
	IncrementMin    float64
	IncrementMax    float64
	Description     string
	ImageURL        string
	Category        string
	TargetPrice     float64
	Buyout          float64
	CurrencySide    string
	IntegerOnly     bool
	SnipeExtension  time.Duration
	SnipeRange      time.Duration
	AlertRole       string
	Note            string
	ChannelOverride string
	ChannelLock     bool
	UseCurrency     bool
	ChannelPrefix   string
}

// ClaimSetup FromMake sure to remove LogChannel and ClaimMessage from auction log
type ClaimSetup struct {
	GuildID         string `gorm:"primaryKey"`
	Category        string
	StaffRole       string
	Instructions    string
	LogChannel      string
	Expiration      string
	DisableClaiming bool
	ChannelPrefix   string
}

type Claim struct {
	MessageID   string `gorm:"primaryKey"`
	ChannelID   string
	GuildID     string
	Item        string
	Type        string
	Winner      string
	Cost        float64
	Host        string
	BidHistory  string
	Note        string
	ImageURL    string
	TicketID    string
	Description string
}

type Giveaway struct {
	MessageID   string `gorm:"primaryKey"`
	ChannelID   string
	GuildID     string
	Item        string
	EndTime     time.Time
	Description string
	Host        string
	Winners     float64
	ImageURL    string
	Finished    bool
}

type GiveawaySetup struct {
	GuildID    string `gorm:"primaryKey"`
	HostRole   string
	AlertRole  string
	LogChannel string
}

type ShopSetup struct {
	GuildID    string `gorm:"primaryKey"`
	HostRole   string
	AlertRole  string
	LogChannel string
}

type Currency struct {
	ID          uint `gorm:"primaryKey"`
	Symbol      string
	RightSide   bool // True if the currency should display on the right side
	UseCurrency bool // True if user's currency is used for transaction.
}

type CurrencySetup struct {
	GuildID  string `gorm:"primaryKey"`
	Currency string
	Side     string
}

type UserProfile struct {
	UserID  string `gorm:"primaryKey;autoIncrement:false"`
	GuildID string `gorm:"primaryKey;autoIncrement:false"`
	Balance float64
}

type Quest struct {
	MessageID string `gorm:"primaryKey;autoIncrement:false"`
}

type Errors struct {
	ErrorID string `gorm:"primaryKey"`
	UserID  string `gorm:"primaryKey;autoIncrement:false"`
}

type Event struct {
	ID          uint `gorm:"primaryKey"`
	BotID       string
	GuildID     string
	Host        string
	Item        string
	ChannelID   *string
	MessageID   *string
	StartTime   *time.Time
	EndTime     *time.Time
	ImageURL    *string
	Description *string
	Note        *string
	AlertRole   *string
}
