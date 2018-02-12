package main

import (
    "time"
    "net/http"
    "fmt"
    // "os"
    "encoding/csv"
    "strings"
    "bytes"
    "log"
    "strconv"
    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"
)

func ReadCSV(c *gin.Context) {
    file, _, err := c.Request.FormFile("upload")
    if err != nil {
        fmt.Println("Error", err)
        c.JSON(422, gin.H{"error": "Error occured"})
    }

    defer file.Close()
    tipe := c.PostForm("tipe")

    reader := csv.NewReader(file)
    record, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error", err)
    }

    db := InitDb()
    defer db.Close()

    for value:= range record{ 
        if tipe == "penjualan"{
            jkeluar,_ := strconv.Atoi(record[value][3])
            hjual,_ := strconv.Atoi(record[value][4])
            waktu := time.Now().In(loc)
            waktuupdt := time.Now().In(loc)
            penjualan := PenjualanDB{
                Waktu: waktu,
                Waktuupdt: waktuupdt,
                Sku: record[value][2],
                Jkeluar: jkeluar,
                Hjual: hjual,
                Catatan: record[value][5],
            }
            
            if penjualan.Sku != "" && penjualan.Jkeluar != 0 && penjualan.Hjual != 0{
                
                err := db.Create(&penjualan)
                if err.RowsAffected == 0 {
                    c.JSON(422, gin.H{"error": "Error occured"})
                }
            }
        } else if tipe == "pembelian"{
            
            waktu := time.Now().In(loc)
            waktuupdt := time.Now().In(loc)
            jpemesanan,_ := strconv.Atoi(record[value][2])
            jditerima,_ := strconv.Atoi(record[value][3])
            hbeli,_ := strconv.Atoi(record[value][4])

            pembelian := PembelianDB{
                Waktu: waktu,
                Waktuupdt: waktuupdt,
                Sku: record[value][1],
                Jpemesanan : jpemesanan,
                Jditerima : jditerima,
                Hbeli : hbeli,
                Kwitansi: record[value][5],
                Catatan: record[value][6],
            }
            
            if pembelian.Sku != "" && pembelian.Jpemesanan != 0 && pembelian.Jditerima != 0 && pembelian.Hbeli != 0 && pembelian.Kwitansi != "" && pembelian.Catatan != "" {
                // pembelian.Waktu = time.Now().In(loc)
                // pembelian.Waktuupdt = time.Now().In(loc)

                err := db.Create(&pembelian)
                if err.RowsAffected == 0 {
                    c.JSON(422, gin.H{"error": "Error occured"})
                } 
            }
        } else if tipe == "barang"{
            barang := BarangDB{
                Sku: record[value][0],
                Nama: record[value][1],
            }
            
            if barang.Sku != "" && barang.Nama != "" {
                
                err := db.Create(&barang)
                if err.RowsAffected == 0 {
                    c.JSON(422, gin.H{"error": "Error occured"})
                } 
            }
        }

    }

    c.JSON(200, gin.H{"success": "Kebaca"})

}

func stringComp(a, b string) string {
    if a != "" {
        return a
    }
    return b
}

func intComp(a, b int) int {
    if a != 0 {
        return a
    }
    return b
}

func Index(c *gin.Context) {
    c.Writer.WriteHeader(http.StatusOK)
    c.Writer.Write([]byte("Welcome!\n"))
}

func PostBarang(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    var barang BarangDB
    c.Bind(&barang)

    if barang.Nama != "" && barang.Sku != "" {
        db.Create(&barang)
        c.JSON(201, gin.H{"success": barang})
    } else {
        c.JSON(422, gin.H{"error": "Fields are empty"})
    }
}

func GetBarangs(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    var barangs []Listbarang
    
    var sumPB = "(SELECT COALESCE(sum(pb.jditerima), 0) as jumlah FROM pembelian_dbs pb where pb.sku = b.sku)"
    var sumPJ = "(SELECT COALESCE(sum(pj.jkeluar), 0) as jumlah FROM penjualan_dbs pj where pj.sku = b.sku)"
    db.Select("b.*, "+sumPB+" - "+sumPJ+" as jumlah").Table("barang_dbs b").Group("b.sku").Scan(&barangs);

    if len(barangs)==0 {
        c.JSON(404, gin.H{"error": "Barang not found"})
    } else {
        c.JSON(200, barangs)
    }
}

func GetLaporanBarangs(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    var barangs []Laporannilaibarang
    
    var sumPB = "(SELECT COALESCE(sum(pb.jditerima), 0) as jumlah FROM pembelian_dbs pb where pb.sku = b.sku)"
    var sumPJ = "(SELECT COALESCE(sum(pj.jkeluar), 0) as jumlah FROM penjualan_dbs pj where pj.sku = b.sku)"
    var sumPBTotal = "(SELECT COALESCE(sum(pb.jditerima * pb.hbeli), 0) as jumlah FROM pembelian_dbs pb where pb.sku = b.sku)"
    db.Select("b.*, "+sumPB+" - "+sumPJ+" as jumlah, "+sumPBTotal+" / "+sumPB+" as ratarata, ("+sumPB+" - "+sumPJ+") * ("+sumPBTotal+" / "+sumPB+") as total").Table("barang_dbs b").Group("b.sku").Scan(&barangs);
    
    var jtotalbarang,tnilai int 

    for i := 0; i < len(barangs); i++ {
        jtotalbarang += barangs[i].Jumlah
        tnilai += barangs[i].Total
    }

    if len(barangs)==0 {
        c.JSON(404, gin.H{"error": "Barang not found"})
    } else {
        csv := c.Params.ByName("csv")
        if csv == "csv"{
            c.Set("data",barangs)
            c.Set("tipe","nilaibarang")
            c.Next()
        } else {
            c.JSON(200, gin.H{"result": gin.H{"summary":gin.H{"jsku":len(barangs),"jtotalbarang":jtotalbarang,"tnilai":tnilai},"detail":barangs}})
        }
    }
}

func GetBarang(c *gin.Context) {
    db := InitDb()
    defer db.Close()
    sku := c.Params.ByName("sku")

    var barang Listbarang

    var sumPB = "(SELECT COALESCE(sum(pb.jditerima), 0) as jumlah FROM pembelian_dbs pb where pb.sku = b.sku)"
    var sumPJ = "(SELECT COALESCE(sum(pj.jkeluar), 0) as jumlah FROM penjualan_dbs pj where pj.sku = b.sku)"
    db.Select("b.*, "+sumPB+" - "+sumPJ+" as jumlah").Table("barang_dbs b").Group("b.sku").Where("b.sku = ?", sku).Order("b.sku ASC").Scan(&barang);
    
    if barang.Sku != "" {
        c.JSON(200, barang)
    } else {
        c.JSON(404, gin.H{"error": "Barang not found"})
    }
}

func UpdateBarang(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    sku := c.Params.ByName("sku")
    var barang BarangDB
    db.First(&barang, sku)

    if barang.Nama != "" {

        if barang.Sku != "" {
            var newBarang BarangDB
            c.Bind(&newBarang)

            result := BarangDB{
                Sku: barang.Sku,
                Nama: stringComp(newBarang.Nama, barang.Nama),
            }

            db.Save(&result)
            c.JSON(200, gin.H{"success": result})
        } else {
            c.JSON(404, gin.H{"error": "Barang not found"})
        }

    } else {
        c.JSON(422, gin.H{"error": "Fields are empty"})
    }

}

func DeleteBarang(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    sku := c.Params.ByName("sku")
    var barang BarangDB
    db.First(&barang, sku)

    if barang.Sku != "" {
        db.Delete(&barang)
        c.JSON(200, gin.H{"success": "Barang #" + sku + " deleted"})
    } else {
        c.JSON(404, gin.H{"error": "Barang not found"})
    }
}

func PostPembelian(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    var pembelian PembelianDB
    c.Bind(&pembelian)
    if pembelian.Sku != "" && pembelian.Jpemesanan != 0 && pembelian.Jditerima != 0 && pembelian.Hbeli != 0 && pembelian.Kwitansi != "" && pembelian.Catatan != "" {
        pembelian.Waktu = time.Now().In(loc)
        pembelian.Waktuupdt = time.Now().In(loc)
        err := db.Create(&pembelian)

        result := Pembelians{
                PembelianDB : pembelian,
                Total       : pembelian.Jditerima * pembelian.Hbeli,
            }

        if err.RowsAffected == 0 {
            c.JSON(422, gin.H{"error": "Error occured"})
        } else {
            c.JSON(201, gin.H{"success": result})
        }
    } else {
        c.JSON(422, gin.H{"error": "Fields are empty"})
    }
}

func GetPembelians(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    var pembelians []Pembelians
    db.Select("pb.*, (pb.jditerima * pb.hbeli) as total").Table("pembelian_dbs pb").Scan(&pembelians);

    if len(pembelians)==0 {
        c.JSON(404, gin.H{"error": "Pembelian not found"})
    } else {
        c.JSON(200, pembelians)
    }
}

func GetPembelian(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    id := c.Params.ByName("id")
    var pembelian PembelianDB
    db.First(&pembelian, id)
    result  := Pembelians{
        PembelianDB : pembelian,
        Total       : pembelian.Jditerima * pembelian.Hbeli,
    }
    fmt.Println(result.Total)
    if result.Id != 0 {
        c.JSON(200, result)
    } else {
        c.JSON(404, gin.H{"error": "Pembelian not found"})
    }
}

func UpdatePembelian(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    id := c.Params.ByName("id")
    var pembelian PembelianDB
    db.First(&pembelian, id)

    if pembelian.Sku != "" {
        if pembelian.Id != 0 {
            var newPembelian PembelianDB
            c.Bind(&newPembelian)
            newPembelian.Waktuupdt = time.Now().In(loc)

            inp := PembelianDB{
                Id          : pembelian.Id,
                Waktu       : pembelian.Waktu,
                Waktuupdt   : newPembelian.Waktuupdt,
                Sku         : stringComp(newPembelian.Sku,pembelian.Sku),
                Jpemesanan  : intComp(newPembelian.Jpemesanan,pembelian.Jpemesanan),
                Jditerima   : intComp(newPembelian.Jditerima,pembelian.Jditerima),
                Hbeli       : intComp(newPembelian.Hbeli,pembelian.Hbeli),
                Kwitansi    : stringComp(newPembelian.Kwitansi,pembelian.Kwitansi),
                Catatan     : stringComp(newPembelian.Catatan,pembelian.Catatan),
            }

            err := db.Model(&pembelian).Updates(&inp)
            result := Pembelians{
                PembelianDB : inp,
                Total       : inp.Jditerima * inp.Hbeli,
            }

            if err.RowsAffected == 0 {
                c.JSON(422, gin.H{"error": err})
            } else {
                c.JSON(200, gin.H{"success": result})
            }
        } else {
            c.JSON(404, gin.H{"error": "Pembelian not found"})
        }

    } else {
        c.JSON(422, gin.H{"error": "Fields are empty"})
    }

}

func DeletePembelian(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    id := c.Params.ByName("id")
    var pembelian PembelianDB
    db.First(&pembelian, id)

    if pembelian.Id != 0 {
        db.Delete(&pembelian)
        c.JSON(200, gin.H{"success": "Pembelian #" + id + " deleted"})
    } else {
        c.JSON(404, gin.H{"error": "Pembelian not found"})
    }
}

func CheckStock (c *gin.Context) {
    db := InitDb()
    defer db.Close()
    var inpBarang PenjualanDB
    c.Bind(&inpBarang)

    var barang Listbarang
    var sumPB = "(SELECT COALESCE(sum(pb.jditerima), 0) as jumlah FROM pembelian_dbs pb where pb.sku = b.sku)"
    var sumPJ = "(SELECT COALESCE(sum(pj.jkeluar), 0) as jumlah FROM penjualan_dbs pj where pj.sku = b.sku)"
    db.Select("b.*, "+sumPB+" - "+sumPJ+" as jumlah").Table("barang_dbs b").Group("b.sku").Where("b.sku = ?", inpBarang.Sku).Order("b.sku ASC").Scan(&barang);
    
    if barang.Jumlah >= inpBarang.Jkeluar {
        c.Set("req",inpBarang)
        c.Next()
    } else {
        c.JSON(404, gin.H{"error": "Out of Stock"})
    }
}

func PostPenjualan(c *gin.Context) {
    db := InitDb()
    defer db.Close()
    penjualan := c.MustGet("req").(PenjualanDB)

    if penjualan.Sku != "" && penjualan.Jkeluar != 0 && penjualan.Hjual != 0{
        penjualan.Waktu = time.Now().In(loc)
        penjualan.Waktuupdt = time.Now().In(loc)
        err := db.Create(&penjualan)

        result := Penjualans{
            PenjualanDB : penjualan,
            Total       : penjualan.Jkeluar * penjualan.Hjual,
        }

        if err.RowsAffected == 0 {
            c.JSON(422, gin.H{"error": "Error occured"})
        } else {
            c.JSON(201, gin.H{"success": result})
        }
    } else {
        c.JSON(422, gin.H{"error": "Fields are empty"})
    }
}

func GetPenjualans(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    var penjualans []Penjualans
    db.Select("pj.*, pj.jkeluar * pj.hjual as total").Table("penjualan_dbs pj").Scan(&penjualans);
    // db.Find(&penjualans)

    if len(penjualans)==0 {
        c.JSON(404, gin.H{"error": "Penjualan not found"})
    } else {
        c.JSON(200, penjualans)
    }
}

func GetLaporanPenjualans(c *gin.Context) {
    db := InitDb()
    defer db.Close()
    tmfirst := c.Params.ByName("tmfirst")
    tmlast := c.Params.ByName("tmlast")
    var penjualans []Laporanpenjualan
    var sumPBTotal = "(SELECT COALESCE(sum(pb.jditerima * pb.hbeli), 0) as jumlah FROM pembelian_dbs pb where pb.sku = pj.sku)"
    var sumPB = "(SELECT COALESCE(sum(pb.jditerima), 0) as jumlah FROM pembelian_dbs pb where pb.sku = pj.sku)"
    db.Select("pj.*, b.nama as nama, (pj.jkeluar * pj.hjual) as total, "+sumPBTotal+" / "+sumPB+" as hbeli, (pj.hjual - "+sumPBTotal+" / "+sumPB+") * pj.jkeluar as laba").Table("penjualan_dbs pj, barang_dbs b").Where("b.sku = pj.sku").Where("pj.waktu BETWEEN ? AND ?",tmfirst,tmlast).Scan(&penjualans);

    var tomzet,tlabakotor,tpenjualan int

    for i := 0; i < len(penjualans); i++ {
        tomzet += penjualans[i].Total
        tlabakotor += penjualans[i].Laba
        if strings.ContainsAny(penjualans[i].Catatan,"Pesanan"){
            penjualans[i].Catatan = strings.Replace(penjualans[i].Catatan, "Pesanan ", "", -1)
            tpenjualan += 1
        } else {
            penjualans[i].Catatan = ""
        }
    }

    if len(penjualans)==0 {
        c.JSON(404, gin.H{"error": "Penjualan not found"})
    } else {
        csv := c.Params.ByName("csv")
        if csv == "csv"{
            c.Set("data",penjualans)
            c.Set("tipe","penjualan")
            c.Next()
        } else {
            c.JSON(200, gin.H{"result": gin.H{"summary":gin.H{"tomzet":tomzet,"tlabakotor":tlabakotor,"tpenjualan":tpenjualan,"tbarang":len(penjualans)},"detail":penjualans}})
        }
    }

}

func GetPenjualan(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    id := c.Params.ByName("id")
    var penjualan PenjualanDB
    db.First(&penjualan, id)

    result  := Penjualans{
        PenjualanDB : penjualan,
        Total       : penjualan.Jkeluar * penjualan.Hjual,
    }
    fmt.Println(result.Total)
    if result.Id != 0 {
        c.JSON(200, result)
    } else {
        c.JSON(404, gin.H{"error": "Penjualan not found"})
    }
}

func UpdatePenjualan(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    id := c.Params.ByName("id")
    
    var penjualan PenjualanDB
    db.First(&penjualan, id)

    if penjualan.Sku != "" {
        if penjualan.Id != 0 {
            newPenjualan := c.MustGet("req").(PenjualanDB)
            newPenjualan.Waktuupdt = time.Now().In(loc)
            
            inp := PenjualanDB{
                Id          : penjualan.Id,
                Waktu       : penjualan.Waktu,
                Waktuupdt   : newPenjualan.Waktuupdt,
                Sku         : stringComp(newPenjualan.Sku, penjualan.Sku),
                Jkeluar     : intComp(newPenjualan.Jkeluar, penjualan.Jkeluar),
                Hjual       : intComp(newPenjualan.Hjual, penjualan.Hjual),
                Catatan     : newPenjualan.Catatan,
            }

            err := db.Model(&penjualan).Updates(&inp)

            result := Penjualans{
                PenjualanDB : inp,
                Total       : inp.Jkeluar * inp.Hjual,
            }

            if err.RowsAffected == 0 {
                c.JSON(422, gin.H{"error": "Error occured"})
            } else {
                c.JSON(200, gin.H{"success": result})
            }
        } else {
            c.JSON(404, gin.H{"error": "Penjualan not found"})
        }

    } else {
        c.JSON(422, gin.H{"error": "Fields are empty"})
    }

}

func DeletePenjualan(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    id := c.Params.ByName("id")
    var penjualan PenjualanDB
    db.First(&penjualan, id)

    if penjualan.Id != 0 {
        db.Delete(&penjualan)
        c.JSON(200, gin.H{"success": "Penjualan #" + id + " deleted"})
    } else {
        c.JSON(404, gin.H{"error": "Penjualan not found"})
    }
}

func ExportCSV(c *gin.Context) {
    tipe := c.MustGet("tipe").(string)
    b := &bytes.Buffer{}
    w := csv.NewWriter(b)
    if tipe == "nilaibarang"{
        rawdata := c.MustGet("data").([]Laporannilaibarang)
        for _, data := range rawdata {
            var record []string
            record = append(record, data.Sku)
            record = append(record, data.Nama)
            record = append(record, strconv.Itoa(data.Jumlah))
            record = append(record, strconv.Itoa(data.Ratarata))
            record = append(record, strconv.Itoa(data.Total))
            if err := w.Write(record); err != nil {
                log.Fatalln("error writing record to csv:", err)
            }
        }
    } else if tipe == "penjualan"{
        rawdata := c.MustGet("data").([]Laporanpenjualan)
        for _, data := range rawdata {
            var record []string
            record = append(record, data.Catatan)
            record = append(record, data.Waktu.String())
            record = append(record, data.Sku)
            record = append(record, data.Nama)
            record = append(record, strconv.Itoa(data.Jkeluar))
            record = append(record, strconv.Itoa(data.Hjual))
            record = append(record, strconv.Itoa(data.Total))
            record = append(record, strconv.Itoa(data.Hbeli))
            record = append(record, strconv.Itoa(data.Laba))
            if err := w.Write(record); err != nil {
                log.Fatalln("error writing record to csv:", err)
            }
        }
    }

    w.Flush()

    if err := w.Error(); err != nil {
        log.Fatal(err)
    }
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename=Laporan"+tipe+".csv")
    c.Data(http.StatusOK, "text/csv", b.Bytes())
}
