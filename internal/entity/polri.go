package entity

// PolriPlateNuberInfo is struct represent data from korlatas get plate number data
type PolriPlateNuberInfo struct {
	WarnaTNKB          interface{} `json:"WarnaTNKB"`
	ModelKendaran      interface{} `json:"ModelKendaraan"`
	KodeJenisKendaraan interface{} `json:"KodeJenisKendaraan"`
	Warna              interface{} `json:"Warna"`
	KodeTNKB           interface{} `json:"KodeTNKB"`
	NIK                interface{} `json:"NIK"`
	Polda              interface{} `json:"Polda"`
	Status             interface{} `json:"Status"`
	Type               interface{} `json:"Type"`
	TahunPembuatan     interface{} `json:"TahunPembuatan"`
	KodeKabKota        interface{} `json:"KodeKabKota"`
	Keterangan         interface{} `json:"Keterangan"`
	Merk               interface{} `json:"Merk"`
	KodePolda          interface{} `json:"KodePolda"`
	IsiCylinder        interface{} `json:"IsiCylinder"`
	KodeModelKendaraan interface{} `json:"KodeModelKendaraan"`
	JenisKendaraan     interface{} `json:"JenisKendaraan"`
	KabKota            interface{} `json:"KabKota"`
}

//PolriCitizenResponseData struct
type PolriCitizenResponseData struct {
	StatusDesc string              `json:"StatusDesc"`
	Payload    PolriCitizenPayload `json:"Payload"`
	StatusCode string              `json:"StatusCode"`
	TraceID    string              `json:"trace_id"`
	Nik        string              `json:"nik"`
	WaktuAkses string              `json:"waktu_akses"`
}

// PolriCitizenPayload struct
type PolriCitizenPayload struct {
	Content          []PolriCitizenData `json:"content"`
	LastPage         bool               `json:"lastPage"`
	NumberOfElements int                `json:"numberOfElements"`
	Sort             interface{}        `json:"sort"`
	TotalElements    int                `json:"totalElements"`
	FirstPage        bool               `json:"firstPage"`
	Number           int                `json:"number"`
	Size             int                `json:"size"`
}

//PolriCitizenData struct
type PolriCitizenData struct {
	NamaLgkp    interface{} `json:"NAMA_LGKP"`
	Agama       interface{} `json:"AGAMA"`
	JenisPkrjn  interface{} `json:"JENIS_PKRJN"`
	PddkAkh     interface{} `json:"PDDK_AKH"`
	TmptLhr     interface{} `json:"TMPT_LHR"`
	StatusKawin interface{} `json:"STATUS_KAWIN"`
	GolDarah    interface{} `json:"GOL_DARAH"`
	JenisKlmin  interface{} `json:"JENIS_KLMIN"`
	Nik         interface{} `json:"NIK"`
	NoKk        interface{} `json:"NO_KK"`
	KabName     interface{} `json:"KAB_NAME"`
	NoRw        interface{} `json:"NO_RW"`
	KecName     interface{} `json:"KEC_NAME"`
	NoRt        interface{} `json:"NO_RT"`
	NoKel       interface{} `json:"NO_KEL"`
	Alamat      interface{} `json:"ALAMAT"`
	NoKec       interface{} `json:"NO_KEC"`
	NoProp      interface{} `json:"NO_PROP"`
	PropName    interface{} `json:"PROP_NAME"`
	NoKab       interface{} `json:"NO_KAB"`
	TglLhr      interface{} `json:"TGL_LHR"`
	KelName     interface{} `json:"KEL_NAME"`
	Foto        interface{} `json:"FOTO"`
	Respon      interface{} `json:"RESPON"`
}
