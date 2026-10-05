package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

func TestValidator_Struct(t *testing.T) {
	t.Parallel()

	validator := validation.New()

	type request struct {
		Name  string  `validate:"required,min=2,max=10"`
		Price float64 `validate:"required,gt=0"`
	}

	err := validator.Struct(request{
		Name:  "",
		Price: 0,
	})

	require.Error(t, err)

	var validationErr *apperrors.ValidationError

	require.ErrorAs(t, err, &validationErr)

	assert.Contains(t, validationErr.Fields, "name")
	assert.Contains(t, validationErr.Fields, "price")
}

func TestValidator_Struct_Valid(t *testing.T) {
	t.Parallel()

	validator := validation.New()

	type request struct {
		Name  string  `validate:"required,min=2,max=10"`
		Price float64 `validate:"required,gt=0"`
	}

	err := validator.Struct(request{
		Name:  "iPhone",
		Price: 999,
	})

	assert.NoError(t, err)
}

func TestValidator_Struct_Messages(t *testing.T) {
	t.Parallel()

	validator := validation.New()

	type request struct {
		Name  string  `validate:"required,min=2,max=10"`
		Price float64 `validate:"required,gt=0"`
	}

	err := validator.Struct(request{
		Name:  "",
		Price: 0,
	})

	require.Error(t, err)

	var validationErr *apperrors.ValidationError

	require.ErrorAs(t, err, &validationErr)

	assert.Equal(t, "is required", validationErr.Fields["name"])
	assert.Equal(t, "is required", validationErr.Fields["price"])
}

func TestValidator_Struct_MinMaxGt(t *testing.T) {
	t.Parallel()

	validator := validation.New()

	type request struct {
		Name  string  `validate:"min=2,max=10"`
		Price float64 `validate:"gt=0"`
	}

	err := validator.Struct(request{
		Name:  "A",
		Price: 0,
	})

	require.Error(t, err)

	var validationErr *apperrors.ValidationError

	require.ErrorAs(t, err, &validationErr)

	assert.Equal(
		t,
		"must be at least 2 characters",
		validationErr.Fields["name"],
	)

	assert.Equal(
		t,
		"must be greater than 0",
		validationErr.Fields["price"],
	)
}

func TestValidator_Struct_Max(t *testing.T) {
	t.Parallel()

	validator := validation.New()

	type request struct {
		Name string `validate:"max=5"`
	}

	err := validator.Struct(request{
		Name: "abcdef",
	})

	require.Error(t, err)

	var validationErr *apperrors.ValidationError

	require.ErrorAs(t, err, &validationErr)

	assert.Equal(
		t,
		"must be at most 5 characters",
		validationErr.Fields["name"],
	)
}

func TestValidator_Struct_UnsupportedMessageTag(t *testing.T) {
	t.Parallel()

	validator := validation.New()

	type request struct {
		Email string `validate:"email"`
	}

	err := validator.Struct(request{
		Email: "invalid-email",
	})

	require.Error(t, err)

	var validationErr *apperrors.ValidationError

	require.ErrorAs(t, err, &validationErr)

	assert.Equal(
		t,
		"is invalid",
		validationErr.Fields["email"],
	)
}

func TestValidator_Struct_UsesJSONFieldName(t *testing.T) {
	t.Parallel()

	validator := validation.New()

	type request struct {
		ProductName string `json:"productName" validate:"required"`
	}

	err := validator.Struct(request{})

	require.Error(t, err)

	var validationErr *apperrors.ValidationError

	require.ErrorAs(t, err, &validationErr)

	assert.Contains(
		t,
		validationErr.Fields,
		"productName",
	)

	assert.NotContains(
		t,
		validationErr.Fields,
		"productname",
	)
}

func TestValidator_Struct_JSONFieldNameWithOptions(t *testing.T) {
	t.Parallel()

	validator := validation.New()

	type request struct {
		ProductName string `json:"productName,omitempty" validate:"required"`
	}

	err := validator.Struct(request{})

	require.Error(t, err)

	var validationErr *apperrors.ValidationError

	require.ErrorAs(t, err, &validationErr)

	assert.Equal(
		t,
		"is required",
		validationErr.Fields["productName"],
	)
}
