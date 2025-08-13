package netx

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetUserIP(c *fiber.Ctx) string {
	if ip := c.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		if idx := strings.Index(ip, ","); idx != -1 {
			return strings.TrimSpace(ip[:idx])
		}
		return strings.TrimSpace(ip)
	}
	if ip := c.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}
	ip := c.IP()
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return strings.TrimSpace(ip)
}
