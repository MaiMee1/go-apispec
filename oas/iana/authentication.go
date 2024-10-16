package iana

type AuthScheme string

const (
	BasicScheme       AuthScheme = "Basic"
	BearerScheme      AuthScheme = "Bearer"
	DigestScheme      AuthScheme = "Digest"
	HobaScheme        AuthScheme = "HOBA"
	MutualScheme      AuthScheme = "Mutual"
	NegotiateScheme   AuthScheme = "Negotiate"
	VapidScheme       AuthScheme = "vapid"
	ScramSha1Scheme   AuthScheme = "SCRAM-SHA-1"
	ScramSha256Scheme AuthScheme = "SCRAM-SHA-256"
)
