package recaptcha

type MockVerifier struct {
	ShouldPass bool
	Err        error
}

func (m *MockVerifier) Verify(token string) (bool, error) {
	if m.Err != nil {
		return false, m.Err
	}
	return m.ShouldPass, nil
}
