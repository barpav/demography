package models

const MimeTypeEnrichedPersonDataV1 = "application/vnd.enrichedPersonData.v1+json"

// Schema: enrichedPersonData.v1
type EnrichedPersonDataV1 struct {
	Id         int64  `json:"id"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Country    string `json:"country,omitempty"`
}
