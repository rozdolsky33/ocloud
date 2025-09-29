package loadbalancer

import (
	x509std "crypto/x509"
	pemenc "encoding/pem"
	"time"
)

// parseCertNotAfter attempts to parse the first certificate in a PEM bundle and returns NotAfter
func parseCertNotAfter(pemData string) (time.Time, bool) {
	data := []byte(pemData)
	for {
		var block *pemenc.Block
		block, data = pemenc.Decode(data)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			c, err := x509std.ParseCertificate(block.Bytes)
			if err == nil {
				return c.NotAfter, true
			}
		}
	}
	return time.Time{}, false
}
