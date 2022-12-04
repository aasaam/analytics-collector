package main

import (
	"encoding/json"
	"net"

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

	record, recordErr := newRecord(c.Query(recordQueryMode), c.Query(recordQueryPublicInstanceID))

	if recordErr != nil {
		return httpErrorResponse(
			c,
			errorInvalidModeOrProjectPublicID,
		)
	}

	ip := getClientIP(c)
	userAgent := c.Get(fiber.HeaderUserAgent)

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

		if record.Mode == recordModeClientError || record.Mode == recordModeClientErrorLegacy {
			record.CID = clientIDNoneSTD([]string{ip.String(), userAgent}, clientIDTypeOther)
		}
	}

	// recordModePageViewJavaScript
	// recordModePageViewAMP
	// recordModeEvent*
	if c.Method() == fiber.MethodPost {
		var postData postRequest
		if postDataErr := json.Unmarshal(c.Body(), &postData); postDataErr != nil {
			return httpErrorResponse(
				c,
				errorBadPOSTBody,
			)
		}

		// changes if in api mode
		if record.Mode == recordModeEventAPI {
			// set api parameters
			apiErr := record.setAPI(&postData)
			if apiErr != nil {
				return httpErrorResponse(
					c,
					*apiErr,
				)
			}
		} else if record.Mode == recordModeClientError { // on client error

			go func() {
				record.ClientErrorMessage = postData.ClientErrorMessage
				record.ClientErrorObject = postData.ClientErrorObject

				finalizeByte, finalizeErr := record.finalize()
				if finalizeErr == nil {
					storage.addClientError(finalizeByte)
					return
				}

				conf.getLogger().
					Error().
					Str("type", errorTypeApp).
					Str("error", finalizeErr.Error()).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Send()
			}()
			defer promMetricClientErrors.Inc()
			return c.JSON(1)
		}

		if record.Mode == recordModeEventAPI && postData.API != nil {
			// updates
			userAgent = postData.API.ClientUserAgent
			ip = net.IP(postData.API.ClientIP)

			// apply updates
			record.UserAgentResult = userAgentParser.parse(userAgent)
			record.GeoResult = geoParser.newResultFromIP(ip)
		}

		record.setPostRequest(&postData, refererParser, geoParser)

		// changes if in api mode
		if record.Mode == recordModeEventAPI {
			// check api key
			apiVerifyError := record.verify(projectsManager, postData.API.PrivateInstanceKey)
			if apiVerifyError != nil {
				return httpErrorResponse(
					c,
					*apiVerifyError,
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
					Str("type", errorTypeApp).
					Str("error", finalizeErr.Error()).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Send()
			}()

			return c.JSON(true)
		}

		postVerifyError := record.verify(projectsManager, "")
		if postVerifyError != nil {
			return httpErrorResponse(
				c,
				*postVerifyError,
			)
		}

		go func() {
			finalizeByte, finalizeErr := record.finalize()
			if finalizeErr == nil {
				storage.addRecord(finalizeByte)
			} else {
				conf.getLogger().
					Error().
					Str("type", errorTypeApp).
					Str("error", finalizeErr.Error()).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Send()
			}
		}()

		return c.JSON(1)

		// recordModePageViewImageLegacy
		// recordModePageViewImageNoScript
		// recordModePageViewAMPImage
	} else if c.Method() == fiber.MethodGet && record.isImage() {
		if record.Mode == recordModeClientErrorLegacy {

			go func() {
				record.ClientErrorMessage = "legacy"
				record.ClientErrorObject = c.Get("err")

				finalizeByte, finalizeErr := record.finalize()
				if finalizeErr == nil {
					storage.addClientError(finalizeByte)
					return
				}

				conf.getLogger().
					Error().
					Str("type", errorTypeApp).
					Str("error", finalizeErr.Error()).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Send()
			}()

		} else {
			record.CID = clientIDNoneSTD([]string{ip.String(), userAgent}, clientIDTypeOther)

			if record.pURL == nil && record.Mode == recordModePageViewImageNoScript {
				imgReferer := getURL(c.Get(fiber.HeaderReferer))
				if imgReferer != nil {
					record.PURL = imgReferer.String()
					record.pURL = imgReferer
				}
			}

			record.setReferer(refererParser, getURL(c.Query(recordQueryRefererURL)))

			getVerifyError := record.verify(projectsManager, "")
			if getVerifyError != nil {
				return httpErrorResponse(
					c,
					*getVerifyError,
				)
			}

			go func() {
				finalizeByte, finalizeErr := record.finalize()
				if finalizeErr == nil {
					storage.addRecord(finalizeByte)
				} else {
					conf.getLogger().
						Error().
						Str("type", errorTypeApp).
						Str("error", finalizeErr.Error()).
						Str("ip", ip.String()).
						Str("method", c.Method()).
						Str("path", c.Path()).
						Send()
				}
			}()
		}

		c.Set(fiber.HeaderContentType, mimetypeGIF)
		return c.Send(singleGifImage)
	}

	return httpErrorResponse(
		c,
		errorRecordNotValid,
	)
}
