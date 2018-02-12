# toko
building toko app using go

DOKUMENTASI API:

IMPORT DATA DARI CSV
METHOD 		: POST
ENDPOINT 	: "api/v1/importCSV"
BODY 		: file.csv, "tipe(penjualan, pembelian, barang)"

CREATE barang
METHOD 		: POST
ENDPOINT 	: "api/v1/barangs"
BODY 		: {
			  "sku" : "SSI-D01220307-L-SAL",
			  "nama" : "Devibav Plain Trump Blouse (L,Salem)"
			}

MENAMPILKAN LIST barang
METHOD 		: GET
ENDPOINT 	: "api/v1/barangs"

MENAMPILKAN DETAIL barang
METHOD 		: GET
ENDPOINT 	: "api/v1/barangs"
PARAMS 		: sku

MELAKUKAN UPDATE barang
METHOD 		: PUT
ENDPOINT 	: "api/v1/barangs"
PARAMS 		: sku

MENGHAPUS DATA barang
METHOD 		: DELETE
ENDPOINT 	: "api/v1/barangs"
PARAMS 		: sku

CREATE pembelian
METHOD 		: POST
ENDPOINT 	: "api/v1/pembelians"
BODY 		: {
			  	"sku" : "SSI-D01220307-L-SAL",
				"jpemesanan" : 40,
				"jditerima" : 40,
				"hbeli" : 80000,
				"kwitansi" : "20171102-29823",
				"catatan" : "2017/11/06 terima 40"
			}

MENAMPILKAN LIST pembelian
METHOD 		: GET
ENDPOINT 	: "api/v1/pembelians"

MENAMPILKAN DETAIL pembelian
METHOD 		: GET
ENDPOINT 	: "api/v1/pembelians"
PARAMS 		: id

MELAKUKAN UPDATE pembelian
METHOD 		: PUT
ENDPOINT 	: "api/v1/pembelians"
PARAMS 		: id

MENGHAPUS DATA pembelian
METHOD 		: DELETE
ENDPOINT 	: "api/v1/pembelians"
PARAMS 		: id

CREATE penjualan
METHOD 		: POST
ENDPOINT 	: "api/v1/penjualans"
BODY 		: {
			  	"sku" : "SSI-D01220307-XL-SAL",
				"jkeluar" :2,
				"hjual" : 115000,
				"catatan" : "Pesanan ID-20180106-052436"
			}

MENAMPILKAN LIST penjualan
METHOD 		: GET
ENDPOINT 	: "api/v1/penjualans"

MENAMPILKAN DETAIL penjualan
METHOD 		: GET
ENDPOINT 	: "api/v1/penjualans"
PARAMS 		: id

MELAKUKAN UPDATE penjualan
METHOD 		: PUT
ENDPOINT 	: "api/v1/penjualans"
PARAMS 		: id

MENGHAPUS DATA penjualan
METHOD 		: DELETE
ENDPOINT 	: "api/v1/penjualans"
PARAMS 		: id

MENAMPILKAN LAPORAN nilaibarang
METHOD 		: GET
ENDPOINT 	: "api/v1/laporan/nilaibarang"

EXPORT LAPORAN nilaibarang CSV
METHOD 		: GET
ENDPOINT 	: "api/v1/laporan/nilaibarang/csv"

MENAMPILKAN LAPORAN penjualan
METHOD 		: GET
ENDPOINT 	: "api/v1/laporan/penjualan"
PARAMS 		: "tanggal awal", "tanggal akhir"

EXPORT LAPORAN penjualan CSV
METHOD 		: GET
ENDPOINT 	: "api/v1/laporan/penjualan"
PARAMS 		: "tanggal awal", "tanggal akhir", "csv"
