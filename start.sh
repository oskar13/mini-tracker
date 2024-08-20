export DB_USER=minitorrentuser
export DB_PASSWORD_FILE=local_db_password.txt
export DB_NAME=minitorrent
export DB_HOST=127.0.0.1
export DB_PORT=3306
export INSTALLER_TOKEN_FILE=installer_token.txt

go run cmd/mini-tracker/main.go