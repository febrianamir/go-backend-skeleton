package constant

const (
	SendVerificationDelayKeyPrefix = "send-verification-delay:%d:%s" // send-verification-delay:[identifier]:[verification_type]
	MfaFlagKeyPrefix               = "mfa-flag:%d"                   // mfa-flag:[identifier]
	AccessTokenKeyPrefix           = "access-token:%s"               // access-token:[access_token]
	OtpTokenKeyPrefix              = "otp-token:%s"                  // otp-token:[otp_token]

	DefaultIssuer   = "DefaultIssuer"
	DefaultAudience = "DefaultAudience"
)
