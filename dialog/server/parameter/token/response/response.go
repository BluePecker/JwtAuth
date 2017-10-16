package response

type (
	Token struct {
		Device string  `json:"device"`
		Addr   string  `json:"addr"`
		TTL    float64 `json:"ttl"`
		Singed string  `json:"singed"`
	}
)
