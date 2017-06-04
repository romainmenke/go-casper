package casper

import (
	"errors"
	"net/http"
	"strings"
)

func Push(w *responseWriter) error {
	c := CasperFromContext(w.ctx)
	if c == nil {
		return errors.New("casper was not defined")
	}

	// Remove casper cookie header if it's already exists.
	if cookies, ok := w.Header()["Set-Cookie"]; ok && len(cookies) != 0 {
		w.Header().Del("Set-Cookie")
		for _, cookieStr := range cookies {
			if strings.Contains(cookieStr, defaultCookieName+"=") {
				continue
			}
			w.Header().Add("Set-Cookie", cookieStr)
		}
	}

	hashValues := HashValuesFromContext(w.ctx)

	var somethingPushed bool

	for target, options := range w.targets {
		h := c.hash([]byte(target))

		// Check the content is already pushed or not.
		if search(hashValues, h) {
			continue
		}

		if !c.skipPush {
			if err := w.Push(target, options); err != nil {
				return err
			}
			somethingPushed = true
		}

		hashValues = append(hashValues, h)
	}

	if !somethingPushed {
		return nil
	}

	cookie, err := c.generateCookie(hashValues)
	if err != nil {
		return err
	}
	http.SetCookie(w, cookie)

	return nil
}
