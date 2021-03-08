package repositories

import (
	"os"
	"project/domain/repositories/testutils"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.Setup(m))
}
