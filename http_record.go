package main

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
)

func responseImage(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, mimetypeGIF)
	return c.Send(singleGifImage)
}

func httpRecord(
	c *fiber.Ctx,
	conf *config,
	refererParser *refererParser,
	geoParser *geoParser,
	userAgentParser *userAgentParser,
	projectsManager *projects,
	redisClient *redis.Client,
) error {
	// no cache at all
	noCache(c)

	record, recordErr := newRecord(c.Query(recordQueryMode), c.Query(recordQueryPublicInstanceID))

	ip := getClientIP(c)
	userAgent := c.Get(fiber.HeaderUserAgent)

	if recordErr != nil {
		blockErr := errorInvalidModeOrProjectPublicID
		blockErr.debug = recordErr.Error()

		promMetricInvalidRequestData.WithLabelValues(blockErr.msg).Inc()

		conf.getLogger().
			Warn().
			Str("part", "init").
			Str("on", blockErr.msg).
			Str("error", recordErr.Error()).
			Str("ip", ip.String()).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Send()

		return httpErrorResponse(
			c,
			blockErr,
		)
	}

	record.setQueryParameters(
		c.Query(recordQueryURL),
		c.Query(recordQueryCanonicalURL),
		c.Query(recordQueryTitle),
		c.Query(recordQueryLang),
		c.Query(recordQueryEntityID),
		c.Query(recordQueryEntityModule),
		c.Query(recordQueryEntityTaxonomyID),
	)

	// in not api mode ip must get from request
	if record.Mode != recordModeEventAPI {
		record.IP = ip
		record.GeoResult = geoParser.newResultFromIP(ip)
		record.UserAgentResult = userAgentParser.parse(userAgent)
		record.CID = clientIDNoneSTD([]string{ip.String(), userAgent}, clientIDTypeOther)
	}

	var postData postRequest

	// recordModePageViewJavaScript
	// recordModePageViewAMP
	// recordModeEvent*
	if c.Method() == fiber.MethodPost {
		if postDataErr := json.Unmarshal(c.Body(), &postData); postDataErr != nil {
			blockErr := errorBadPOSTBody

			promMetricInvalidRequestData.WithLabelValues(blockErr.msg).Inc()

			conf.getLogger().
				Info().
				Str("part", "parse_body").
				Str("on", blockErr.msg).
				Str("error", postDataErr.Error()).
				Str("ip", ip.String()).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Send()

			return httpErrorResponse(
				c,
				errorBadPOSTBody,
			)
		}

		if postData.Page != nil {
			record.setReferer(refererParser, getURL(postData.Page.RefererURL))
		}

		if record.Mode != recordModeEventAPI {
			record.setPostRequest(&postData, refererParser, geoParser)
		}

		if record.Mode == recordModeClientError {
			record.ClientErrorMessage = postData.ClientErrorMessage
			record.ClientErrorObject = postData.ClientErrorObject
		}

	} else if c.Method() == fiber.MethodGet {

		// try detect page url via referer header of image
		if record.pURL == nil && record.Mode == recordModePageViewImageNoScript {
			imgReferer := getURL(c.Get(fiber.HeaderReferer))
			if imgReferer != nil {
				record.PURL = imgReferer.String()
				record.pURL = imgReferer
			}
		}

		record.setReferer(refererParser, getURL(c.Query(recordQueryRefererURL)))

		if record.Mode == recordModeClientErrorLegacy {
			legacyErrorObject := c.Query(recordQueryErrorVeryLegacy)
			if legacyErrorObject != "" {
				record.ClientErrorMessage = "veryLegacy"
				record.ClientErrorObject = legacyErrorObject
			}
		}
	}

	// changes if in api mode
	if record.Mode == recordModeEventAPI {

		// set api parameters
		apiErr := record.setAPI(projectsManager, userAgentParser, geoParser, &postData)
		if apiErr != nil {

			promMetricInvalidRequestData.WithLabelValues(apiErr.msg).Inc()

			conf.getLogger().
				Warn().
				Str("part", "api").
				Str("on", apiErr.msg).
				Str("error", apiErr.msg).
				Str("ip", ip.String()).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Send()

			return httpErrorResponse(
				c,
				*apiErr,
			)
		}

	} else {

		verifyErr := record.verify(projectsManager, "")
		if verifyErr != nil {

			promMetricInvalidRequestData.WithLabelValues(verifyErr.msg).Inc()

			conf.getLogger().
				Warn().
				Str("part", "verify").
				Str("on", verifyErr.msg).
				Str("error", verifyErr.msg).
				Str("ip", ip.String()).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Send()

			return httpErrorResponse(
				c,
				*verifyErr,
			)
		}
	}

	recordBytes, recordBytesErr := record.finalize()

	if recordBytesErr != nil {

		promMetricInvalidRequestData.WithLabelValues(recordBytesErr.msg).Inc()

		conf.getLogger().
			Error().
			Str("part", "finalize").
			Str("on", recordBytesErr.msg).
			Str("error", recordBytesErr.debug).
			Str("ip", ip.String()).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Send()

		return httpErrorResponse(
			c,
			*recordBytesErr,
		)
	}

	_, redErr := redisClient.LPush(context.Background(), redisKeyRecords, recordBytes).Result()

	if redErr != nil {
		blockErr := errorInternalDependencyFailed
		blockErr.debug = redErr.Error()

		promMetricInvalidRequestData.WithLabelValues(blockErr.msg).Inc()

		conf.getLogger().
			Error().
			Str("part", "finalize").
			Str("on", blockErr.msg).
			Str("error", blockErr.debug).
			Str("ip", ip.String()).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Send()

		return httpErrorResponse(
			c,
			blockErr,
		)
	}

	if record.isImage() {
		return responseImage(c)
	}

	return c.JSON(1)
}
