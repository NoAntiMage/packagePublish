package constant

const Secret = "574dcac971c9c242701c0778608a5110"
const Password = "123"

const (
	LoginTokenKey = "PackageServer:LoginToken:"
	JwtTokenKey   = "PackageServer:Jwt:"
)

type direction string

const (
	To   direction = "to"
	From direction = "from"
)
