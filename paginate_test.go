package fiberpaginate

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// go test -run Test_PaginateWithQueries
func Test_PaginateWithQueries(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?page=2&limit=20", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer func() {
		closeErr := body.Close()
		require.Equal(t, nil, closeErr)
	}()

	bodyBytes, err := io.ReadAll(body)
	require.Equal(t, nil, err)

	var respBody Response
	require.Equal(t, nil, json.Unmarshal(bodyBytes, &respBody))

	require.Equal(t, 2, respBody.Page)
	require.Equal(t, 20, respBody.Limit)
}

// go test -run Test_PaginateCheckDefaultsWhenNoQueries
func Test_PaginateCheckDefaultsWhenNoQueries(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer func() {
		closeErr := body.Close()
		require.Equal(t, nil, closeErr)
	}()

	bodyBytes, err := io.ReadAll(body)
	require.Equal(t, nil, err)

	var respBody Response
	require.Equal(t, nil, json.Unmarshal(bodyBytes, &respBody))

	require.Equal(t, 1, respBody.Page)
	require.Equal(t, 10, respBody.Limit)
}

// go test -run Test_PaginateCheckDefaultsWhenNoLimit
func Test_PaginateCheckDefaultsWhenNoLimit(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?page=2", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer func() {
		closeErr := body.Close()
		require.Equal(t, nil, closeErr)
	}()

	bodyBytes, err := io.ReadAll(body)
	require.Equal(t, nil, err)

	var respBody Response
	require.Equal(t, nil, json.Unmarshal(bodyBytes, &respBody))

	require.Equal(t, 2, respBody.Page)
	require.Equal(t, 10, respBody.Limit)
}

// go test -run Test_PaginateCheckDefaultsWhenNoPage
func Test_PaginateCheckDefaultsWhenNoPage(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?limit=20", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer func() {
		closeErr := body.Close()
		require.Equal(t, nil, closeErr)
	}()

	bodyBytes, err := io.ReadAll(body)
	require.Equal(t, nil, err)

	var respBody Response
	require.Equal(t, nil, json.Unmarshal(bodyBytes, &respBody))

	require.Equal(t, 1, respBody.Page)
	require.Equal(t, 20, respBody.Limit)
}

// go test -run Test_Paginate_Next
func Test_Paginate_Next(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(New(Config{
		Next: func(_ fiber.Ctx) bool {
			return true
		},
	}))

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// go test -run Test_PaginateConfigDefaultPageDefaultLimit
func Test_PaginateConfigDefaultPageDefaultLimit(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(New(Config{
		DefaultPage:  100,
		DefaultLimit: 1000,
	}))

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer func() {
		closeErr := body.Close()
		require.Equal(t, nil, closeErr)
	}()

	bodyBytes, err := io.ReadAll(body)
	require.Equal(t, nil, err)

	var respBody Response
	require.Equal(t, nil, json.Unmarshal(bodyBytes, &respBody))

	require.Equal(t, 100, respBody.Page)
	require.Equal(t, 1000, respBody.Limit)
}

// go test -run Test_PaginateConfigPageKeyLimitKey
func Test_PaginateConfigPageKeyLimitKey(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(New(Config{
		PageKey:  "p",
		LimitKey: "l",
	}))

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?p=2&l=20", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer func() {
		closeErr := body.Close()
		require.Equal(t, nil, closeErr)
	}()

	bodyBytes, err := io.ReadAll(body)
	require.Equal(t, nil, err)

	var respBody Response
	require.Equal(t, nil, json.Unmarshal(bodyBytes, &respBody))

	require.Equal(t, 2, respBody.Page)
	require.Equal(t, 20, respBody.Limit)
}

// go test -run Test_PaginateNegativeDefaultPageDefaultLimitValues
func Test_PaginateNegativeDefaultPageDefaultLimitValues(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(New(Config{
		DefaultPage:  -1,
		DefaultLimit: -1,
	}))

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer func() {
		closeErr := body.Close()
		require.Equal(t, nil, closeErr)
	}()

	bodyBytes, err := io.ReadAll(body)
	require.Equal(t, nil, err)

	var respBody Response
	require.Equal(t, nil, json.Unmarshal(bodyBytes, &respBody))

	require.Equal(t, 1, respBody.Page)
	require.Equal(t, 10, respBody.Limit)
}

// go test -run Test_PaginateFromContextWithoutNew
func Test_PaginateFromContextWithoutNew(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.Equal(t, nil, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
