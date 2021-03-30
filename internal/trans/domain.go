package trans

import (
	"github.com/opendigitalpay-io/open-pay/internal/domain"
	"github.com/opendigitalpay-io/open-pay/internal/storage"
	"github.com/opendigitalpay-io/open-pay/internal/tcc"
)

type Transfer struct {
	ID            uint64
	OrderID       uint64
	SourceID      uint64
	DestinationID uint64
	Type          Type
	Amount        int64
	Currency      string
	Status        domain.STATUS
	Metadata      map[string]interface{}
	CreatedAt     int64
	UpdatedAt     int64
	// FIXME: try to split these field out into another struct
	transferTxns []tcc.Interface
	repository   storage.Repository
}

type Type string

const (
	ORDER  Type = "ORDER"
	TOP_UP Type = "TOP_UP"
)

var types = [...]string{
	"ORDER",
	"TOP_UP",
}

func (t *Type) String() string {
	x := string(*t)
	for _, v := range types {
		if v == x {
			return x
		}
	}
	return ""
}

type TransferDTO struct {
	OrderID       uint64
	SourceID      uint64
	DestinationID uint64
	Type          Type
	Amount        int64
	Currency      string
	Metadata      map[string]interface{}
}

func Create(dto TransferDTO) Transfer {
	return Transfer{
		OrderID:       dto.OrderID,
		SourceID:      dto.SourceID,
		DestinationID: dto.DestinationID,
		Type:          dto.Type,
		Amount:        dto.Amount,
		Status:        domain.CREATED,
		Currency:      dto.Currency,
		Metadata:      dto.Metadata,
	}
}

func (t *Transfer) OnTrySuccessCallback() {

}

func (t *Transfer) OnTryFailCallback() {

}

func (t *Transfer) OnCommitSuccessCallback() {

}

func (t *Transfer) OnCommitFailCallback() {

}

func (t *Transfer) OnCancelSuccessCallback() {

}
func (t *Transfer) OnCancelFailCallback() {

}
