package bare

import (
	horus "khepri.dev/horus"
	ent "khepri.dev/horus/ent"
)

func ToProtoToken(e *ent.Token) (*horus.Token, error) {
	return toProtoToken(e)
}
