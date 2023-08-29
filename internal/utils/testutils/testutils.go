package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unable to identify current directory (needed to load .env.test)")
		os.Exit(1)
	}
	basepath := filepath.Dir(file)

	err := godotenv.Load(filepath.Join(basepath, "../../../.env"))
	if err!= nil {
		fmt.Fprint(os.Stderr, "Unable to load .env file")
	}
}