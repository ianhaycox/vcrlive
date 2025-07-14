package secretsstore

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGet(t *testing.T) {
	t.Run("Test Get is successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storer := NewMockSecretsStorer(ctrl)
		storer.EXPECT().Get("test").MaxTimes(1).MinTimes(1).Return("test", nil)

		secretsstore := NewSecretsStore(storer)

		result, err := secretsstore.Get("test")

		assert.NoError(t, err)
		assert.Equal(t, result, "test")
	})

	t.Run("Test Get errors successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storer := NewMockSecretsStorer(ctrl)
		storer.EXPECT().Get("test").MaxTimes(1).MinTimes(1).Return("", errors.New("test error"))

		secretsstore := NewSecretsStore(storer)

		_, err := secretsstore.Get("test")

		assert.Error(t, err)
	})
}

func TestSet(t *testing.T) {
	t.Run("Test Set is successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storer := NewMockSecretsStorer(ctrl)
		storer.EXPECT().Set("testname", "testvalue").MaxTimes(1).MinTimes(1).Return(nil)

		secretsstore := NewSecretsStore(storer)

		err := secretsstore.Set("testname", "testvalue")

		assert.NoError(t, err)
	})

	t.Run("Test Set errors successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storer := NewMockSecretsStorer(ctrl)
		storer.EXPECT().Set("testname", "testvalue").MaxTimes(1).MinTimes(1).Return(errors.New("test error"))

		secretsstore := NewSecretsStore(storer)

		err := secretsstore.Set("testname", "testvalue")

		assert.Error(t, err)
	})
}
