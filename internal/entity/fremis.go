package entity

// FRemisCandidate represents data format from FRemis API
type FRemisCandidate struct {
	FaceID     string  `json:"face_id"`
	Similarity float64 `json:"similarity"`
	Variation  string  `json:"variation"`
}

// FaceEnrollment represent data format response from FRemis face enrollmetn API
type FaceEnrollment struct {
	FaceID    string `json:"face_id"`
	Variation string `json:"variation"`
}

type FaceEmbedings struct {
	Embeddings [][]float64 `json:"embeddings"`
}
