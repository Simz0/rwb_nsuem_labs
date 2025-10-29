package models

import (
	"time"
)

type Order struct {
	OrderUid          string    `json:"order_uid" example:"b563feb7b2b84b6test" gorm:"primaryKey"`
	TrackNumber       string    `json:"track_number" example:"WBILMTESTTRACK"`
	Entry             string    `json:"entry" example:"WBIL"`
	DeliveryID        uint      `json:"-"`
	Delivery          Delivery  `json:"delivery" gorm:"foreignKey:DeliveryID"`
	PaymentID         uint      `json:"-"`
	Payment           Payment   `json:"payment" gorm:"foreignKey:PaymentID"`
	Items             []Item    `json:"items" gorm:"foreignKey:OrderUID"`
	Locale            string    `json:"locale" example:"en"`
	InternalSignature string    `json:"internal_signature" example:""`
	CustomerId        string    `json:"customer_id" example:"test"`
	DeliveryService   string    `json:"delivery_service" example:"meest"`
	ShardKey          string    `json:"shardkey" example:"9"`
	SMId              uint      `json:"sm_id" example:"99"`
	DateCreated       time.Time `json:"date_created" example:"2021-11-26T06:22:19Z"`
	OofShard          string    `json:"oof_shard" example:"1"`
}

type Delivery struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name" example:"Test Testov"`
	Phone   string `json:"phone" example:"+9720000000"`
	Zip     string `json:"zip" example:"2639809"`
	City    string `json:"city" example:"Kiryat Mozkin"`
	Address string `json:"address" example:"Ploshad Mira 15"`
	Region  string `json:"region" example:"Kraiot"`
	Email   string `json:"email" example:"test@gmail.com"`
}

type Payment struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Transaction  string `json:"transaction" example:"b563feb7b2b84b6test"`
	RequestId    string `json:"request_id" example:""`
	Currency     string `json:"currency" example:"USD"`
	Provider     string `json:"provider" example:"wbpay"`
	Amount       int32  `json:"amount" example:"1817"`
	PaymentDT    int64  `json:"payment_dt" example:"1637907727"`
	Bank         string `json:"bank" example:"alpha"`
	DeliveryCost int32  `json:"delivery_cost" example:"1500"`
	GoodsTotal   int32  `json:"goods_total" example:"317"`
	CustomFee    int32  `json:"custom_fee" example:"0"`
}

type Item struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	OrderUID    string `json:"-" gorm:"size:255"`
	ChrtId      uint32 `json:"chrt_id" example:"9934930"`
	TrackNumber string `json:"track_number" example:"WBILMTESTTRACK"`
	Price       int32  `json:"price" example:"453"`
	Rid         string `json:"rid" example:"ab4219087a764ae0btest"`
	Name        string `json:"name" example:"Mascaras"`
	Sale        int    `json:"sale" example:"30"`
	Size        string `json:"size" example:"0"`
	TotalPrice  int32  `json:"total_price" example:"317"`
	NMId        uint32 `json:"nm_id" example:"2389212"`
	Brand       string `json:"brand" example:"Vivienne Sabo"`
	Status      int    `json:"status" example:"202"`
}

// Альтернатива JSONB для асболютной статики, но в теории фильтрация сложнее будет (хоть она и не нужна)
// Для масштабирования более приятен первый вариант, поэтому я оставил именно его

// import (
// 	"time"
// 	"encoding/json"
// )

// type Order struct {
// 	OrderUid          string    `json:"order_uid" example:"b563feb7b2b84b6test" gorm:"primaryKey"`
// 	TrackNumber       string    `json:"track_number" example:"WBILMTESTTRACK"`
// 	Entry             string    `json:"entry" example:"WBIL"`
// 	Delivery          JSONB     `json:"delivery" gorm:"type:jsonb"`
// 	Payment           JSONB     `json:"payment" gorm:"type:jsonb"`
// 	Items             JSONB     `json:"items" gorm:"type:jsonb"`
// 	Locale            string    `json:"locale" example:"en"`
// 	InternalSignature string    `json:"internal_signature" example:""`
// 	CustomerId        string    `json:"customer_id" example:"test"`
// 	DeliveryService   string    `json:"delivery_service" example:"meest"`
// 	ShardKey          string    `json:"shardkey" example:"9"`
// 	SMId              uint      `json:"sm_id" example:"99"`
// 	DateCreated       time.Time `json:"date_created" example:"2021-11-26T06:22:19Z"`
// 	OofShard          string    `json:"oof_shard" example:"1"`
// }

// type JSONB map[string]interface{}

// func (j JSONB) Value() (driver.Value, error) {
// 	return json.Marshal(j)
// }

// func (j *JSONB) Scan(value interface{}) error {
// 	if value == nil {
// 		*j = nil
// 		return nil
// 	}
// 	data, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}
// 	return json.Unmarshal(data, &j)
// }
