package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetObjectSizeDataID(dataID uint64) (uint32, uint32) {
	var size uint32

	err := database.Postgres.QueryRow(`SELECT size FROM datastore.objects WHERE data_id=$1`, dataID).Scan(&size)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nex.ResultCodes.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return 0, nex.ResultCodes.DataStore.Unknown
	}

	return size, 0
}
