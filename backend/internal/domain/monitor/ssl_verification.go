package monitor

import "time"

type SslVerification struct {
	isEnabled          bool
	verifyExpirationIn time.Duration
}

func (s *SslVerification) IsEnabled() bool {
	return s.isEnabled
}

func (s *SslVerification) VerifyExpirationIn() time.Duration {
	return s.verifyExpirationIn
}
