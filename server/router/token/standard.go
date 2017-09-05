package token

type Standard interface {
    // 生成jwt
    Generate(userId, device, address string) (string, error)
    
    // 校验jwt
    Auth(jwt string) (string, error)
    
    // 更新jwt
    Upgrade(jwt string) (string, error)
}