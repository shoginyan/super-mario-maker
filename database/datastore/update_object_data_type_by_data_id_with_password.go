package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func UpdateObjectDataTypeByDataIDWithPassword(dataID uint64, dataType uint16, password uint64) uint32 {
	var updatePassword uint64
	var underReview bool

	err := database.Postgres.QueryRow(`SELECT update_password, under_review FROM datastore.objects WHERE data_id=$1 AND upload_completed=TRUE AND deleted=FALSE`, dataID).Scan(
		&updatePassword,
		&underReview,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nex.ResultCodes.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nex.ResultCodes.DataStore.Unknown
	}

	if updatePassword != 0 && updatePassword != password {
		return nex.ResultCodes.DataStore.InvalidPassword
	}

	if underReview {
		return nex.ResultCodes.DataStore.UnderReviewing
	}

	_, err = database.Postgres.Exec(`UPDATE datastore.objects SET data_type=$1 WHERE data_id=$2`, dataType, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.ResultCodes.DataStore.Unknown
	}

	return 0
}
