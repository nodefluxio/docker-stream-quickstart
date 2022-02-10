package presenter

import "gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"

// PolriSearchPlateResponse is struct represent data from korlatas get plate number data
type PolriSearchPlateResponse struct {
	WarnaTNKB          interface{} `json:"warna_tnkb"`
	ModelKendaran      interface{} `json:"model_kendaraan"`
	KodeJenisKendaraan interface{} `json:"kode_jenis_kendaraan"`
	Warna              interface{} `json:"warna"`
	KodeTNKB           interface{} `json:"kode_tnkb"`
	NIK                interface{} `json:"nik"`
	Polda              interface{} `json:"polda"`
	Status             interface{} `json:"status"`
	Type               interface{} `json:"type"`
	TahunPembuatan     interface{} `json:"tahun_pembuatan"`
	KodeKabKota        interface{} `json:"kode_kab_kota"`
	Keterangan         interface{} `json:"keterangan"`
	Merk               interface{} `json:"merk"`
	KodePolda          interface{} `json:"kode_polda"`
	IsiCylinder        interface{} `json:"isi_cylinder"`
	KodeModelKendaraan interface{} `json:"kode_model_kendaraan"`
	JenisKendaraan     interface{} `json:"jenis_kendaraan"`
	KabKota            interface{} `json:"kab_kota"`
}

//PolriCitizenResponse struct
type PolriCitizenResponse struct {
	NamaLgkp    interface{} `json:"nama_lgkp"`
	Agama       interface{} `json:"agama"`
	JenisPkrjn  interface{} `json:"jenis_pkrjn"`
	PddkAkh     interface{} `json:"pddk_akh"`
	TmptLhr     interface{} `json:"tmpt_lhr"`
	StatusKawin interface{} `json:"status_kawin"`
	GolDarah    interface{} `json:"gol_darah"`
	JenisKlmin  interface{} `json:"jenis_klmin"`
	Nik         interface{} `json:"nik"`
	NoKk        interface{} `json:"no_kk"`
	KabName     interface{} `json:"kab_name"`
	NoRw        interface{} `json:"no_rw"`
	KecName     interface{} `json:"kec_name"`
	NoRt        interface{} `json:"no_rt"`
	NoKel       interface{} `json:"no_kel"`
	Alamat      interface{} `json:"alamat"`
	NoKec       interface{} `json:"no_kec"`
	NoProp      interface{} `json:"no_prop"`
	PropName    interface{} `json:"prop_name"`
	NoKab       interface{} `json:"no_kab"`
	TglLhr      interface{} `json:"tgl_lhr"`
	KelName     interface{} `json:"kel_name"`
	Foto        interface{} `json:"foto"`
	Respon      interface{} `json:"respon,omitempty"`
}

// PolriFaceResultResponse public struct
type PolriFaceResultResponse struct {
	Similiarity    entity.SeagateSimilarityData `json:"similarity"`
	Token          string                       `json:"token"`
	Probability    float64                      `json:"probability"`
	DukcapilStatus entity.SeagateDukcapilStatus `json:"dukcapil_status"`
	DukcapilData   *entity.SeagateDukcapilData  `json:"dukcapil_data"`
}
