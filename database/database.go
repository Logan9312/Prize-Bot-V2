package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DatabaseConnect(password, host, env string) {
	fmt.Println("Connecting to Database...")
	defer fmt.Println("Bot has finished attempting to connect to the database!")

	if env == "prod" {
		DB = ProdDB(password, host)
	} else if env == "local" {
		DB = LocalDB()
	}

	err := DB.AutoMigrate(Event{}, AuctionSetup{}, Auction{}, AuctionQueue{}, GiveawaySetup{}, Giveaway{}, ClaimSetup{}, CurrencySetup{}, Claim{}, DevSetup{}, UserProfile{}, ShopSetup{}, WhiteLabels{})
	if err != nil {
		fmt.Println(err)
	}

}

func LocalDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("/tmp/test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("Error creating Sqlite Database", err)
	}

	return db
}

func ProdDB(password, host string) *gorm.DB {
	dbuser := "auctionbot"
	port := "3306"
	dbname := "auction"

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, dbuser, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func GetAuctionSettings(guildID string) (AuctionSetup, error) {

	auctionSettings := AuctionSetup{}

	result := DB.First(&auctionSettings, guildID)

	//TODO Test if this still works before fetching data
	if auctionSettings.ChannelPrefix == nil {
		auctionSettings.ChannelPrefix = Ptr("💸│")
	}
	if result.Error != nil {
		return auctionSettings, fmt.Errorf("Error getting auction settings: %w", result.Error)
	}

	return auctionSettings, nil
}

func GetAuctionData(channelID string) (Auction, error) {
	auction := Auction{}

	result := DB.Preload("Event").Where("event.channel_id = ?", channelID).First(&auction)
	if result.Error != nil {
		return auction, fmt.Errorf("Error getting auction settings: %w", result.Error)
	}

	return auction, nil
}

func Ptr[T any](v T) *T {
	return &v
}
