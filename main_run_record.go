package main

import (
	"github.com/urfave/cli/v2"
)

func mainRecordRoutine(
	c *cli.Context,
	conf *config,
	redisClient *redisClient,
) {

	recordsCount := redisClient.countRecords()
	if recordsCount == 0 {
		conf.getLogger().
			Debug().
			Msg("queue is empty")
		return
	}

	// clickhouseConn, clickhouseCtx, clickhouseConnErr := clickhouseGetConnection(
	// 	c.String("clickhouse-servers"),
	// 	c.String("clickhouse-database"),
	// 	c.String("clickhouse-username"),
	// 	c.String("clickhouse-password"),
	// 	c.Int("clickhouse-max-execution-time"),
	// 	c.Int("clickhouse-dial-timeout"),
	// 	c.Bool("test-mode"),
	// 	c.Bool("clickhouse-compression-lz4"),
	// 	c.Int("clickhouse-max-idle-conns"),
	// 	c.Int("clickhouse-max-open-conns"),
	// 	c.Int("clickhouse-conn-max-lifetime"),
	// 	c.Int("clickhouse-max-block-size"),
	// 	nil,
	// 	nil,
	// )

	// if clickhouseConnErr != nil {
	// 	conf.getLogger().
	// 		Error().
	// 		Str("type", errorTypeApp).
	// 		Str("on", "clickhouse-connection").
	// 		Str("error", clickhouseConnErr.Error()).
	// 		Send()
	// 	return
	// }

	// // startTime := time.Now().UnixMilli()
	// // inserts := 0

	// records, recordsErr := redisClient.popRecord()
	// if recordsErr != nil {
	// 	conf.getLogger().
	// 		Error().
	// 		Str("type", errorTypeApp).
	// 		Str("on", "clickhouse-connection").
	// 		Str("error", recordsErr.Error()).
	// 		Send()
	// 	return
	// }

	// /**
	//  * Records
	//  */
	// go func() {
	// 	for {
	// 		func() {
	// 			startTime := time.Now().UnixMilli()
	// 			inserts := 0

	//
	// 			//
	// 			// records
	// 			//
	// 			if storage.recordCount > 0 {
	// 				records := storage.getRecords()

	// 				recordsBatch, recordsBatchErr := clickhouseConn.PrepareBatch(
	// 					clickhouseCtx, clickhouseInsertRecords,
	// 				)
	// 				if recordsBatchErr != nil {
	// 					conf.getLogger().
	// 						Error().
	// 						Str("type", errorTypeApp).
	// 						Str("on", "clickhouse-connection").
	// 						Str("error", recordsBatchErr.Error()).
	// 						Send()
	// 					time.Sleep(clickhouseInterval)
	// 					return
	// 				}

	// 				for _, recordByte := range records {
	// 					recordByteReader := bytes.NewReader(recordByte)

	// 					var rec record
	// 					recordDecodeErr := gob.NewDecoder(recordByteReader).Decode(&rec)
	// 					if recordDecodeErr != nil {
	// 						conf.getLogger().
	// 							Error().
	// 							Str("type", errorTypeApp).
	// 							Str("on", "record-decode").
	// 							Str("error", recordDecodeErr.Error()).
	// 							Send()
	// 						continue
	// 					}

	// 					if rec.EventCount > 0 {
	// 						for i := 0; i < rec.EventCount; i++ {
	// 							ECategory := rec.Events[i].ECategory
	// 							EAction := rec.Events[i].EAction
	// 							ELabel := rec.Events[i].ELabel
	// 							EIdent := rec.Events[i].EIdent
	// 							EValue := rec.Events[i].EValue
	// 							insertErr := insertRecordBatch(recordsBatch, rec, ECategory, EAction, ELabel, EIdent, EValue)
	// 							if insertErr != nil {
	// 								conf.getLogger().
	// 									Error().
	// 									Str("type", errorTypeApp).
	// 									Str("on", "record-insert").
	// 									Str("error", insertErr.Error()).
	// 									Send()
	// 							}
	// 							inserts += 1
	// 						}
	// 					} else {
	// 						insertErr := insertRecordBatch(recordsBatch, rec, "", "", "", "", 0)
	// 						if insertErr != nil {
	// 							conf.getLogger().
	// 								Error().
	// 								Str("type", errorTypeApp).
	// 								Str("on", "record-insert").
	// 								Str("error", insertErr.Error()).
	// 								Send()
	// 						}
	// 						inserts += 1
	// 					}
	// 				}

	// 				recordsBatchSendErr := recordsBatch.Send()
	// 				if recordsBatchSendErr != nil {
	// 					conf.getLogger().
	// 						Error().
	// 						Str("type", errorTypeApp).
	// 						Str("on", "record-batch-send").
	// 						Str("error", recordsBatchSendErr.Error()).
	// 						Send()
	// 				}

	// 				storage.cleanRecords()
	// 			}

	// 			//
	// 			// client errors
	// 			//
	// 			if storage.clientErrorCount > 0 {
	// 				clientErrors := storage.getClientErrors()

	// 				clientErrorsBatch, clientErrorsBatchErr := clickhouseConn.PrepareBatch(
	// 					clickhouseCtx, clickhouseInsertClientErrors,
	// 				)
	// 				if clientErrorsBatchErr != nil {
	// 					conf.getLogger().
	// 						Error().
	// 						Str("type", errorTypeApp).
	// 						Str("on", "clickhouse-connection").
	// 						Str("error", clientErrorsBatchErr.Error()).
	// 						Send()
	// 					time.Sleep(clickhouseInterval)
	// 					return
	// 				}

	// 				for _, clientErrorByte := range clientErrors {
	// 					clientErrorByteReader := bytes.NewReader(clientErrorByte)

	// 					var ce record
	// 					clientErrorDecodeErr := gob.NewDecoder(clientErrorByteReader).Decode(&ce)
	// 					if clientErrorDecodeErr != nil {
	// 						conf.getLogger().
	// 							Error().
	// 							Str("type", errorTypeApp).
	// 							Str("on", "client-error-decode").
	// 							Str("error", clientErrorDecodeErr.Error()).
	// 							Send()
	// 						continue
	// 					}

	// 					insertErr := insertClientErrBatch(clientErrorsBatch, ce)

	// 					if insertErr != nil {
	// 						conf.getLogger().
	// 							Error().
	// 							Str("type", errorTypeApp).
	// 							Str("on", "client-error-insert").
	// 							Str("error", insertErr.Error()).
	// 							Send()
	// 					}
	// 					inserts += 1
	// 				}

	// 				clientErrorsBatchSendErr := clientErrorsBatch.Send()
	// 				if clientErrorsBatchSendErr != nil {
	// 					conf.getLogger().
	// 						Error().
	// 						Str("type", errorTypeApp).
	// 						Str("on", "client-error-batch-send").
	// 						Str("error", clientErrorsBatchSendErr.Error()).
	// 						Send()
	// 				}

	// 				storage.cleanRecords()
	// 			}

	// 			if inserts > 0 {
	// 				endTime := time.Now().UnixMilli()
	// 				inSeconds := (float64(endTime) - float64(startTime)) / 1000
	// 				conf.getLogger().
	// 					Debug().
	// 					Msg(fmt.Sprintf("Insert %d item(s) in %.2f seconds(s)", inserts, inSeconds))
	// 			}
	// 		}()
	// 		time.Sleep(clickhouseInterval)
	// 	}
	// }()

	// conn, connErr := pgx.Connect(context.Background(), c.String("postgis-uri"))
	// if connErr != nil {
	// 	return connErr
	// }

	// defer conn.Close(context.Background())

	// geoParser, geoParserErr := newGeoParser(conn, c.String("mmdb-city-path"), c.String("mmdb-asn-path"))
	// if geoParserErr != nil {
	// 	return geoParserErr
	// }

	// refererParser := newRefererParser()
	// userAgentParser := newUserAgentParser()

	// app := newHTTPServer(conf, geoParser, refererParser, userAgentParser, projectsManager, redisClient)
	// return app.Listen(c.String("listen"))
}
