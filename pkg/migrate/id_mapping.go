package migrate

import (
	"github.com/guregu/null"
	"path/filepath"
	"strings"
)

type NamingKind int

const (
	NamingKindUUID NamingKind = iota
	NamingKindWxID
	NamingKindDevice
)

var dirNames = [...]string{
	"user-id",
	"wechat",
	"device-token",
}

func (x NamingKind) String() string {
	if x > NamingKindDevice || x < NamingKindUUID {
		return ""
	}

	return dirNames[x]
}

const StmtSaveMapping = `
INSERT INTO premium.apple_id_mapping
SET original_transaction_id = :original_transaction_id,
	ftc_id = :ftc_id,
	device_token = :device_token,
	wx_union_id = :wx_union_id
ON DUPLICATE KEY UPDATE
	ftc_id = IFNULL(ftc_id, :ftc_id),
	device_token = IFNULL(device_token, :device_token),
	wx_union_id = IFNULL(wx_union_id, :wx_union_id)`

type IDMapping struct {
	TxID        string      `db:"original_transaction_id"`
	FtcID       null.String `db:"ftc_id"`
	DeviceToken null.String `db:"device_token"`
	UnionID     null.String `db:"wx_union_id"`
	AbsFilePath string
}

func NewIDMapping(fileName string, k NamingKind) IDMapping {
	m := IDMapping{
		TxID:        "",
		FtcID:       null.String{},
		DeviceToken: null.String{},
		UnionID:     null.String{},
		AbsFilePath: fileName,
	}

	baseName := filepath.Base(fileName)
	ext := filepath.Ext(fileName)

	id := strings.TrimSuffix(baseName, ext)

	switch k {
	case NamingKindUUID:
		m.FtcID = null.StringFrom(id)
	case NamingKindWxID:
		m.UnionID = null.StringFrom(id)
	case NamingKindDevice:
		m.DeviceToken = null.StringFrom(id)
	}

	return m
}
