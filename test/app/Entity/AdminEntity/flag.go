package AdminEntity

import "github.com/dunpju/higo-gin/higo"

const (
	FlagDelete higo.Flag = iota + 1
	FlagLastLogin
	FlagUpdateInfo
	FlagUnifyUpdate
	FlagUpdateMobile
)
