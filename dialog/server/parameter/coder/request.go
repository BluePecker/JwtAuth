package coder

type (
	Encode struct {
		Device string `json:"device"`
		Addr   string `json:"addr"`
		Unique string `json:"unique"`
	}

	Decode struct {
		JsonWebToken string `json:"jwt"`
	}
)
