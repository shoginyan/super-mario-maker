package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func IsObjectAvailable(dataID uint64) uint32 {
	var underReview bool

	err := database.Postgres.QueryRow(`SELECT under_review FROM datastore.objects WHERE data_id=$1 AND upload_completed=TRUE AND deleted=FALSE`, dataID).Scan(&underReview)
	if err != nil {
		if err == sql.ErrNoRows {
			return nex.ResultCodes.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.ResultCodes.DataStore.Unknown
	}

	if underReview {
		return nex.ResultCodes.DataStore.UnderReviewing
	}

	return 0
}
