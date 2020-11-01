package test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Wilder60/ShadowKeep/test/mocks"
)

func TestAdapterGetEventsInRangeSuccess(T *testing.T) {

	resp := httptest.NewRecorder()
	c, r := gin.CreateTestContext(resp)

	c.Set("UserID", val.UserID)
	c.Set("db", mocks.NewStorageMock())
	r.ServeHTTP(resp, c.Request)

}

func TestAdapterGetEventsInRangeFailure(T *testing.T) {
	resp := httptest.NewRecorder()
	c, r := gin.CreateTestContext(resp)
	c.Set("UserID", "")
	c.Set("db", nil)
}
