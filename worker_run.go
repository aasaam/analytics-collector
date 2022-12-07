package main

import (
	"context"
	"encoding/json"
	"time"
)

type workerRunResult struct {
	e            error
	errorState   string
	timeTaken    float64
	records      int64
	clientErrors int64
}

func workerRun(
	clickhouseConfig *clickhouseConfig,
	conf *config,
	redisConnection string,
) *workerRunResult {

	r := workerRunResult{
		e:            nil,
		timeTaken:    0,
		records:      0,
		clientErrors: 0,
	}

	redisClient, redisClientErr := redisGetClient(redisConnection)
	if redisClientErr != nil {
		r.e = redisClientErr
		r.errorState = "redisClientErr"
		return &r
	}

	recordSize := redisClient.LLen(context.Background(), redisKeyRecords).Val()

	if recordSize == 0 {
		return &r
	}

	clickhouseConn, clickhouseCtx, clickhouseInitErr := clickhouseGetConnection(clickhouseConfig)

	if clickhouseInitErr != nil {
		r.e = clickhouseInitErr
		r.errorState = "clickhouseInitErr"
		return &r
	}

	var insertRecords int64 = 0
	var insertClientErrors int64 = 0

	recordsBatch, recordsBatchErr := clickhouseConn.PrepareBatch(
		clickhouseCtx, clickhouseSQLInsertRecords,
	)

	clientErrorsBatch, clientErrorsBatchErr := clickhouseConn.PrepareBatch(
		clickhouseCtx, clickhouseSQLInsertClientErrors,
	)

	if recordsBatchErr != nil {
		r.e = recordsBatchErr
		r.errorState = "recordsBatchErr"
		return &r
	}

	if clientErrorsBatchErr != nil {
		r.e = clientErrorsBatchErr
		r.errorState = "clientErrorsBatchErr"
		return &r
	}

	startTime := time.Now().UnixMilli()

	for {
		item := redisClient.LPop(context.Background(), redisKeyRecords).Val()
		if item == "" {
			conf.getLogger().Debug().Msg("no item found")
			return &r
		}

		var rec record
		err := json.Unmarshal([]byte(item), &rec)

		if err != nil {
			redisClient.RPush(context.Background(), redisKeyRecords, item)
			conf.getLogger().Error().Str("err", err.Error())
			continue
		}

		if rec.isClientError() {
			insertClientErr := clickhouseInsertClientErrBatch(clientErrorsBatch, rec)
			if insertClientErr != nil {
				redisClient.RPush(context.Background(), redisKeyRecords, item)
				r.e = insertClientErr
				r.errorState = "insertClientErr"
				return &r
			}
			insertClientErrors++
		} else {
			if rec.EventCount > 0 {
				for i := 0; i < rec.EventCount; i++ {

					ECategory := rec.Events[i].ECategory
					EAction := rec.Events[i].EAction
					ELabel := rec.Events[i].ELabel
					EIdent := rec.Events[i].EIdent
					EValue := rec.Events[i].EValue

					insertRecordErr := clickhouseInsertRecordBatch(recordsBatch, rec, ECategory, EAction, ELabel, EIdent, EValue)
					if insertRecordErr != nil {
						redisClient.RPush(context.Background(), redisKeyRecords, item)
						r.e = insertRecordErr
						r.errorState = "insertRecordErr"
						return &r
					}
					insertRecords += 1
				}
			} else {
				insertRecordErr := clickhouseInsertRecordBatch(recordsBatch, rec, "", "", "", "", 0)
				if insertRecordErr != nil {
					redisClient.RPush(context.Background(), redisKeyRecords, item)
					r.e = insertRecordErr
					r.errorState = "insertRecordErr"
					return &r
				}
				insertRecords += 1
			}

		}

		if redisClient.LLen(context.Background(), redisKeyRecords).Val() == 0 {
			break
		}
	}

	if insertRecords > 0 {
		recordsBatchSendErr := recordsBatch.Send()
		if recordsBatchSendErr != nil {
			r.e = recordsBatchSendErr
			r.errorState = "recordsBatchSendErr"
		}
	}
	if insertClientErrors > 0 {
		clientErrorsBatchSendErr := clientErrorsBatch.Send()
		if clientErrorsBatchSendErr != nil {
			r.e = clientErrorsBatchSendErr
			r.errorState = "insertClientErr"
		}
	}

	r.clientErrors = insertClientErrors
	r.records = insertRecords
	r.timeTaken = toFixed(float64(time.Now().UnixMilli()-startTime)/1000, 3)

	return &r

}
