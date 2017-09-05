package token

type(
    GenerateRequest struct {
        Device string `json:"device"`
        Addr   string `json:"addr"`
        Unique string `json:"unique"`
    }
    
    AuthRequest struct {
        JsonWebToken string `json:"jwt"`
    }
)
