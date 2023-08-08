package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Settings interface {
	*AuctionSetup | *ClaimSetup | *CurrencySetup | *DevSetup | *GiveawaySetup | *ShopSetup
}

func DatabaseConnect(password, host, env string) {
	fmt.Println("Connecting to Database...")
	defer fmt.Println("Bot has finished attempting to connect to the database!")

	if env == "prod" {
		DB = ProdDB(password, host).Session(&gorm.Session{FullSaveAssociations: true})
	} else if env == "local" {
		DB = LocalDB().Session(&gorm.Session{FullSaveAssociations: true})
	}

	err := DB.AutoMigrate(Event{}, AuctionSetup{}, Auction{}, AuctionQueue{}, GiveawaySetup{}, Giveaway{}, ClaimSetup{}, CurrencySetup{}, Claim{}, DevSetup{}, UserProfile{}, ShopSetup{}, WhiteLabels{})
	if err != nil {
		fmt.Println(err)
	}

}

func LocalDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
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
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func GetAuctionSettings(guildID string) (AuctionSetup, error) {

	auctionSettings := AuctionSetup{
		GuildID: guildID,
	}

	result := DB.FirstOrInit(&auctionSettings)
	if result.Error != nil {
		return auctionSettings, fmt.Errorf("Error getting auction settings: %w", result.Error)
	}

	//TODO Test if this still works before fetching data
	if auctionSettings.ChannelPrefix == nil {
		auctionSettings.ChannelPrefix = Ptr("💸│")
	}

	return auctionSettings, nil
}

func SaveSettings[T Settings](data map[string]any, model T) error {
	return DB.Model(model).Save(data).Error
}

func GetAuctionData(channelID string) (Auction, error) {
	var auction Auction

	err := DB.Joins("JOIN events ON events.id = auctions.event_id").Where("events.channel_id = ?", channelID).Preload("Event").First(&auction).Error
	return auction, err
}

func SaveAuction(auction *Auction) error {
	return DB.Save(auction).Error
}

func Ptr[T any](v T) *T {
	return &v
}
