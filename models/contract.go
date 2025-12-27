package models

import "encoding/xml"

// ContractDetail represents the complete contract detail
type ContractDetail struct {
	XMLName            xml.Name            `xml:"http://gsph.sub.com/cust/types contractDetail"`
	Contract           Contract            `xml:"contract"`
	ContractAttributes *ContractAttributes `xml:"contractAttributes,omitempty"`
	ContractNo         int                 `xml:"contractno,omitempty"`
	Person             *Person             `xml:"person,omitempty"`
	StdAddr            *StdAddr            `xml:"stdAddr,omitempty"`
	PaymentInfo        *PaymentInfo        `xml:"paymentInfo,omitempty"`
	Counting           int                 `xml:"counting,omitempty"`
	Present            int                 `xml:"present,omitempty"`
	Status             int                 `xml:"status,omitempty"`
	Delete             int                 `xml:"delete,omitempty"`
	Memo               string              `xml:"memo,omitempty"`
	InvoiceGroup       int                 `xml:"invoicegroup,omitempty"`
	TaxIDNo            string              `xml:"taxIdNo,omitempty"`
	IDNo               string              `xml:"idNo,omitempty"`
}

// Contract represents the basic contract info
type Contract struct {
	Href       string `xml:"href,attr,omitempty"`
	ID         *int   `xml:"id,omitempty"`
	Name       string `xml:"name"`
	ValidFrom  string `xml:"xValidFrom"`
	ValidUntil string `xml:"xValidUntil"`
	FilialID   string `xml:"filialId,omitempty"`
}

// ContractAttributes represents contract attributes
type ContractAttributes struct {
	AutoBlockDays     int    `xml:"autoBlockDays"`
	PrePayment        int    `xml:"prePayment"`
	DiscountCash      int    `xml:"discountCash"`
	Discount          int    `xml:"discount"`
	DiscountValue     int    `xml:"discountValue"`
	FeeType           int    `xml:"feeType"`
	FlatFee           int    `xml:"flatFee"`
	VAT               string `xml:"vat"`
	FlatFeeType       int    `xml:"flatFeeType"`
	FlatFeeCalc       int    `xml:"flatFeeCalc"`
	FlatFeeFirstMonth int    `xml:"flatFeeFirstMonth"`
	FlatFeeLastMonth  int    `xml:"flatFeeLastMonth"`
}

// Person represents person information
type Person struct {
	Title        string `xml:"title"`
	FirstName    string `xml:"firstName"`
	Birthday     string `xml:"birthday"`
	Lang         int    `xml:"lang"`
	ContractLang int    `xml:"contractLang"`
	MatchCode    string `xml:"matchCode,omitempty"`
}

// StdAddr represents standard address
type StdAddr struct {
	Street  string `xml:"street"`
	Town    string `xml:"town"`
	Postbox string `xml:"postbox"`
}

// PaymentInfo represents payment information
type PaymentInfo struct {
	PaymentType int       `xml:"paymentType"`
	BankAcct    *BankAcct `xml:"bankacct,omitempty"`
	CCNo1       *CCInfo   `xml:"ccno1,omitempty"`
	CCNo2       *CCInfo   `xml:"ccno2,omitempty"`
}

// BankAcct represents bank account information
type BankAcct struct {
	AcctNo                 string `xml:"acctno"`
	BankCode               string `xml:"bankcode"`
	BankName               string `xml:"bankname"`
	SEPAMandatRef          string `xml:"sepamandatref"`
	SEPAMandatDate         string `xml:"sepamandatdate"`
	SEPAPreNotification    int    `xml:"sepaprenotification"`
	SEPAIBAN               string `xml:"sepaiban"`
	SEPACreditorIdentifier string `xml:"sepacreditoridentifier"`
}

// CCInfo represents credit card information
type CCInfo struct {
	CardType       int    `xml:"cardtype"`
	CardNo         string `xml:"cardno"`
	CardName       string `xml:"cardname"`
	CardValidUntil string `xml:"cardvaliduntil"`
}

// CreateContractRequest for creating a new contract
type CreateContractRequest struct {
	ID         *int   // Optional - nil if 3rd party should generate
	Name       string // Required
	ValidFrom  string // Required - Format: "2021-01-01"
	ValidUntil string // Required - Format: "2021-12-31"
}

// ToXML converts CreateContractRequest to ContractDetail for XML marshaling
func (r CreateContractRequest) ToXML() ContractDetail {
	return ContractDetail{
		Contract: Contract{
			ID:         r.ID,
			Name:       r.Name,
			ValidFrom:  r.ValidFrom,
			ValidUntil: r.ValidUntil,
		},
	}
}

// ErrorResponse represents XML error response from API
type ErrorResponse struct {
	XMLName xml.Name `xml:"http://gsph.sub.com/cust/types errorResponse"`
	Error   struct {
		ErrCode      string `xml:"errCode"`
		ShortMsg     string `xml:"shortMsg"`
		Message      string `xml:"message"`
		CauseMessage string `xml:"causeMessage"`
	} `xml:"error"`
}
