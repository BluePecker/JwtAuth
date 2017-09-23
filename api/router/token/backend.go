package token

import "github.com/BluePecker/JwtAuth/api/types/token"

type Backend interface {
    // 生成jwt
    Generate(req token.GenerateRequest) (string, error)
    
    // 校验jwt
    Auth(req token.AuthRequest) (interface{}, error)
}