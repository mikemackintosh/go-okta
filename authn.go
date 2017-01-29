package okta

import (
	"time"
)

type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorSummary string `json:"errorSummary"`
	ErrorLink    string `json:"errorLink"`
	ErrorID      string `json:"errorId"`
	ErrorCauses  []struct {
		ErrorSummary string `json:"errorSummary"`
	} `json:"errorCauses"`
}

type AuthnRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RelayState string `json:"relayState"`
	Options    struct {
		MultiOptionalFactorEnroll bool `json:"multiOptionalFactorEnroll"`
		WarnBeforePasswordExpired bool `json:"warnBeforePasswordExpired"`
	} `json:"options"`
}

type AuthnResponse struct {
	StateToken string    `json:"stateToken"`
	ExpiresAt  time.Time `json:"expiresAt"`
	Status     string    `json:"status"`
	RelayState string    `json:"relayState"`
	Embedded   struct {
		User struct {
			ID              string    `json:"id"`
			PasswordChanged time.Time `json:"passwordChanged"`
			Profile         struct {
				Login     string `json:"login"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Locale    string `json:"locale"`
				TimeZone  string `json:"timeZone"`
			} `json:"profile"`
		} `json:"user"`
		Factors []struct {
			ID         string `json:"id"`
			FactorType string `json:"factorType"`
			Provider   string `json:"provider"`
			VendorName string `json:"vendorName"`
			Profile    struct {
				CredentialID string `json:"credentialId"`
			} `json:"profile"`
			Links struct {
				Verify struct {
					Href  string `json:"href"`
					Hints struct {
						Allow []string `json:"allow"`
					} `json:"hints"`
				} `json:"verify"`
			} `json:"_links"`
		} `json:"factors"`
		Policy struct {
			AllowRememberDevice             bool `json:"allowRememberDevice"`
			RememberDeviceLifetimeInMinutes int  `json:"rememberDeviceLifetimeInMinutes"`
			RememberDeviceByDefault         bool `json:"rememberDeviceByDefault"`
		} `json:"policy"`
	} `json:"_embedded"`
	Links struct {
		Cancel struct {
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"cancel"`
	} `json:"_links"`
}
