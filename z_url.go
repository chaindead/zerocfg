package zerocfg

import (
	"fmt"
	"net/url"
)

type urlValue struct {
	val **url.URL
}

func newURLValue(val *url.URL, p **url.URL) Value {
	*p = val
	return &urlValue{val: p}
}

func (u *urlValue) Set(s string) error {
	if s == "" {
		*(u.val) = nil
		return nil
	}
	parsedURL, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("parsing URL %q: %w", s, err)
	}
	*(u.val) = parsedURL
	return nil
}

func (u *urlValue) Type() string {
	return "url"
}

func (u *urlValue) String() string {
	if u.val == nil || *(u.val) == nil {
		return ""
	}
	return (*(u.val)).String()
}

func URL(name string, defValStr string, desc string, opts ...OptNode) **url.URL {
	var initialURL *url.URL
	if defValStr != "" {
		parsedURL, err := url.Parse(defValStr)
		if err != nil {
			panic(fmt.Sprintf("invalid default URL value %q for %q: %v", defValStr, name, err))
		}
		initialURL = parsedURL
	}

	ptr := new(*url.URL)
	*ptr = initialURL

	return Any(name, initialURL, desc, func(val *url.URL, p **url.URL) Value {

		return newURLValue(val, p)
	}, opts...)
}
