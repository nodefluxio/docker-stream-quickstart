
 curl --location --request POST 'https://api.polri.go.id/korlantas/getbynopol' --header 'Authorization: Basic BASE64' --header 'Content-Type: application/json' --data-raw '{"nopol": "B 4720 KAD"}'

http code : 200
```json
[
  {
    "STATUS": "1",
    "Merk": "HONDA",
    "Type": "ASD ZXCASD",
    "TahunPembuatan": "2016",
    "IsiCylinder": "150.00",
    "Warna": "HITAM MERAH",
    "KodeJenisKendaraan": "4",
    "KodeModelKendaraan": "700",
    "KodeKabKota": "803",
    "KodePolda": "6",
    "KodeTNKB": "1",
    "NIK": "1234567891234",
    "JenisKendaraan": "SEPEDA MOTOR",
    "ModelKendaraan": "SPD MTR SOLO",
    "KabKota": "BEKASI, KOTA",
    "Polda": "METRO JAYA",
    "WarnaTNKB": "HITAM",
    "Keterangan": ""
  }
]
```

http code : 200
```json
[{ "STATUS": "0", "DESKRIPSI": "NRKB TIDAK DITEMUKAN DALAM DATABASE KAMI..." }]
```

http code : 400
```json
{"message":"nopol tidak valid, contoh nopol: B-1234-ABC atau B 1234 ABC atau B1234ABC atau B1234"}
```

http code : 401
```json
{"message":"Unauthorized"}
```

## Search By NIK

curl --location --request POST 'https://api.polri.go.id/dukcapil/getbynik2' --header 'Authorization: Basic BASE64' --header 'Content-Type: application/json' --data-raw '{"nik":"32712312312312321"}' -i

http code : 200
```json
{
  "StatusDesc": "Success",
  "Payload": {
    "content": [
      {
        "NAMA_LGKP": "Agus",
        "AGAMA": "ISLAM",
        "JENIS_PKRJN": "KARYAWAN SWASTA",
        "PDDK_AKH": "TAMAT SMA/SEDERAJAT",
        "TMPT_LHR": "CIAMIS",
        "STATUS_KAWIN": "KAWIN",
        "GOL_DARAH": "TIDAK TAHU",
        "JENIS_KLMIN": "Laki-Laki",
        "NIK": 312345678910,
        "NO_KK": 312345678910,
        "KAB_NAME": "KOTA BEKASI",
        "NO_RW": 10,
        "KEC_NAME": "JATIASIH",
        "NO_RT": 5,
        "NO_KEL": 1004,
        "ALAMAT": "PERUM PERUMAHAN",
        "NO_KEC": 9,
        "NO_PROP": 32,
        "PROP_NAME": "JAWA BARAT",
        "NO_KAB": 75,
        "TGL_LHR": "1994-11-15",
        "KEL_NAME": "JATIRASA",
        "FOTO": "b64"
      }
    ],
    "lastPage": true,
    "numberOfElements": 1,
    "sort": null,
    "totalElements": 1,
    "firstPage": true,
    "number": 0,
    "size": 1
  },
  "StatusCode": "00",
  "trace_id": "51cda30a6b00c3090c5682ef2565a58e",
  "nik": "312345678910",
  "waktu_akses": "20210715 05:04:26.311"
}
```

http code : 200
```json
{
  "StatusDesc": "Success",
  "Payload": {
    "content": [{ "RESPON": "Data Tidak Ditemukan" }],
    "lastPage": true,
    "numberOfElements": 1,
    "sort": null,
    "totalElements": 1,
    "firstPage": true,
    "number": 0,
    "size": 1
  },
  "StatusCode": "00",
  "trace_id": "525eba1b9427f4cfa50e554ec095c358",
  "nik": "3123456789101",
  "waktu_akses": "20210715 05:07:08.668"
}
```

http code : 400
```json
{
  "message":"request tidak valid"
}

```

http code : 401
```json
{"message":"Unauthorized"}
```