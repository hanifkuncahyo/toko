package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

var loc, _ = time.LoadLocation("Asia/Jakarta")

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {

	router := gin.Default()

	router.Use(Cors())

	v1 := router.Group("api/v1")
	{
		v1.GET("/", Index)
		v1.POST("/importCSV", ReadCSV)
		v1.POST("/barangs", PostBarang)
		v1.GET("/barangs", GetBarangs)
		v1.GET("/barangs/:sku", GetBarang)
		v1.PUT("/barangs/:sku", UpdateBarang)
		v1.DELETE("/barangs/:sku", DeleteBarang)
		v1.GET("/laporan/nilaibarang", GetLaporanBarangs)
		v1.GET("/laporan/nilaibarang/:csv", GetLaporanBarangs, ExportCSV)
		v1.GET("/laporan/penjualan/:tmfirst/:tmlast", GetLaporanPenjualans)
		v1.GET("/laporan/penjualan/:tmfirst/:tmlast/:csv", GetLaporanPenjualans, ExportCSV)
		v1.POST("/pembelians", PostPembelian)
		v1.GET("/pembelians", GetPembelians)
		v1.GET("/pembelians/:id", GetPembelian)
		v1.PUT("/pembelians/:id", UpdatePembelian)
		v1.DELETE("/pembelians/:id", DeletePembelian)
		v1.POST("/penjualans", CheckStock, PostPenjualan)
		v1.GET("/penjualans", GetPenjualans)
		v1.GET("/penjualans/:id", GetPenjualan)
		v1.PUT("/penjualans/:id", CheckStock, UpdatePenjualan)
		v1.DELETE("/penjualans/:id", DeletePenjualan)
	}
	
	router.Run(":8080")

}
