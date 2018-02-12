package main

import (
	"time"
	"github.com/jinzhu/gorm"
	// "github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type BarangDB struct {
	// gorm.Model
	Sku     string	`gorm:"primary_key" form:"sku" json:"sku"`
	Nama	string	`sql:"not null" form:"nama" json:"nama"`
}

type PembelianDB struct {
	// Id 			int
	Id     		int	`gorm:"primary_key,AUTO_INCREMENT" form:"id" json:"id"`
	Waktu		time.Time	`sql:"type:datetime" form:"waktu" json:"waktu"`
	Waktuupdt	time.Time	`sql:"type:datetime" form:"waktuupdt" json:"waktuupdt"`
	Sku			string	`sql:"type:varchar(255) REFERENCES barang_dbs(sku) not null" form:"sku" json:"sku"`
	Jpemesanan	int	`sql:"not null" form:"jpemesanan" json:"jpemesanan"`
	Jditerima	int	`sql:"not null" form:"jditerima" json:"jditerima"`
	Hbeli		int	`sql:"not null" form:"hbeli" json:"hbeli"`
	Kwitansi	string	`sql:"not null" form:"kwitansi" json:"kwitansi"`
	Catatan		string	`sql:"not null" form:"catatan" json:"catatan"`
}

type PenjualanDB struct {
	Id     		int	`gorm:"primary_key,AUTO_INCREMENT" form:"id" json:"id"`
	Waktu		time.Time	`sql:"type:datetime" form:"waktu" json:"waktu"`
	Waktuupdt	time.Time	`sql:"type:datetime" form:"waktuupdt" json:"waktuupdt"`
	Sku			string	`sql:"type:varchar(255) REFERENCES barang_dbs(sku) not null" form:"sku" json:"sku"`
	Jkeluar		int	`sql:"not null" form:"jkeluar" json:"jkeluar"`
	Hjual		int	`sql:"not null" form:"hjual" json:"hjual"`
	Catatan		string	`sql:"type:varchar(255)" form:"catatan" json:"catatan"`
}


type Listbarang struct {
	BarangDB
	Jumlah	int `json:"jumlah"`
}

type Laporannilaibarang struct {
	Listbarang
	Ratarata 	int `json:"ratarata"`
	Total		int `json:"total"`
}

type Pembelians struct {
	PembelianDB
	Total 		int `json:"total"`
}

type Penjualans struct {
	PenjualanDB
	Total 		int `json:"total"`
}

type Laporanpenjualan struct {
	Penjualans
	Nama 	string 	`json:"nama"`
	Hbeli 	int 	`json:"hbeli"`
	Laba	int 	`json:"laba"`
}



func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", "data.db")
	db.Exec("PRAGMA foreign_keys = ON")
	// Display SQL queries
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table
	if !db.HasTable(&BarangDB{}) {
		db.CreateTable(&BarangDB{})
		// db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&BarangDB{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&BarangDB{})
	}

	if !db.HasTable(&PembelianDB{}) {
		db.CreateTable(&PembelianDB{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&PembelianDB{})
	}

	if !db.HasTable(&PenjualanDB{}) {
		db.CreateTable(&PenjualanDB{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&PenjualanDB{})
	}

	return db
}