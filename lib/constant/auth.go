package constant

const (
	SendVerificationDelayKeyPrefix = "send-verification-delay:%d:%s" // send-verification-delay:[identifier]:[verification_type]
	MfaFlagKeyPrefix               = "mfa-flag:%d"                   // mfa-flag:[identifier]
	AccessTokenKeyPrefix           = "access-token:%s"               // access-token:[access_token]
	OtpTokenKeyPrefix              = "otp-token:%s"                  // otp-token:[otp_token]
	SendOtpCtrKeyPrefix            = "send-otp-ctr:%d:%s"            // send-otp-ctr:[identifier]:[otp_type]
	SendOtpDelayKeyPrefix          = "send-otp-delay:%d:%s"          // send-otp-delay:[identifier]:[otp_type]

	DefaultIssuer   = "DefaultIssuer"
	DefaultAudience = "DefaultAudience"

	OtpTypeLogin = "LOGIN"

	OtpChannelEmail = "EMAIL"
)
