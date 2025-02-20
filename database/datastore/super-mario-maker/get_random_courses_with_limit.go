package datastore_smm_db

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/lib/pq"
)

func GetRandomCoursesWithLimit(limit int) ([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, uint32) {
	courses := make([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, 0)

	rows, err := database.Postgres.Query(`
		SELECT
			object.data_id,
			object.owner,
			object.size,
			object.name,
			object.data_type,
			object.meta_binary,
			object.permission,
			object.permission_recipients,
			object.delete_permission,
			object.delete_permission_recipients,
			object.period,
			object.refer_data_id,
			object.flag,
			object.tags,
			object.creation_date,
			object.update_date,
			ranking.value
		FROM datastore.objects object
		JOIN datastore.object_custom_rankings ranking
		ON
			object.data_id = ranking.data_id AND
			object.upload_completed = TRUE AND
			object.deleted = FALSE AND
			object.under_review = FALSE AND
			ranking.application_id = 0
		ORDER BY RANDOM()
		LIMIT $1
	`, limit)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nex.ResultCodes.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.ResultCodes.DataStore.Unknown
	}

	defer rows.Close()

	for rows.Next() {
		course := datastore_super_mario_maker_types.NewDataStoreCustomRankingResult()
		course.Order = 0 // * Order is ALWAYS 0
		course.MetaInfo = datastore_types.NewDataStoreMetaInfo()
		course.MetaInfo.Permission = datastore_types.NewDataStorePermission()
		course.MetaInfo.DelPermission = datastore_types.NewDataStorePermission()
		course.MetaInfo.CreatedTime = nex.NewDateTime(0)
		course.MetaInfo.UpdatedTime = nex.NewDateTime(0)
		course.MetaInfo.ReferredTime = nex.NewDateTime(0)
		course.MetaInfo.ExpireTime = nex.NewDateTime(0x9C3f3E0000) // * 9999-12-31T00:00:00.000Z. This is what the real server sends
		course.MetaInfo.Ratings = make([]*datastore_types.DataStoreRatingInfoWithSlot, 0)

		var createdDate time.Time
		var updatedDate time.Time

		err := rows.Scan(
			&course.MetaInfo.DataID,
			&course.MetaInfo.OwnerID,
			&course.MetaInfo.Size,
			&course.MetaInfo.Name,
			&course.MetaInfo.DataType,
			&course.MetaInfo.MetaBinary,
			&course.MetaInfo.Permission.Permission,
			pq.Array(&course.MetaInfo.Permission.RecipientIDs),
			&course.MetaInfo.DelPermission.Permission,
			pq.Array(&course.MetaInfo.DelPermission.RecipientIDs),
			&course.MetaInfo.Period,
			&course.MetaInfo.ReferDataID,
			&course.MetaInfo.Flag,
			pq.Array(&course.MetaInfo.Tags),
			&createdDate,
			&updatedDate,
			&course.Score,
		)
		if err != nil {
			globals.Logger.Error(err.Error())
			continue
		}

		ratings, errCode := datastore_db.GetObjectRatingsWithSlotByDataID(course.MetaInfo.DataID)
		if errCode != 0 {
			return nil, errCode
		}

		course.MetaInfo.Ratings = ratings

		course.MetaInfo.CreatedTime.FromTimestamp(createdDate)
		course.MetaInfo.UpdatedTime.FromTimestamp(updatedDate)
		course.MetaInfo.ReferredTime.FromTimestamp(createdDate) // * This is what the real server does

		courses = append(courses, course)
	}

	return courses, 0
}
