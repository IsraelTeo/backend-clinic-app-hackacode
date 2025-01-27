package mapper

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ParseID(c echo.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("the id must not be less than 1 and must be numeric")
	}
	return uint(id), nil
}
