package main

import (
	"github.com/gofiber/fiber/v2"
)

func httpRecord(
	c *fiber.Ctx,
	conf *config,
	refererParser *refererParser,
	geoParser *geoParser,
	userAgentParser *userAgentParser,
	projectsManager *projects,
	st *storage,
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
	record.setQueryParameters(
		c.Query(record_query_url),
		c.Query(record_query_canonical),
		c.Query(record_query_title),
		c.Query(record_query_lang),
		c.Query(record_query_entity_id),
		c.Query(record_query_entity_module),
		c.Query(record_query_entity_taxonomy_id),
	)

	userAgent := c.Get(fiber.HeaderUserAgent)

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
			return httpErrorResponse(
				c,
				fiber.StatusBadRequest,
				fiber.StatusBadRequest,
			)
		}

		record.setPostRequest(&postData, refererParser, geoParser)

		if record.Mode == recordModeEventAPI {
			// check api key
			apiVerifyError := record.verify(projectsManager, c.Get(api_key_header))
			if apiVerifyError != nil {
				return httpErrorResponse(
					c,
					apiVerifyError.msg,
					apiVerifyError.code,
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

			// updates
			userAgent = postData.API.ClientUserAgent
			ip = record.IP
			cid := clientIDFromOther([]string{ip.String(), userAgent})

			// apply updates
			record.UserAgentResult = userAgentParser.parse(userAgent)
			record.GeoResult = geoParser.newResultFromIP(ip)
			record.CID = cid
		}

		if record.isClientError() {
			defer prometheusClientErrors.Inc()
			defer conf.getLogger().
				Error().
				Str("type", error_type_client).
				Str("ip", ip.String()).
				Str("err", postData.ClientErrorObject).
				Msg(postData.ClientErrorMessage)
			return c.JSON(1)
		}

		getVerifyError := record.verify(projectsManager, "")
		if getVerifyError != nil {
			return httpErrorResponse(
				c,
				getVerifyError.msg,
				getVerifyError.code,
			)
		}

		go func() {
			st.addRecord(record.finalize())
		}()

		return c.JSON(1)

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

			st.addRecord(record.finalize())
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
