package models

type ShortLink struct {
	Url        string `json:"url"`
	PartnerUrl string `json:"partner_url"`
}

type ShortLinksResult struct {
	Links []ShortLink `json:"links"`
}

type ShortLinksResponse struct {
	Status string           `json:"code"`
	Result ShortLinksResult `json:"result"`
}
