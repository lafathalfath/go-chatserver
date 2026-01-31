package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

func Env(name string) string {
	if err := godotenv.Load(); err != nil {
		panic("[!] .env file not found")
	}
	return os.Getenv(name)
}