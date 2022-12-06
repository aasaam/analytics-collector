package main

import (
	"github.com/go-redis/redis/v9"
)

type workerRunResult struct {
	e            error
	errorState   string
	timeTakenMS  int64
	records      int64
	clientErrors int64
}

func workerRun(
	clickhouseConfig *clickhouseConfig,
	conf *config,
	redisClient *redis.Client,
) *workerRunResult {

	return nil
	// startTime := time.Now().UnixMilli()
	// r := workerRunResult{
	// 	e:            nil,
	// 	timeTakenMS:  0,
	// 	records:      0,
	// 	clientErrors: 0,
	// }

	// var insertRecords int64 = 0
	// var insertClientErrors int64 = 0

	// storage.Lock()
	// defer storage.Unlock()

	// // no storage data check
	// if storage.recordCount == 0 && storage.clientErrorCount == 0 {
	// 	return &r
	// }

	// clickhouseConn, clickhouseCtx, clickhouseInitErr := clickhouseGetConnection(clickhouseConfig)

	// if clickhouseInitErr != nil {
	// 	r.e = clickhouseInitErr
	// 	r.errorState = "clickhouseInitErr"
	// 	return &r
	// }

	// if storage.recordCount > 0 {

	// 	recordsBatch, recordsBatchErr := clickhouseConn.PrepareBatch(
	// 		clickhouseCtx, clickhouseSQLInsertRecords,
	// 	)
	// 	if recordsBatchErr != nil {
	// 		r.e = recordsBatchErr
	// 		r.errorState = "recordsBatchErr"
	// 		return &r
	// 	}

	// 	for _, recordByte := range storage.getRecords() {
	// 		recordByteReader := bytes.NewReader(recordByte)

	// 		var rec record
	// 		recordDecodeErr := gob.NewDecoder(recordByteReader).Decode(&rec)
	// 		if recordDecodeErr != nil {
	// 			r.e = recordDecodeErr
	// 			r.errorState = "recordDecodeErr"
	// 			return &r
	// 		}

	// 		if !validatePublicInstanceIDRegex(rec.PublicInstanceID) {
	// 			defer promMetricInvalidProcessData.Inc()
	// 			continue
	// 		}

	// 		if rec.EventCount > 0 {
	// 			for i := 0; i < rec.EventCount; i++ {
	// 				ECategory := rec.Events[i].ECategory
	// 				EAction := rec.Events[i].EAction
	// 				ELabel := rec.Events[i].ELabel
	// 				EIdent := rec.Events[i].EIdent
	// 				EValue := rec.Events[i].EValue
	// 				insertRecordErr := clickhouseInsertRecordBatch(recordsBatch, rec, ECategory, EAction, ELabel, EIdent, EValue)
	// 				if insertRecordErr != nil {
	// 					r.e = insertRecordErr
	// 					r.errorState = "insertRecordErr"
	// 					return &r
	// 				}
	// 				insertRecords += 1
	// 			}
	// 		} else {
	// 			insertRecordErr := clickhouseInsertRecordBatch(recordsBatch, rec, "", "", "", "", 0)
	// 			if insertRecordErr != nil {
	// 				r.e = insertRecordErr
	// 				r.errorState = "insertRecordErr"
	// 				return &r
	// 			}
	// 			insertRecords += 1
	// 		}

	// 	}

	// 	recordsBatchSendErr := recordsBatch.Send()
	// 	if recordsBatchSendErr != nil {
	// 		r.e = recordsBatchSendErr
	// 		r.errorState = "recordsBatchSendErr"
	// 		return &r
	// 	}

	// 	storage.cleanRecords()
	// 	r.records = insertRecords
	// }

	// if storage.clientErrorCount > 0 {

	// 	clientErrorsBatch, clientErrorsBatchErr := clickhouseConn.PrepareBatch(
	// 		clickhouseCtx, clickhouseSQLInsertClientErrors,
	// 	)

	// 	if clientErrorsBatchErr != nil {
	// 		r.e = clientErrorsBatchErr
	// 		r.errorState = "clientErrorsBatchErr"
	// 		return &r
	// 	}

	// 	for _, clientErrorByte := range storage.getClientErrors() {
	// 		clientErrorByteReader := bytes.NewReader(clientErrorByte)
	// 		var ce record
	// 		clientErrorDecodeErr := gob.NewDecoder(clientErrorByteReader).Decode(&ce)

	// 		if clientErrorDecodeErr != nil {
	// 			r.e = clientErrorDecodeErr
	// 			r.errorState = "clientErrorDecodeErr"
	// 			return &r
	// 		}

	// 		if !validatePublicInstanceIDRegex(ce.PublicInstanceID) {
	// 			defer promMetricInvalidProcessData.Inc()
	// 			continue
	// 		}

	// 		insertClientErr := clickhouseInsertClientErrBatch(clientErrorsBatch, ce)
	// 		if insertClientErr != nil {
	// 			r.e = insertClientErr
	// 			r.errorState = "insertClientErr"
	// 			return &r
	// 		}

	// 		insertClientErrors += 1
	// 	}

	// 	clientErrorsBatchSendErr := clientErrorsBatch.Send()
	// 	if clientErrorsBatchSendErr != nil {
	// 		r.e = clientErrorsBatchSendErr
	// 		r.errorState = "insertClientErr"
	// 		return &r
	// 	}
	// 	storage.cleanRecords()
	// 	r.clientErrors = insertClientErrors
	// }

	// r.timeTakenMS = time.Now().UnixMilli() - startTime
	// return &r
}
