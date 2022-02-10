package entity

// SeagateGetTokenRequest struct
type SeagateGetTokenRequest struct {
	Embedding string `json:"embedding"`
	Limit     uint64 `json:"limit"`
}

//SeagateGetTokenResult struct
type SeagateGetTokenResult struct {
	Token string `json:"token"`
	OK    bool   `json:"ok"`
}

// SeagateFaceSearchResult public struct
type SeagateFaceSearchResult struct {
	OK      bool                 `json:"ok"`
	Message string               `json:"message"`
	Data    []*SeagateSearchData `json:"data"`
}

// SeagateSearchData public struct
type SeagateSearchData struct {
	Similiarity SeagateSimilarityData `json:"similarity"`
	Token       string                `json:"token"`
	Probability float64               `json:"probability"`
	Dukcapil    SeagateDukcapilData   `json:"dukcapil"`
}

// SeagateSimilarityData struct
type SeagateSimilarityData struct {
	Distance    float64 `json:"distance"`
	EmbeddingID string  `json:"embedding_id"`
}

// SeagateSimilarityData struct
type SeagateDukcapilStatus struct {
	Ok     bool   `json:"ok"`
	Respon string `json:"respon"`
}

// SeagateDukcapilData struct
type SeagateDukcapilData struct {
	Ok                 bool   `json:"ok,omitempty"`
	Respon             string `json:"respon,omitempty"`
	NamaLengkap        string `json:"nama_lgkp"`
	StatusHubKel       string `json:"stat_hbkel"`
	Agama              string `json:"agama"`
	JenisPekerjaan     string `json:"jenis_pkrjn"`
	PendidikanAkhir    string `json:"pddk_akh"`
	TempatLahir        string `json:"tmpt_lhr"`
	StatusKawin        string `json:"status_kawin"`
	GolonganDarah      string `json:"gol_darah"`
	NikIbu             int64  `json:"nik_ibu,omitempty"`
	JenisKelamin       string `json:"jenis_klmin"`
	NomorKartuKeluarga int64  `json:"no_kk"`
	NIK                int64  `json:"nik"`
	Foto               string `json:"foto"`
	NamaKabupaten      string `json:"kab_name"`
	NamaLengkapAyah    string `json:"nama_lgkp_ayah"`
	NomorRW            int    `json:"no_rw"`
	NamaKecamatan      string `json:"kec_name"`
	NomorRT            int    `json:"no_rt"`
	NomorKelurahan     int    `json:"no_kel"`
	Alamat             string `json:"alamat"`
	NomorKecamatan     int    `json:"no_kec"`
	NikAyah            int64  `json:"nik_ayah,omitempty"`
	NomorProvinsi      int    `json:"no_prop"`
	NamaLengkapIBU     string `json:"nama_lgkp_ibu"`
	NamaProvinsi       string `json:"prop_name"`
	NomorKabupaten     int    `json:"no_kab"`
	TanggalLahir       string `json:"tgl_lhr"`
	NamaKelurahan      string `json:"kel_name"`
}
