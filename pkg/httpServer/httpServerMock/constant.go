package httpServerMock

const (
	envJwtRsa256PublicKey           = "JWT_RSA256_PUBLIC_KEY"
	envJwtRsa256PrivateKey          = "JWT_RSA256_PRIVATE_KEY"
)

const (
	tokenHeaderName = "Authorization"
	tokenKind       = "Bearer"
	tokenDelimiter  = " "
)

// Faker RSA private and public key use in unit test
const (
	jwtRsa256PrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgGzNqUxynJYXEOGGr7b5ag4RJ0hooOHj6LrLJ49iC1Yrhgbumsvw
itkEUyke7Tptmjyvo3puKUbzxmtjygNcfZAUfgEOpIheXJSkrf3jIw3B8q+pxXzu
lCsNUwiew0ku/AGB968su3bvX0wuUI99dx2GeeABiiPR+322JGbrLK/PAgMBAAEC
gYBYBdpdequr0WVazzgA868VSlNZhSwDE/sIg6qxmURKplN78DVToHr0L0eIEPkj
N+B8ECxVtCG4wSdZYhXgukbp48bT1C/Q5rH3vdnSJ0HAqK3u6u/Off3CazWdAhNr
H9Pyyh0s/Js3PtAljCoJfGQJgePlqhr8M2S24VLmWESaoQJBAMRyMOMMwzUHBeLt
LpH1Apb+sbvuUN5mIhuAJDhDShSy2Xe9+rFor9hRKXOyoECb04PsMGbAmvMMf/7a
KqxoTTsCQQCNya+N+jSlXum0cEq+iH0uyLFsmdhPs1qT9c9UtdE8oSyvPWXxz3xH
eh2SHlxxkzO7dCzN3kQkXG30vKFutU59AkB/2spSnBXYx29fWHs857f9ylqnM95S
QSrltyrbq3/lpNnBA1bMbJQ2N+zArnt6UXECpZCC78xpb7NGjXvEpkXNAkBfD6Ce
/Oh9EzR3IG5MbjAXxMCHwmGvld0dpElcTwY4swrFdtG5nNWDCpF23No332xouowr
fmCLTFkAI9PL6Mz1AkAfb0nmtUqOBFmtrcOgAJpDIzG6TlWAezn0tpzzVgV9Y/Ph
OVdFuWal6LmxsKQi+RsOm+3TvVpjOZIszNCiqT6B
-----END RSA PRIVATE KEY-----
`
	jwtRsa256PublicKey = `
-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgGzNqUxynJYXEOGGr7b5ag4RJ0ho
oOHj6LrLJ49iC1YrhgbumsvwitkEUyke7Tptmjyvo3puKUbzxmtjygNcfZAUfgEO
pIheXJSkrf3jIw3B8q+pxXzulCsNUwiew0ku/AGB968su3bvX0wuUI99dx2GeeAB
iiPR+322JGbrLK/PAgMBAAE=
-----END PUBLIC KEY-----
`
)
