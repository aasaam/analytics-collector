package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func httpRecord(
	c *fiber.Ctx,
	conf *config,
	refererParser *refererParser,
	geoParser *geoParser,
	userAgentParser *userAgentParser,
	projectsManager *projects,
	storage *storage,
) error {
	// no cache at all
	noCache(c)

	record, recordErr := newRecord(c.Query(record_query_mode), c.Query(record_query_public_instace_id))

	if recordErr != nil {
		return httpErrorResponse(
			c,
			error_invalid_mode_or_project_public_id.msg,
			error_invalid_mode_or_project_public_id.code,
		)
	}

	ip := getClientIP(c)
	userAgent := c.Get(fiber.HeaderUserAgent)

	record.setQueryParameters(
		c.Query(record_query_url),
		c.Query(record_query_canonical),
		c.Query(record_query_title),
		c.Query(record_query_lang),
		c.Query(record_query_entity_id),
		c.Query(record_query_entity_module),
		c.Query(record_query_entity_taxonomy_id),
	)

	// in not api mode ip must get from request
	if record.Mode != recordModeEventAPI {
		record.IP = ip
		record.GeoResult = geoParser.newResultFromIP(ip)
		record.UserAgentResult = userAgentParser.parse(userAgent)
	}

	// recordModePageViewJavaScript
	// recordModePageViewAMP
	// recordModeEvent*
	if c.Method() == fiber.MethodPost {
		var postData postRequest
		if postDataErr := c.BodyParser(&postData); postDataErr != nil {
			fmt.Println(postDataErr)
			return httpErrorResponse(
				c,
				fiber.StatusBadRequest,
				fiber.StatusBadRequest,
			)
		}

		if record.Mode == recordModeEventAPI && postData.API != nil {
			// updates
			userAgent = postData.API.ClientUserAgent
			ip = record.IP
			cid := clientIDFromOther([]string{ip.String(), userAgent})

			// apply updates
			record.UserAgentResult = userAgentParser.parse(userAgent)
			record.GeoResult = geoParser.newResultFromIP(ip)
			record.CID = cid
		}

		setPostRequestErr := record.setPostRequest(&postData, refererParser, geoParser)

		if setPostRequestErr != nil {
			return httpErrorResponse(
				c,
				setPostRequestErr.msg,
				setPostRequestErr.code,
			)
		}

		if record.Mode == recordModeClientError { // on client error

			go func() {
				cid := clientIDFromOther([]string{ip.String(), userAgent})
				record.CID = cid

				conf.getLogger().
					Debug().
					Str("type", error_type_client).
					Str("error", postData.ClientErrorObject).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Msg(postData.ClientErrorMessage)

				finalizeByte, finalizeErr := record.finalize()
				if finalizeErr == nil {
					storage.addClientError(finalizeByte)
					return
				}

				conf.getLogger().
					Error().
					Str("type", error_type_app).
					Str("error", finalizeErr.Error()).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Send()
			}()

			defer conf.getLogger().
				Debug().
				Str("type", error_type_client).
				Str("ip", ip.String()).
				Str("err", postData.ClientErrorObject).
				Msg(postData.ClientErrorMessage)
			defer promMetricClientErrors.Inc()
			return c.JSON(1)
		}

		// changes if in api mode
		if record.Mode == recordModeEventAPI {
			if postData.API == nil {
				return httpErrorResponse(
					c,
					error_api_fields_missed.msg,
					error_api_fields_missed.code,
				)
			}

			// set api parameters
			apiErr := record.setAPI(postData.API)
			if apiErr != nil {
				return httpErrorResponse(
					c,
					apiErr.msg,
					apiErr.code,
				)
			}

			// check api key
			apiVerifyError := record.verify(projectsManager, postData.API.PrivateInstanceKey)
			if apiVerifyError != nil {
				return httpErrorResponse(
					c,
					apiVerifyError.msg,
					apiVerifyError.code,
				)
			}

			go func() {
				finalizeByte, finalizeErr := record.finalize()
				if finalizeErr == nil {
					storage.addRecord(finalizeByte)
					return
				}
				conf.getLogger().
					Error().
					Str("type", error_type_app).
					Str("error", finalizeErr.Error()).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Send()
			}()

			return c.JSON(true)

		} else {
			getVerifyError := record.verify(projectsManager, "")

			if getVerifyError != nil {
				return httpErrorResponse(
					c,
					getVerifyError.msg,
					getVerifyError.code,
				)
			}

			go func() {
				finalizeByte, finalizeErr := record.finalize()
				if finalizeErr == nil {
					storage.addRecord(finalizeByte)
				} else {
					conf.getLogger().
						Error().
						Str("type", error_type_app).
						Str("error", finalizeErr.Error()).
						Str("ip", ip.String()).
						Str("method", c.Method()).
						Str("path", c.Path()).
						Send()
				}
			}()

			return c.JSON(1)
		}

		// recordModePageViewImageLegacy
		// recordModePageViewImageNoScript
		// recordModePageViewAMPImage
	} else if c.Method() == fiber.MethodGet {
		cid := clientIDFromOther([]string{ip.String(), userAgent})
		record.CID = cid

		if record.PURL == nil && record.Mode == recordModePageViewImageNoScript {
			imgReferer := getURL(c.Get(fiber.HeaderReferer))
			if imgReferer != nil {
				record.PURL = imgReferer
			}
		}

		record.setReferer(refererParser, getURL(c.Query(record_query_referer)))

		getVerifyError := record.verify(projectsManager, "")
		if getVerifyError != nil {
			return httpErrorResponse(
				c,
				getVerifyError.msg,
				getVerifyError.code,
			)
		}

		go func() {
			finalizeByte, finalizeErr := record.finalize()
			if finalizeErr == nil {
				storage.addRecord(finalizeByte)
			} else {
				conf.getLogger().
					Error().
					Str("type", error_type_app).
					Str("error", finalizeErr.Error()).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Send()
			}
		}()

		// image response single gif
		if record.isImage() {
			c.Set(fiber.HeaderContentType, mimetype_gif)
			return c.Send(single_gif_image)
		}

		return c.JSON(1)
	}

	return httpErrorResponse(
		c,
		error_record_not_valid.msg,
		error_record_not_valid.code,
	)
}
