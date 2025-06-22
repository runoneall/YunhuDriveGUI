package structs

type QiNiuHostResponse struct {
	Hosts []struct {
		Region string `json:"region"`
		TTL    int    `json:"ttl"`
		IO     struct {
			Domains []string `json:"domains"`
		} `json:"io"`
		IOSrc struct {
			Domains []string `json:"domains"`
		} `json:"io_src"`
		Up struct {
			Domains []string `json:"domains"`
			Old     []string `json:"old"`
		} `json:"up"`
		UC struct {
			Domains []string `json:"domains"`
		} `json:"uc"`
		RS struct {
			Domains []string `json:"domains"`
		} `json:"rs"`
		RSF struct {
			Domains []string `json:"domains"`
		} `json:"rsf"`
		API struct {
			Domains []string `json:"domains"`
		} `json:"api"`
		S3 struct {
			Domains     []string `json:"domains"`
			RegionAlias string   `json:"region_alias"`
		} `json:"s3"`
	} `json:"hosts"`
	TTL int `json:"ttl"`
}

type QiNiuUploadResponse struct {
	Key    string      `json:"key"`
	Hash   string      `json:"hash"`
	Fsize  int64       `json:"fsize"`
	Avinfo interface{} `json:"avinfo"`
}
