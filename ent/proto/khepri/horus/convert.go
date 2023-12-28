package horus

import ent "khepri.dev/horus/ent"

func ToProtoAccountList(e []*ent.Account) ([]*Account, error) {
	return toProtoAccountList(e)
}

func ToProtoSilo(e *ent.Silo) (*Silo, error) {
	return toProtoSilo(e)
}
